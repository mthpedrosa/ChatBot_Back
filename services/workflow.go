package services

import (
	"autflow_back/interfaces"
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/repositories"
	"autflow_back/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workflow struct {
	workflowRepository      *repositories.Workflows
	metaRepository          *repositories.Metas
	customerRepository      *repositories.Customers
	sessionRepository       *repositories.Session
	conversarionsRepository *repositories.Conversations
	logger                  utils.Logger
	openaiRepository        interfaces.OpenAIClientRepository
	whatsappRepository      interfaces.WhatsappRepository
}

func NewWorkflow(workflow *repositories.Workflows,
	meta *repositories.Metas,
	customer *repositories.Customers,
	session *repositories.Session,
	conversarions *repositories.Conversations,
	logger utils.Logger, openai interfaces.OpenAIClientRepository, whatsapp interfaces.WhatsappRepository) *Workflow {
	return &Workflow{
		workflowRepository:      workflow,
		metaRepository:          meta,
		logger:                  logger,
		customerRepository:      customer,
		sessionRepository:       session,
		conversarionsRepository: conversarions,
		//guanabaraRepository:     guanabara,
		openaiRepository:   openai,
		whatsappRepository: whatsapp,
	}
}

func (r *Workflow) Create(ctx context.Context, dt *dto.CreateWorkflowDTO) (string, error) {
	workflow := dt.ToWorkflow()

	if workflow.PhoneMetaId != "" {
		//metasRepo := repositories.NewMetaRepository(r.db)

		metaValida, erro := r.metaRepository.FindPhoneID(ctx, workflow.PhoneMetaId)
		if erro != nil {
			return "", erro

		}

		if metaValida == nil {
			return "", errors.New("Conta meta não existente")

		}

	}

	//repositorio := repositories.NewWorkflowsRepository(r.db)
	//idCriado, erro := r.workflowRepository.Insert(ctx, workflow)
	idCriado, erro := r.workflowRepository.Insert(ctx, workflow)
	if erro != nil {
		return "", erro
	}

	return idCriado, nil
}

func (r *Workflow) Find(ctx context.Context, query string) ([]models.Workflow, error) {

	//repositorio := repositories.NewWorkflowsRepository(r.db)
	workflows, erro := r.workflowRepository.Find(ctx, query)
	if erro != nil {
		return nil, erro
	}

	return workflows, nil
}

func (r *Workflow) FindId(ctx context.Context, id string) (models.Workflow, error) {

	workflow, erro := r.workflowRepository.FindId(ctx, id)
	if erro != nil {
		return models.Workflow{}, erro
	}

	return *workflow, nil
}

func (r *Workflow) Edit(ctx context.Context, id string, dt *dto.CreateWorkflowDTO) error {
	newWorkflow := dt.ToWorkflow()
	return r.workflowRepository.Edit(ctx, id, newWorkflow)
}

func (r *Workflow) Delete(ctx context.Context, id string) error {
	return r.workflowRepository.Delete(ctx, id)
}

// ----- go through workflow -----

// Identifies active workflow by meta id
func (r *Workflow) IdentifyWorkflow(ctx context.Context, payload models.WebhookPayload, meta models.Meta) (models.Workflow, error) {
	//Checks if this meta account has the correct id
	fmt.Print("id meta que chegou : " + payload.Entry[0].Changes[0].Value.MetaData.PhoneNumberId)
	var check bool
	check = false
	for _, phones := range meta.PhonesMeta {
		fmt.Println("id meta da conta : " + phones.Id)
		if phones.Id == payload.Entry[0].Changes[0].Value.MetaData.PhoneNumberId {
			check = true
		}
	}

	if !check {
		return models.Workflow{}, errors.New("Não foi possivel localizar a conta meta")
	}

	query := fmt.Sprintf("phone_meta_id=" + payload.Entry[0].Changes[0].Value.MetaData.PhoneNumberId + "&active=true")
	workflow, err := r.workflowRepository.Find(ctx, query)
	if err != nil {
		return models.Workflow{}, err
	}
	if len(workflow) == 0 {
		return models.Workflow{}, errors.New("Nenhum worklow ativo")
	}

	return workflow[0], nil
}

