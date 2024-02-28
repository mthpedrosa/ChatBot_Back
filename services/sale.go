package services

import (
	"autflow_back/models"
	"autflow_back/models/dto"
	"autflow_back/utils"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// adicionar botoes nos horarios - selecionar ok, vou adicionar o filtro
// deixar fluxo liso e continuar processo de venda
// formatação de alguns textos

// sales flow
func (r *Workflow) FlowSale(flowData dto.FlowData, arguments string) error {
	fmt.Println("Dentro do fluxo de vendas")
	if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "current_stage", "string", "sale"); err != nil {
		return err
	}

	// Process argument fields coming from CHATGPT
	if err := r.processArguments(flowData.Ctx, flowData.Session, arguments); err != nil {
		return err
	}

	// Search for the updated session
	sessionAtt, err := r.sessionRepository.FindId(flowData.Ctx, flowData.Session.ID.Hex())
	if err != nil {
		return err
	}

	// Current sales stage
	stage, err := utils.OtherFields(sessionAtt.OtherFields, "sale_stage")

	// Departure - Stage 1 lists the options, stage 2 selects the choice or lists the rest
	if stage == "1" || stage == "2" {
		flowData.Session = *sessionAtt
		err = r.departureStage(flowData, stage)
		if err != nil {
			return err
		}
	}

	sessionAtt, err = r.sessionRepository.FindId(flowData.Ctx, flowData.Session.ID.Hex())
	if err != nil {
		return err
	}

	// Arrival - Stage 3 lists the options, stage 4 selects the choice or lists the rest
	stage, err = utils.OtherFields(sessionAtt.OtherFields, "sale_stage")
	if stage == "3" || stage == "4" {
		flowData.Session = *sessionAtt
		err = r.arrivalStage(flowData, stage)
		if err != nil {
			return err
		}
	}

	sessionAtt, err = r.sessionRepository.FindId(flowData.Ctx, flowData.Session.ID.Hex())
	if err != nil {
		return err
	}
	stage, err = utils.OtherFields(sessionAtt.OtherFields, "sale_stage")

	// Valid route
	if stage == "5" || stage == "6" {
		flowData.Session = *sessionAtt
		err = r.validationSchedules(flowData, stage)
		if err != nil {
			return err
		}
	}

	return nil
}

// Process fields coming from the assistant
func (r *Workflow) processArguments(ctx context.Context, session models.Session, arguments string) error {
	if arguments != "" {
		// Convert a string JSON for struct
		var tripRequest models.TripRequest
		if err := json.Unmarshal([]byte(arguments), &tripRequest); err != nil {
			fmt.Println("Erro ao fazer o parsing do JSON:", err)
			return err
		}

		// Atualizar campos relacionados à viagem (departure, arrival, etc.)
		if tripRequest.Departure != "" {
			if err := r.updateSessionField(ctx, session.ID.Hex(), "departure_name", "string", tripRequest.Departure); err != nil {
				return err
			}
		}

		if tripRequest.Arrival != "" {
			if err := r.updateSessionField(ctx, session.ID.Hex(), "arrival_name", "string", tripRequest.Arrival); err != nil {
				return err
			}
		}

		if tripRequest.Data != "" {
			if err := r.updateSessionField(ctx, session.ID.Hex(), "departure_data", "string", tripRequest.Data); err != nil {
				return err
			}
		}

		// Atualizar campo "sale_stage"
		if err := r.updateSessionField(ctx, session.ID.Hex(), "sale_stage", "string", "1"); err != nil {
			return err
		}
	}
	return nil
}

