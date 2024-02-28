package repositories

import (
	"autflow_back/interfaces"
	"autflow_back/models"
	"autflow_back/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RespostaMedia struct {
	URL string `json:"url"`
}

type whatsappClient struct {
	httpClient *resty.Client
	cache      *redis.Client
}

func NewWhatsappRepository(cache *redis.Client) interfaces.WhatsappRepository {
	client := resty.New().
		SetBaseURL(viper.GetString("WP_URL")).
		SetHeader("Content-Type", "application/json")
	return &whatsappClient{
		httpClient: client,
		cache:      cache,
	}
}

func (w *whatsappClient) InteractiveMessage(ctx context.Context, text string, buttonsArray []models.Button, customer models.Customer, meta models.MetaIds) error {
	// Building buttons
	var buttons []map[string]interface{}
	for _, button := range buttonsArray {
		buttonData := map[string]interface{}{
			"type": "reply",
			"reply": map[string]interface{}{
				"id":    button.ID,    // Button ID
				"title": button.Title, // Button Title
			},
		}
		buttons = append(buttons, buttonData)
	}

	res, err := w.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+meta.Token).
		SetBody(map[string]interface{}{
			"recipient_type": "individual",
			"to":             customer.WhatsAppID,
			"type":           "interactive",
			"interactive": map[string]interface{}{
				"type": "button",
				"body": map[string]interface{}{
					"text": text,
				},
				"action": map[string]interface{}{
					"buttons": buttons,
				},
			},
			"messaging_product": "whatsapp",
		}).
		Post(meta.PhoneID + "/messages")

	if err != nil {
		return fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", res.Status())
	}

	fmt.Println("Responsta do envio de mensagem")
	fmt.Println(res)
	return nil
}

func (w *whatsappClient) InteractiveMessageList(ctx context.Context, customer models.Customer, meta models.MetaIds, bodyText string, rows []models.Row) error {
	res, err := w.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+meta.Token).
		SetBody(map[string]interface{}{
			"messaging_product": "whatsapp",
			"recipient_type":    "individual",
			"to":                customer.WhatsAppID, // Substitua por um número de telefone específico se necessário
			"type":              "interactive",
			"interactive": map[string]interface{}{
				"type": "list",
				"header": map[string]interface{}{
					"type": "text",
					"text": "",
				},
				"body": map[string]interface{}{
					"text": bodyText,
				},
				"footer": map[string]interface{}{
					"text": "",
				},
				"action": map[string]interface{}{
					"button": "Opções",
					"sections": []map[string]interface{}{
						{
							"title": "Opções",
							"rows":  rows,
						},
					},
				},
			},
		}).
		SetDebug(true).
		Post(meta.PhoneID + "/messages")

	if err != nil {
		return fmt.Errorf("error sending interactive message: %v", err)
	}

	if res.Error() != nil {
		return fmt.Errorf("error in response: %v", res.Error())
	}

	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", res.Status())
	}

	return nil
}

func (w *whatsappClient) SimpleMessage(ctx context.Context, messageSend string, customer models.Customer, meta models.MetaIds) error {
	res, err := w.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+meta.Token).
		SetBody(map[string]interface{}{
			"messaging_product": "whatsapp",
			"to":                customer.WhatsAppID,
			"text": map[string]interface{}{
				"body": messageSend,
			},
		}).
		Post(meta.PhoneID + "/messages")

	if err != nil {
		return fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", res.Status())
	}

	return nil
}

func (w *whatsappClient) InteractiveListMessage(ctx context.Context, customer models.Customer, meta models.MetaIds) error {
	url := viper.Get("WP_ENDPOINT").(string) + meta.PhoneID + "/messages"

	// Corpo da requisição (body)
	requestBody := map[string]interface{}{
		"recipient_type": "individual",
		"to":             customer.WhatsAppID,
		"type":           "interactive",
		"interactive": map[string]interface{}{
			"type": "list",
			"list": map[string]interface{}{
				"button": map[string]interface{}{
					"title": "Selecione uma opção:",
					"items": []map[string]interface{}{
						{
							"media": map[string]interface{}{
								"url":          "https://example.com/image.png", // URL da imagem (opcional)
								"content-type": "image/png",                     // Tipo de conteúdo (opcional)
							},
							"description": "Descrição da opção 1",
							"title":       "Opção 1",
							"buttons": []map[string]interface{}{
								{
									"type":  "reply",
									"title": "Responder 1",
									"reply": map[string]interface{}{
										"payload": "payload_opcao_1",
									},
								},
							},
						},
						{
							"media": map[string]interface{}{
								"url":          "https://example.com/image2.png", // URL da imagem (opcional)
								"content-type": "image/png",                      // Tipo de conteúdo (opcional)
							},
							"description": "Descrição da opção 2",
							"title":       "Opção 2",
							"buttons": []map[string]interface{}{
								{
									"type":  "reply",
									"title": "Responder 2",
									"reply": map[string]interface{}{
										"payload": "payload_opcao_2",
									},
								},
							},
						},
					},
				},
			},
		},
		"messaging_product": "whatsapp",
	}

	// Codificar o corpo da requisição para JSON
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// Realizar a requisição POST
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificar o status da resposta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status de resposta inesperado: %s", resp.Status)
	}

	return nil
}

func (w *whatsappClient) GetUrlMedia(ctx context.Context, mediaID string, accessToken string) (string, error) {
	res, err := w.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		Get(mediaID)

	if err != nil {
		return "", fmt.Errorf("error fetching locations: %v", err)
	}

	if res.Error() != nil {
		return "", fmt.Errorf("error in response: %v", res.Error())
	}

	// Check the response status
	if res.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("unexpected response status: %s", res.Status())
	}

	// Define a structure (struct) for the expected JSON in the response
	var mediaResponse struct {
		URL string `json:"url"`
	}

	// Decode the JSON from the response body
	err = json.Unmarshal(res.Body(), &mediaResponse)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON response: %v", err)
	}

	// Access the "url" field in the structure
	url := mediaResponse.URL

	return url, nil
}

func (w *whatsappClient) DownloadMedia(ctx context.Context, url, accessToken string, name string) (string, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response was successful (status 200)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Falha ao baixar o arquivo: status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	extension := utils.GetExtension(contentType)

	if extension == "" {
		return "", fmt.Errorf("tipo de conteúdo não suportado: %s", contentType)
	}

	out, err := os.Create("../temp_files/" + name + "." + extension)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Download concluído com sucesso!")
	return (name + "." + extension), nil
}