func (r *Workflow) RunWorkflow(ctx context.Context, payload models.WebhookPayload, meta models.Meta, workflow models.Workflow) error {
	var session models.Session
	var idConversation string

	// checking if the client exists, if not it creates
	customerExist, err, customer, idCustomer := r.customerExist(ctx, payload)
	if err != nil {
		return err
	}

	if !customerExist {
		session, idConversation, err = r.startConversationSession(ctx, idCustomer, workflow.ID.Hex())
		if err != nil {
			return err
		}

		fmt.Println("Cliente não existia, aqui deve ter sido criado uma conversa e uma sessão")
	} else {
		session, idConversation, err = r.getInfoSessions(ctx, customer.ID.Hex(), workflow.ID.Hex())
		if err != nil {
			return err
		}
	}

	fmt.Println(session)
	fmt.Println(idConversation)

	// identify the type of message received and its content
	message, err := r.identifyMessage(ctx, payload)
	if err != nil {
		return err
	}

	//Get meta and workflow information
	var business string
	for _, phones := range meta.PhonesMeta {
		if phones.Id == workflow.PhoneMetaId {
			business = phones.BusinessId
		}
	}

	metaTokens := models.MetaIds{
		PhoneID:    workflow.PhoneMetaId,
		Token:      meta.Token,
		BusinessId: business,
	}

	//Convert audio for text
	if message.Type == "audio" {
		fmt.Println("Entro no audio")
		responseMedia, err := r.whatsappRepository.GetUrlMedia(ctx, message.Content, metaTokens.Token)
		if err != nil {
			return err
		}

		nameFile, err := r.whatsappRepository.DownloadMedia(ctx, responseMedia, metaTokens.Token, message.Content)
		if err != nil {
			return err
		}

		messageAudio, err := r.openaiRepository.ConvertAudioToText(ctx, nameFile)
		if err != nil {
			return err
		}

		fmt.Println("Audio trauzido para : " + messageAudio)
		message.Content = messageAudio

		//remove file
		err = os.Remove("../temp_files" + nameFile)
	}

	// Check session stage
	step, err := utils.OtherFields(session.OtherFields, "current_stage")
	if err != nil {
		return err
	}
	fmt.Println(step)

	// Data struct mount
	var flowData dto.FlowData
	flowData.Ctx = ctx
	flowData.Customer = customer
	flowData.Session = session
	flowData.Message = message
	flowData.MetaTokens = metaTokens

	// Allow the customer to return to the main menu in chatgpt
	if message.Content == "v" || message.Content == "V" {
		if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "current_stage", "string", "assistant"); err != nil {
			return err
		}

		var arguments = `{"first_message":false}`
		err = r.SendMenu(flowData, arguments)
		return nil
	}

	// Identifies whether to proceed to the assistant
	if step == "assistant" || step == "" {
		//Send to gpt assistant
		err, threadsIds := r.gpt_assistant(flowData)
		if err != nil {
			//_ = r.whatsappRepository.SimpleMessage(ctx, "Desculpe, parece que tivemos um problema técnico. Vou reiniciar nosso processo para garantir que tudo funcione corretamente. Por favor, me informe novamente, qual informação você está buscando? Se precisar de assistência em algum tópico específico, estou aqui para ajudar!", customer, metaTokens)
			_ = r.whatsappRepository.SimpleMessage(ctx, err.Error(), customer, metaTokens)

			if threadsIds.ThreadId != "" && threadsIds.RunId != "" {
				_, err = r.openaiRepository.CancelRun(ctx, threadsIds.ThreadId, threadsIds.RunId)
			}

			return err
		}
	}

	// Identify whether the customer is in the sales stage
	/*if step == "sale" {
		err = r.FlowSale(flowData, "")
		if err != nil {
			_ = r.whatsappRepository.SimpleMessage(ctx, err.Error(), customer, metaTokens)
		}
	}*/

	return nil
}