// List the options compatible with the city and save the chosen one
func (r *Workflow) departureStage(flowData dto.FlowData, stage string) error {
	var messageResponse string
	var formattedMessage string

	//Get departure Options
	departureName, err := utils.OtherFields(flowData.Session.OtherFields, "departure_name")
	if err != nil {
		return err
	}

	if stage == "1" {
		fmt.Println("Dentro de departure stage 1")
		responseFunc, err := r.guanabaraRepository.GetCity(flowData.Ctx, departureName)
		if err != nil {
			return err
		}

		if len(responseFunc) > 0 {
			var button []models.Button
			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "departure_options", "string", strings.Join(responseFunc, ";")); err != nil {
				return err
			}

			var controle int
			// build a formatted message (up to the fifth element)
			for index, messageArray := range responseFunc {
				if index >= 5 {
					break
				}
				controle++

				formattedMessage += fmt.Sprintf("*%d - %s*", index+1, strings.SplitN(messageArray, "-", 2)[1])

				if index < 4 {
					formattedMessage += "\n"
				}
			}
			var complement string
			if len(responseFunc) > 5 {
				qtd := len(responseFunc) - controle
				complement = fmt.Sprintf("\n•+ para ver mais opções (foram encontradas mais %d opções)", qtd)
				button = []models.Button{
					{ID: "+", Title: fmt.Sprintf("Listar (%d opções)", qtd), NextNode: "node1"},
				}
			}

			messageResponse = fmt.Sprintf("Selecione a sua origem: 📍 \n%s \n\n *Confirma pra mim?* Só digitar o número acima. \n \n Ou, se quiser, digite:%s \n*• V para voltar.*", formattedMessage, complement)
			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "sale_stage", "string", "2"); err != nil {
				return err
			}

			err = r.whatsappRepository.InteractiveMessage(flowData.Ctx, messageResponse, button, flowData.Customer, flowData.MetaTokens)
			//err = r.whatsappRepository.SimpleMessage(flowData.Ctx, messageResponse, flowData.Customer, flowData.MetaTokens)
			fmt.Println(messageResponse)

		} else {
			flowData.Message.Content = "origem_nao_encontrada"
			err = r.errorAssitantReturn(flowData)
			return nil
		}

	}

	//Saves the selected departure id
	if stage == "2" {
		departure_options_string, err := utils.OtherFields(flowData.Session.OtherFields, "departure_options")
		departure_optinons := strings.Split(departure_options_string, ";")
		if err != nil {
			return err
		}

		if flowData.Message.Content == "+" {
			for index := 5; index < len(departure_optinons); index++ {
				formattedMessage += fmt.Sprintf("*%d - %s* \n", index+1, strings.SplitN(departure_optinons[index], "-", 2)[1])
			}
			messageResponse = fmt.Sprintf("%s \n Ou, se quiser, digite:\n*• V para voltar.*", formattedMessage)
			err = r.whatsappRepository.SimpleMessage(flowData.Ctx, messageResponse, flowData.Customer, flowData.MetaTokens)
			return nil
		} else {
			// We need to check if the client message is a number
			numberArray, err := utils.ExtractAndConvertToInt(flowData.Message.Content)
			if err != nil {
				err = fmt.Errorf("Não identifiquei a opção selecionada, tente novamente.")
				return err
			}

			departureID := strings.Split(departure_optinons[numberArray-1], "-")[0]
			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "departure_select", "string", departureID); err != nil {
				return err
			}

			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "departure_name", "string", strings.SplitN(departure_optinons[numberArray-1], "-", 2)[1]); err != nil {
				return err
			}

			// Update field "sale_stage"
			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "sale_stage", "string", "3"); err != nil {
				return err
			}
		}
	}

	return nil
}

