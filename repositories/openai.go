package repositories

import (
	"autflow_back/interfaces"
	"autflow_back/models"
	"autflow_back/src/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type openaiClient struct {
	httpClient *resty.Client
	cache      *redis.Client
}

func NewOpenAiRepository(cache *redis.Client) interfaces.OpenAIClientRepository {
	client := resty.New().
		SetBaseURL(viper.GetString("GPT_URL")).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+viper.GetString("GPT_GUANABARA"))

	return &openaiClient{
		httpClient: client,
		cache:      cache,
	}
}

func (o *openaiClient) ConvertAudioToText(ctx context.Context, filePath string) (string, error) {
	// Abrir o arquivo de áudio
	file, err := os.Open("../temp_files/" + filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Criar um buffer para a requisição multipart
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Adicionar o arquivo de áudio ao formulário multipart
	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	// Adicionar outros campos ao formulário multipart
	_ = writer.WriteField("model", "whisper-1")

	// Finalizar o formulário multipart
	err = writer.Close()
	if err != nil {
		return "", err
	}

	res, err := o.httpClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", writer.FormDataContentType()).
		SetBody(body).
		Post("/audio/transcriptions")

	if err != nil {
		return "", fmt.Errorf("erro ao buscar locais: %v", err)
	}

	if res.Error() != nil {
		return "", fmt.Errorf("erro na resposta: %v", res.Error())
	}

	// Verificar o status da resposta
	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("status de resposta inesperado: %s", res.Status())
	}

	// Ler o corpo da resposta do novo response (res)
	responseBytes := res.Body()

	// Criar uma estrutura para extrair apenas o campo "text" do JSON
	var jsonResponse struct {
		Text string `json:"text"`
	}

	// Decodificar o JSON no corpo da resposta e extrair o campo "text"
	err = json.Unmarshal(responseBytes, &jsonResponse)
	if err != nil {
		return "", fmt.Errorf("erro ao decodificar JSON: %v", err)
	}

	// Pegar apenas o conteúdo do campo "text"
	responseText := jsonResponse.Text

	return responseText, nil
}

func OpenAIChatCompletion(text string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	// Corpo da requisição
	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]interface{}{
			{"role": "user", "content": text},
		},
		"temperature": 0.7,
	}

	// Codificar o corpo da requisição para JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Criar a requisição POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+config.ChatGPTAPI)
	req.Header.Set("Content-Type", "application/json")

	// Realizar a requisição HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Verificar o status da resposta
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status de resposta inesperado: %s", resp.Status)
	}

	// Ler a resposta da API
	var response models.ChatCompletion
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	// Acessa o conteúdo da primeira escolha
	content := response.Choices[0].Message.Content
	//content = removerQuebrasDeLinha(content)
	//content = removerEspacosEmBranco(content)
	//content = removerEspacosEmBrancoInicioFim(content)

	// Decodifica a string JSON para uma estrutura ou mapa
	/*var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(content), &jsonData); err != nil {
		return err
	}*/

	//fmt.Println("Origem", jsonData["origem"])
	//fmt.Println("Destino", jsonData["destino"])
	fmt.Printf(content)
	return content, nil
}

func GetChatGPTResponse(text string) (string, error) {
	// URL da API do OpenAI
	url := "https://api.openai.com/v1/chat/completions"

	// Corpo da requisição para o OpenAI
	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo-1106",
		"response_format": map[string]string{
			"type": "json_object",
		},
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Você é um assistente que vai identificar os possíveis campos no texto que vai receber: cidade de origem , destino data de ida, data de volta , retorne um JSON com esses campos (origem, destino, dataIda, dataVolta) e se não encontrar os valores retorne os campos nulos",
			},
			{
				"role":    "user",
				"content": text,
			},
		},
	}

	// Codificar o corpo da requisição para JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("Erro ao codificar a requisição: %v", err)
	}

	// Criar a requisição HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return "", fmt.Errorf("Erro ao criar requisição HTTP: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.ChatGPTAPI)
	req.Header.Set("Content-Type", "application/json")

	// Realizar a requisição HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Erro ao fazer a requisição: %v", err)
	}
	defer resp.Body.Close()

	// Ler a resposta da requisição
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Erro ao ler resposta: %v", err)
	}

	var response models.ChatCompletion
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	// Acessa o conteúdo da primeira escolha
	content := string(response.Choices[0].Message.Content)

	// Retorna o conteúdo da primeira escolha
	return content, nil
}

// / Assistant ///