func (r *Workflow) customerExist(ctx context.Context, payload models.WebhookPayload) (bool, error, models.Customer, string) {
	if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 && len(payload.Entry[0].Changes[0].Value.Contacts) > 0 {
		waID := payload.Entry[0].Changes[0].Value.Contacts[0].WAID
		fmt.Println("WA ID:", waID)

		customer, err := r.customerRepository.Find(ctx, "whatsapp_id="+waID)
		if err != nil || len(customer) == 0 {
			// Customer not found, we created it
			newCustomer := models.Customer{
				WhatsAppID: waID,
				Phone:      waID,
				Name:       payload.Entry[0].Changes[0].Value.Contacts[0].Profile.Name,
			}

			id, erro := r.customerRepository.Insert(ctx, newCustomer)
			if erro != nil {
				return false, erro, models.Customer{}, "" // Error when adding a new customer
			}

			newCustomer.ID, err = primitive.ObjectIDFromHex(id)
			if err != nil {
				return false, err, models.Customer{}, ""
			}
			return false, nil, newCustomer, id // Customer created with success
		}
		return true, nil, customer[0], "" // Customer found
	}

	return false, errors.New("Falha ao ler body da solicitação"), models.Customer{}, "" // Request body reading error
}

func (r *Workflow) startConversationSession(ctx context.Context, idCustomer, workflowId string) (models.Session, string, error) {
	var conversations models.Conversation
	conversations.CustomerId = idCustomer
	conversations.WorkflowID = workflowId

	idConversation, err := r.conversarionsRepository.Insert(ctx, conversations)
	if err != nil {
		return models.Session{}, "", err
	}

	var session models.Session
	session.ConversationId = idConversation
	session.CustomerID = idCustomer
	session.WorkflowId = workflowId
	session.Status = "in_progress"
	idSession, err := r.sessionRepository.Insert(ctx, session)
	if err != nil {
		return models.Session{}, "", err
	}

	return idSession, idConversation, nil
}

func (r *Workflow) getInfoSessions(ctx context.Context, idCustomer, workflowId string) (models.Session, string, error) {
	var newSession models.Session
	var idConversation string

	var query = fmt.Sprintf("customer_id=%s&workflow_id=%s", idCustomer, workflowId)
	sessions, err := r.sessionRepository.Find(ctx, query)
	if err != nil {
		return models.Session{}, "", err
	}

	fmt.Println("Secoes encontradas: ")
	fmt.Println(sessions)
	for _, session := range sessions {
		fmt.Println("Dentro do for")
		if session.Status == "in_progress" {
			fmt.Println("Procurando  sessoes em progresso : " + session.ID.Hex())
			newSession = session
			idConversation = session.ConversationId
		}
	}

	if newSession.ID.Hex() == "000000000000000000000000" || len(sessions) == 0 {
		var conversations models.Conversation
		conversations.CustomerId = idCustomer
		conversations.WorkflowID = workflowId

		/*idConversation, err := r.conversarionsRepository.Insert(ctx, conversations)
		if err != nil {
			return "", "", err
		}*/

		//Não cria uma nova conversa em toda sessão
		fmt.Println("vou pesquisar")
		conversation, err := r.conversarionsRepository.Find(ctx, "customer_id="+idCustomer)
		if err != nil {
			return models.Session{}, "", err
		}
		idConversation = conversation[0].ID.Hex()

		var newSession models.Session
		newSession.ConversationId = idConversation
		newSession.CustomerID = idCustomer
		newSession.WorkflowId = workflowId
		newSession.Status = "in_progress"
		newSession, err = r.sessionRepository.Insert(ctx, newSession)
		if err != nil {
			return models.Session{}, "", err
		}

	}

	return newSession, idConversation, nil
}