// List the options compatible with the city and save the chosen one
func (r *Workflow) arrivalStage(flowData dto.FlowData, stage string) error {
	var formattedMessage string

	//Get departure Options
	departureName, err := utils.OtherFields(flowData.Session.OtherFields, "departure_name")
	if err != nil {
		return err
	}

	//Get Arrival Options
	arrivalName, err := utils.OtherFields(flowData.Session.OtherFields, "arrival_name")
	if stage == "3" {
		var messageResponse string
		var formattedMessage string

		fmt.Println("Dentro de arrivalstage ")
		responseFunc, err := r.guanabaraRepository.GetCity(flowData.Ctx, arrivalName)
		if err != nil {
			return err
		}

		if len(responseFunc) > 0 {
			var button []models.Button
			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "arrival_options", "string", strings.Join(responseFunc, ";")); err != nil {
				return err
			}

			var controle int
			// Itere pelo array de strings e construa a mensagem formatada (até o quinto elemento)
			for index, messageArray := range responseFunc {
				if index >= 5 {
					break // Saia do loop depois do quinto elemento
				}
				controle++

				formattedMessage += fmt.Sprintf("*%d* - %s", index+1, strings.SplitN(messageArray, "-", 2)[1])

				if index < 4 {
					formattedMessage += "\n" // Adicione uma quebra de linha, exceto para o quinto elemento
				}
			}

			var complement string
			if len(responseFunc) > 5 {
				qtd := len(responseFunc) - controle
				complement = fmt.Sprintf("\n •+ para ver mais opções (foram encontradas mais %d opções)", qtd)
				button = []models.Button{
					{ID: "+", Title: fmt.Sprintf("Listar (%d opções)", qtd), NextNode: "node1"},
				}
			}

			messageResponse = fmt.Sprintf("Você vai sair de %s.📍\n  \n Selecione o seu *destino*: 📍 \n%s \n\n*Confirma pra mim?* Só digitar o número acima. \n \n Ou, se quiser, digite: %s \n*• V para voltar.*", departureName, formattedMessage, complement)
			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "sale_stage", "string", "4"); err != nil {
				return err
			}

			err = r.whatsappRepository.InteractiveMessage(flowData.Ctx, messageResponse, button, flowData.Customer, flowData.MetaTokens)
			fmt.Println(messageResponse)
		} else {
			flowData.Message.Content = "destino_nao_encontrada"
			err = r.errorAssitantReturn(flowData)
			return nil
		}
	}

	//Saves the selected arrival id
	if stage == "4" {
		arrival_options_string, err := utils.OtherFields(flowData.Session.OtherFields, "arrival_options")
		arrival_optinons := strings.Split(arrival_options_string, ";")
		if err != nil {
			return err
		}

		if flowData.Message.Content == "+" {
			for index := 4; index < len(arrival_optinons); index++ {
				formattedMessage += fmt.Sprintf("*%d* - %s \n", index+1, strings.SplitN(arrival_optinons[index], "-", 2)[1])
			}

			err = r.whatsappRepository.SimpleMessage(flowData.Ctx, formattedMessage, flowData.Customer, flowData.MetaTokens)
			return nil
		} else {
			// Precisamos verificar se a mensagem do cliente é um numero
			numberArray, err := utils.ExtractAndConvertToInt(flowData.Message.Content)
			if err != nil {
				err = fmt.Errorf("Não identifiquei a opção selecionada, tente novamente.")
				return err
			}

			arrivalID := strings.Split(arrival_optinons[numberArray-1], "-")[0]
			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "arrival_select", "string", arrivalID); err != nil {
				return err
			}

			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "arrival_name", "string", strings.SplitN(arrival_optinons[numberArray-1], "-", 2)[1]); err != nil {
				return err
			}

			// Atualizar campo "sale_stage"
			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "sale_stage", "string", "5"); err != nil {
				return err
			}
		}
	}

	return nil
}