// Create a new thread in gpt
func (o *openaiClient) CreateThread(ctx context.Context) (*models.ThreadResponse, error) {
	res, err := o.httpClient.R().
		SetContext(ctx).
		SetHeader("OpenAI-Beta", "assistants=v1").
		SetBody(map[string]interface{}{}).
		Post("/threads")

	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return nil, fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("criar thread unexpected response status: %s", res.Status())
	}

	var threadResp models.ThreadResponse
	err = json.Unmarshal(res.Body(), &threadResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling locations: %v", err)
	}

	return &threadResp, nil
}

// Add a message to thread
func (o *openaiClient) PostMessage(ctx context.Context, threadID, message string) (string, error) {
	res, err := o.httpClient.R().
		SetContext(ctx).
		SetHeader("OpenAI-Beta", "assistants=v1").
		SetBody(map[string]interface{}{
			"role":    "user",
			"content": message,
		}).
		Post("/threads/" + threadID + "/messages")

	if err != nil {
		return "", fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return "", fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("adicionar mensagen unexpected response status: %s", res.Status())
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(res.Body(), &responseMap)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling locations: %v", err)
	}

	id, ok := responseMap["id"].(string)
	if !ok {
		return "", fmt.Errorf("ID not found in response")
	}

	return id, nil
}

// Run the thread
func (o *openaiClient) StartThreadRun(ctx context.Context, threadID string) (string, error) {
	// Get the current date and time
	currentTime := time.Now()
	currentDay := currentTime.Day()
	currentMonth := currentTime.Month()
	currentYear := currentTime.Year()
	fmt.Println(fmt.Sprintf("A Data atual é: %02d/%02d/%d", currentDay, currentMonth, currentYear))

	res, err := o.httpClient.R().
		SetContext(ctx).
		SetHeader("OpenAI-Beta", "assistants=v1").
		SetBody(map[string]interface{}{
			"assistant_id":            viper.Get("ASSISTANT_ID").(string),
			"additional_instructions": fmt.Sprintf("Hoje é dia: %d-%02d-%02d", currentYear, currentMonth, currentDay),
		}).
		Post("/threads/" + threadID + "/runs")

	if err != nil {
		return "", fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return "", fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("iniciar run unexpected response status: %s", res.Status())
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(res.Body(), &responseMap)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling locations: %v", err)
	}

	id, ok := responseMap["id"].(string)
	if !ok {
		return "", fmt.Errorf("ID not found in response")
	}

	return id, nil
}

// Search the current status of the Thread
func (o *openaiClient) GetThreadRunStatus(ctx context.Context, threadID, runID string) (*models.ThreadRun, error) {
	res, err := o.httpClient.R().
		SetContext(ctx).
		SetHeader("OpenAI-Beta", "assistants=v1").
		Get("/threads/" + threadID + "/runs/" + runID)

	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return nil, fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Status thread unexpected response status: %s", res.Status())
	}

	var threadRun models.ThreadRun
	err = json.Unmarshal(res.Body(), &threadRun)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling locations: %v", err)
	}

	return &threadRun, nil
}

// Search for the last message in the thread
func (o *openaiClient) GetThreadMessages(ctx context.Context, threadID string) ([]models.MessageThread, error) {
	res, err := o.httpClient.R().
		SetContext(ctx).
		SetHeader("OpenAI-Beta", "assistants=v1").
		Get("/threads/" + threadID + "/messages")

	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return nil, fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("mensagens unexpected response status: %s", res.Status())
	}

	var response struct {
		Data []models.MessageThread `json:"data"`
	}
	err = json.Unmarshal(res.Body(), &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// Reply gpt function
func (o *openaiClient) PostToolOutputs(ctx context.Context, threadID, runID, callID string, arrayRespone []models.CallResponse) (string, error) {

	res, err := o.httpClient.R().
		SetContext(ctx).
		SetHeader("OpenAI-Beta", "assistants=v1").
		SetBody(map[string]interface{}{
			"tool_outputs": arrayRespone,
		}).
		Post("/threads/" + threadID + "/runs/" + runID + "/submit_tool_outputs")

	if err != nil {
		return "", fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return "", fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("postar respostas unexpected response status: %s", res.Status())
	}

	fmt.Println("Requisição POST bem-sucedida")
	return "", nil
}

// Cancel thread run
func (o *openaiClient) CancelRun(ctx context.Context, threadID, runID string) (string, error) {
	fmt.Println("Vamos cancelar um run")
	res, err := o.httpClient.R().
		SetContext(ctx).
		SetHeader("OpenAI-Beta", "assistants=v1").
		Post("/threads/" + threadID + "/runs/" + runID + "/cancel")

	if err != nil {
		return "", fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return "", fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("cancelar run unexpected response status: %s", res.Status())
	}

	fmt.Println("Requisição POST bem-sucedida")
	return "", nil
}