func (r *Workflow) gpt_assistant(flowData dto.FlowData) (error, models.ThreadIds) {
	var status string
	var callID string
	var response string
	var threadsId models.ThreadIds
	var ingoreReturn bool

	//Checks if the client has a thread created
	threadID, err := utils.OtherFields(flowData.Customer.OtherFields, "idThread")
	if err != nil {
		return err, threadsId
	}

	if threadID == "" {
		fmt.Println("Vamos criar a thread")
		thread, err := r.openaiRepository.CreateThread(flowData.Ctx)
		if err != nil {
			return err, threadsId
		}

		fmt.Println("Thread Criada :")
		fmt.Print(thread)
		threadID = thread.ID

		var field models.Fields
		field.Name = "idThread"
		field.Type = "string"
		field.Value = threadID

		fmt.Println("Campo a ser criado/atualizado: " + field.Name)
		fmt.Println("valor a ser criado/atualizado: " + field.Value)

		err = r.customerRepository.UpdateCustomerField(flowData.Ctx, flowData.Session.CustomerID, field)
		if err != nil {
			return err, threadsId
		}
	}

	// add message to thread
	messageID, err := r.openaiRepository.PostMessage(flowData.Ctx, threadID, flowData.Message.Content)
	if err != nil {
		return err, threadsId
	}
	fmt.Println("Id da mensagem adicionada: " + messageID)

	// Start the thread
	runID, err := r.openaiRepository.StartThreadRun(flowData.Ctx, threadID)
	if err != nil {
		return err, threadsId
	}
	fmt.Println("Id do run: " + runID)

	//Define the struc to return on error
	threadsId.ThreadId = threadID
	threadsId.MessageId = messageID
	threadsId.RunId = runID

	// Query thread run status
	for {
		fmt.Println("No for")
		threadJson, err := r.openaiRepository.GetThreadRunStatus(flowData.Ctx, threadID, runID)
		if err != nil {
			return err, threadsId
		}

		status = threadJson.Status
		fmt.Println("Status:", status)

		if status == "completed" || status == "cancelled" || status == "failed" {
			fmt.Println("Status:", threadJson)
			break
		}

		if status == "requires_action" {
			var arrayRespone []models.CallResponse

			for _, call := range threadJson.RequiredAction.SubmitToolOutputs.ToolCalls {
				var format models.CallResponse
				callID = call.ID
				fmt.Println("Numero da chamada : " + callID)

				switch call.Function.Name {
				case "get_schedules":
					//err = r.FlowSale(flowData, call.Function.Arguments)
					if err != nil {
						return err, models.ThreadIds{}
					}

					format = models.CallResponse{
						ToolCallID: callID,
						OutPut:     `{success:"true"}`,
					}

					arrayRespone = append(arrayRespone, format)
					ingoreReturn = true
				case "save_name":
					var result map[string]string

					err = json.Unmarshal([]byte(call.Function.Arguments), &result)
					if err != nil {
						format = models.CallResponse{
							ToolCallID: callID,
							OutPut:     `{success:"false"}`,
						}
					}

					var field models.Fields
					field.Name = "name_provided"
					field.Type = "string"
					field.Value = result["first_name"] + " " + result["last_name"]
					err = r.customerRepository.UpdateCustomerField(flowData.Ctx, flowData.Session.CustomerID, field)

					format = models.CallResponse{
						ToolCallID: callID,
						OutPut:     `{success:"true"}`,
					}

					arrayRespone = append(arrayRespone, format)

					break

				case "send_menu":
					err = r.SendMenu(flowData, call.Function.Arguments)
					if err != nil {
						return err, models.ThreadIds{}
					}

					format = models.CallResponse{
						ToolCallID: callID,
						OutPut:     `{success:"true"}`,
					}

					arrayRespone = append(arrayRespone, format)
					ingoreReturn = true
				default:
					return fmt.Errorf("Função nao encontrada"), threadsId
				}
			}

			if len(arrayRespone) > 0 {
				_, err = r.openaiRepository.PostToolOutputs(flowData.Ctx, threadID, runID, callID, arrayRespone)
				if err != nil {
					return err, threadsId
				}
			}
		}
		time.Sleep(100 * time.Millisecond) // test
	}

	if status != "completed" && !ingoreReturn {
		return errors.New("error status : " + status), models.ThreadIds{}
	}

	if !ingoreReturn {
		// Search chat messages
		mensagens, err := r.openaiRepository.GetThreadMessages(flowData.Ctx, threadID)
		if err != nil {
			return err, threadsId
		}

		fmt.Println("Mensagens da busca:")
		fmt.Println(mensagens)
		for _, message := range mensagens {
			if message.RunID == runID {
				response = message.Content[0].Text.Value
			}
		}

		// We send the chat response to the customer
		err = r.whatsappRepository.SimpleMessage(flowData.Ctx, response, flowData.Customer, flowData.MetaTokens)
		if err != nil {
			return err, threadsId
		}
	}

	return nil, models.ThreadIds{}
}