// Validates that the routes are valid and searches the timetables
func (r *Workflow) validationSchedules(flowData dto.FlowData, stage string) error {
	departureSelect, err := utils.OtherFields(flowData.Session.OtherFields, "departure_select")
	arrivalSelect, err := utils.OtherFields(flowData.Session.OtherFields, "arrival_select")
	dataDeparture, err := utils.OtherFields(flowData.Session.OtherFields, "departure_data")
	departureName, err := utils.OtherFields(flowData.Session.OtherFields, "departure_name")
	arrivalName, err := utils.OtherFields(flowData.Session.OtherFields, "arrival_name")
	shedules_options, err := utils.OtherFields(flowData.Session.OtherFields, "shedules_options")
	if err != nil {
		return err
	}

	if stage == "5" {
		if strings.Contains(departureSelect, ":group") || strings.Contains(arrivalSelect, ":group") {
			fmt.Println("A string contém ':group'.")
		} else {
			fmt.Print("Route Validation: ")
			responseFunc, err := r.guanabaraRepository.RouteValidation(flowData.Ctx, strings.TrimSpace(departureSelect), strings.TrimSpace(arrivalSelect))
			if err != nil {
				fmt.Println("erro route validation: ", err)
				return err
			}

			if !responseFunc {
				flowData.Message.Content = "rota_nao_encontrada" // Comando que vai ser enviado ao assistant no retorno
				err := r.errorAssitantReturn(flowData)
				if err != nil {
					return err
				}
				//err = r.whatsappRepository.SimpleMessage(flowData.Ctx, "Atualmente não atuamos na rota escolhida, vamos tentar novamente?", flowData.Customer, flowData.MetaTokens)
			}
		}

		shedules, err := r.guanabaraRepository.GetSchedules(flowData.Ctx, departureSelect, arrivalSelect, dataDeparture)
		if err != nil {
			return err
		}

		if len(shedules) > 0 {
			var controle int

			//first message
			firstMessage := fmt.Sprintf("Pronto, encontrei *%d* opções de %s para %s no dia %s. \n\nOrdenei por preço, mas você pode escolher mudar a ordem para se guiar por *horário* ou *classe*. Clique no botão para selecionar um horário de partida e a classe do serviço!Caso não goste de nenhuma das opções, envie *V* para retornar a etapa de busca", len(shedules), departureName, arrivalName, dataDeparture)
			err = r.whatsappRepository.SimpleMessage(flowData.Ctx, firstMessage, flowData.Customer, flowData.MetaTokens)

			// convert to string
			shedulesString, err := json.Marshal(shedules)
			if err != nil {
				return err
			}

			// save shedules options
			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "shedules_options", "string", string(shedulesString)); err != nil {
				return err
			}

			// checks the need for filters and sends it to the client
			var filterClass []string
			var filterHourAM int
			var filterHourPM int

			for _, trip := range shedules {
				//filter class
				utils.AddStringIfNotExists("filterClass:"+trip.ClassOfServiceName, &filterClass)

				//filter hour
				//datetime := strings.Split(trip.DepartureTime, ":")[0]
				num, err := strconv.Atoi(strings.Split(trip.DepartureTime, ":")[0])
				if err != nil {
					return nil
				}
				if num < 12 {
					filterHourAM++
				} else {
					//if datetime.After(midday)
					filterHourPM++
				}
			}

			// Send list with filters class
			if len(filterClass) > 0 {
				var row []models.Row
				for _, class := range filterClass {
					var option models.Row
					option.ID = class
					option.Title = strings.Split(class, ":")[1]
					row = append(row, option)
				}
				err = r.whatsappRepository.InteractiveMessageList(flowData.Ctx, flowData.Customer, flowData.MetaTokens, "*Filtrar por Classe* 💺", row)
			}

			//Send list with filter hour
			var rowHour []models.Row
			if filterHourAM > 0 {
				var option models.Row
				option.ID = "filterHour:AM"
				option.Title = "00:01 - 12:00"
				rowHour = append(rowHour, option)
			}
			if filterHourPM > 0 {
				var option models.Row
				option.ID = "filterHour:PM"
				option.Title = "12:01 - 00:00"
				rowHour = append(rowHour, option)
			}
			if len(rowHour) > 0 {
				err = r.whatsappRepository.InteractiveMessageList(flowData.Ctx, flowData.Customer, flowData.MetaTokens, "*Filtrar por Horário* 🕐", rowHour)
			}

			// list top 5 of array
			for index, messageArray := range shedules {
				if index >= 5 {
					break
				}
				controle++
				format := fmt.Sprintf("*%d*. %s \n 📍 %s - %s \n 🏁 %s - %s \n %s", index+1, messageArray.CompanyName, messageArray.DepartureTime, messageArray.DepartureLocation.Name, messageArray.ArrivalTime, messageArray.ArrivalLocation.Name, messageArray.ClassOfServiceName)

				button := []models.Button{
					{ID: messageArray.ControlNumber, Title: "SELECIONAR", NextNode: ""},
				}
				err = r.whatsappRepository.InteractiveMessage(flowData.Ctx, format, button, flowData.Customer, flowData.MetaTokens)
			}

			if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "sale_stage", "string", "6"); err != nil {
				return err
			}

			//Final message for this step, allows you to go back or see more options
			buttons := []models.Button{}

			if len(shedules) > 5 {
				qtd := len(shedules) - controle
				moreButton := models.Button{
					ID:       "+",
					Title:    fmt.Sprintf("Ver mais %d opções", qtd),
					NextNode: "",
				}

				buttons = append(buttons, moreButton)
			}

			message := "Escolha uma das opções acima \n\n Ou selecione uma das opções abaixo:"
			buttons = append(buttons, models.Button{ID: "v", Title: "Voltar", NextNode: "node1"})

			err = r.whatsappRepository.InteractiveMessage(flowData.Ctx, message, buttons, flowData.Customer, flowData.MetaTokens)
		} else {
			flowData.Message.Content = "rota_nao_encontrada"
			err := r.errorAssitantReturn(flowData)
			if err != nil {
				return err
			}
		}
	} else if stage == "6" {

		// na minha cabeca faz sentido termos dois campos de viagem, as consultadas e uma que seria as filtradas, basicamente quando o cliente filtrar alguma opção nas atualizamos a segunda e sempre olhamos ela caso o cliente envie um numero na resposta
		// quando ele clicar em um botao sempre vamos olhar o id do control number que esta no botão

		if flowData.Message.Content == "+" {
			var arrayShedules []models.TripInfoResponse
			err := json.Unmarshal([]byte(shedules_options), &arrayShedules)
			if err != nil {
				return err
			}

			for i := 4; i < len(arrayShedules); i++ {
				fmt.Println("dentro do array de shedules")
				format := fmt.Sprintf("*%d*. %s \n 📍 %s - %s \n 🏁 %s - %s \n %s", i+1, arrayShedules[i].CompanyName, arrayShedules[i].DepartureTime, arrayShedules[i].DepartureLocation.Name, arrayShedules[i].ArrivalTime, arrayShedules[i].ArrivalLocation.Name, arrayShedules[i].ClassOfServiceName)

				button := []models.Button{
					{ID: arrayShedules[i].ControlNumber, Title: "SELECIONAR", NextNode: ""},
				}
				err = r.whatsappRepository.InteractiveMessage(flowData.Ctx, format, button, flowData.Customer, flowData.MetaTokens)
			}
		}

		// identifies the filter
		if strings.Contains(flowData.Message.Content, "filterClass:") {
			filter := strings.Split(flowData.Message.Content, ":")[1]
			fmt.Println("Filtro da pesquisa ", filter)

			var arrayShedules []models.TripInfoResponse
			err := json.Unmarshal([]byte(shedules_options), &arrayShedules)
			if err != nil {
				return err
			}

			var arrayFilter []models.TripInfoResponse
			for _, trip := range arrayShedules {
				if trip.ClassOfServiceName == filter {
					arrayFilter = append(arrayFilter, trip)
				}
			}

			for index, messageArray := range arrayFilter {
				format := fmt.Sprintf("*%d*. %s \n 📍 %s - %s \n 🏁 %s - %s \n %s", index+1, messageArray.CompanyName, messageArray.DepartureTime, messageArray.DepartureLocation.Name, messageArray.ArrivalTime, messageArray.ArrivalLocation.Name, messageArray.ClassOfServiceName)

				button := []models.Button{
					{ID: messageArray.ControlNumber, Title: "SELECIONAR", NextNode: ""},
				}
				err = r.whatsappRepository.InteractiveMessage(flowData.Ctx, format, button, flowData.Customer, flowData.MetaTokens)
			}

		}

		// identifies the filter Hour
		if strings.Contains(flowData.Message.Content, "filterHour:") {
			filter := strings.Split(flowData.Message.Content, ":")[1]

		}

		return nil
	}

	return nil
}

// Returns to the gpt flow passing a message
func (r *Workflow) errorAssitantReturn(flowData dto.FlowData) error {
	if err := r.updateSessionField(flowData.Ctx, flowData.Session.ID.Hex(), "current_stage", "string", "assistant"); err != nil {
		return err
	}

	err, threadsIds := r.gpt_assistant(flowData)
	if err != nil {
		//_ = r.whatsappRepository.SimpleMessage(ctx, "Desculpe, parece que tivemos um problema técnico. Vou reiniciar nosso processo para garantir que tudo funcione corretamente. Por favor, me informe novamente, qual informação você está buscando? Se precisar de assistência em algum tópico específico, estou aqui para ajudar!", customer, metaTokens)
		_ = r.whatsappRepository.SimpleMessage(flowData.Ctx, err.Error(), flowData.Customer, flowData.MetaTokens)

		if threadsIds.ThreadId != "" && threadsIds.RunId != "" {
			_, err = r.openaiRepository.CancelRun(flowData.Ctx, threadsIds.ThreadId, threadsIds.RunId)
		}

		return err
	}

	return nil
}

// Send message for customer with menu
func (r *Workflow) sendMenu(flowData dto.FlowData, arguments string) error {
	var tripRequest models.TripRequest
	var nodeMenu models.Node
	var messageStart = "Oi! Tudo bem? \n \nEu sou a Ana, a assistente virtual do Grupo Guanabara! Tô aqui pra te ajudar a viver experiências incríveis de viagem. 🤗 \n\n"
	var messageMenu = "Sobre o que vamos falar hoje? *Você quer conferir ou comprar bilhetes de viagem? Já comprou e precisa de ajuda? Ou quer saber mais sobre algum assunto?* É só me contar!"
	var finalMessage string
	if err := json.Unmarshal([]byte(arguments), &tripRequest); err != nil {
		fmt.Println("Erro ao fazer o parsing do JSON:", err)
		return err
	}

	if tripRequest.FirstMessage {
		finalMessage = messageStart + messageMenu
	} else {
		finalMessage = messageMenu
	}

	button1 := models.Button{ID: "pesquisar_passagem", Title: "COMPRAR PASSAGEM", NextNode: "node1"}
	button2 := models.Button{ID: "meus_bilhetes", Title: "MEUS BILHETES", NextNode: "node2"}
	button3 := models.Button{ID: "outros_assuntos", Title: "OUTROS ASSUNTOS", NextNode: "node3"}
	nodeMenu.Parameters.Buttons = []models.Button{button1, button2, button3}

	r.whatsappRepository.InteractiveMessage(flowData.Ctx, finalMessage, []models.Button{button1, button2, button3}, flowData.Customer, flowData.MetaTokens)
	return nil
}