// Identifies the type of message received - If it's audio, download it and get treatment
func (r *Workflow) identifyMessage(ctx context.Context, payload models.WebhookPayload) (models.MessagePayload, error) {
	var messagePayload models.MessagePayload
	// Determine the message type
	messageType := payload.Entry[0].Changes[0].Value.Messages[0].Type

	if messageType == "interactive" {
		if payload.Entry[0].Changes[0].Value.Messages[0].Interactive.Type == "button_reply" {
			messagePayload.Type = " button"
			messagePayload.Content = payload.Entry[0].Changes[0].Value.Messages[0].Interactive.ButtonReply.Id
		} else if payload.Entry[0].Changes[0].Value.Messages[0].Interactive.Type == "list_reply" {
			messagePayload.Type = " list"
			messagePayload.Content = payload.Entry[0].Changes[0].Value.Messages[0].Interactive.ListReply.Id
		}
	} else if messageType == "text" {
		messagePayload.Type = " text"
		messagePayload.Content = utils.RemoveInvalidCharacters(payload.Entry[0].Changes[0].Value.Messages[0].Text.Body)

	} else if messageType == "audio" {
		messagePayload.Type = "audio"
		messagePayload.Content = payload.Entry[0].Changes[0].Value.Messages[0].Audio.Id
	}
	return messagePayload, nil
}

// Facilitates the creation of fields
func (r *Workflow) updateSessionField(ctx context.Context, sessionId string, name string, fieldType string, value string) error {
	field := models.Fields{
		Name:  name,
		Type:  fieldType,
		Value: value,
	}

	return r.sessionRepository.UpdateSessionField(ctx, sessionId, field)
}

// Send message for customer with menu
func (r *Workflow) SendMenu(flowData dto.FlowData, arguments string) error {

	var nodeMenu models.Node
	var messageStart = "Oi! Tudo bem? \n \nEu sou a Ana, a assistente virtual do Grupo Guanabara! Tô aqui pra te ajudar a viver experiências incríveis de viagem. 🤗 \n\n"
	var messageMenu = "Sobre o que vamos falar hoje? *Você quer conferir ou comprar bilhetes de viagem? Já comprou e precisa de ajuda? Ou quer saber mais sobre algum assunto?* É só me contar!"
	var finalMessage string

	finalMessage = messageStart + messageMenu

	button1 := models.Button{ID: "pesquisar_passagem", Title: "COMPRAR PASSAGEM", NextNode: "node1"}
	button2 := models.Button{ID: "meus_bilhetes", Title: "MEUS BILHETES", NextNode: "node2"}
	button3 := models.Button{ID: "outros_assuntos", Title: "OUTROS ASSUNTOS", NextNode: "node3"}
	nodeMenu.Parameters.Buttons = []models.Button{button1, button2, button3}

	r.whatsappRepository.InteractiveMessage(flowData.Ctx, finalMessage, []models.Button{button1, button2, button3}, flowData.Customer, flowData.MetaTokens)
	return nil
}
