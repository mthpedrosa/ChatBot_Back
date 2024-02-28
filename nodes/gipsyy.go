package nodes

/*
import (
	"autflow_back/src/database"
	"autflow_back/models"
	"autflow_back/src/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Location struct {
	ID          int     `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Latitude    *string `json:"latitude,omitempty"`
	Longitude   *string `json:"longitude,omitempty"`
}

type Price struct {
	OriginalPriceValue    float64 `json:"original_price_value"`
	PriceValue            float64 `json:"price_value"`
	PassengerType         int     `json:"passenger_type"`
	PassengerTypeName     string  `json:"passenger_type_name"`
	AvailableSeats        int     `json:"avaliable_seats"`
	PromoCodeError        bool    `json:"promo_code_error"`
	PromoCodeMessage      *string `json:"promo_code_message,omitempty"`
	PriceRouteConditionID int     `json:"price_route_condition_id"`
	PricingName           string  `json:"pricing_name"`
}

type ConnectionRoute struct {
	ControlNumber      string   `json:"control_number"`
	DepartureTime      string   `json:"departure_time"`
	ArrivalTime        string   `json:"arrival_time"`
	ArrivalDay         int      `json:"arrival_day"`
	IDCompany          int      `json:"id_company"`
	CompanyName        string   `json:"company_name"`
	ClassOfServiceName string   `json:"class_of_service_name"`
	BPE                bool     `json:"bpe"`
	Duration           string   `json:"duration"`
	DepartureLocation  Location `json:"departure_location"`
	ArrivalLocation    Location `json:"arrival_location"`
}

type Trip struct {
	ControlNumber              string            `json:"control_number"`
	DepartureTime              string            `json:"departure_time"`
	DepartureLocation          Location          `json:"departure_location"`
	ArrivalTime                string            `json:"arrival_time"`
	ArrivalLocation            Location          `json:"arrival_location"`
	ArrivalDay                 int               `json:"arrival_day"`
	IDCompany                  int               `json:"id_company"`
	CompanyName                string            `json:"company_name"`
	PriceValue                 float64           `json:"price_value"`
	Prices                     []Price           `json:"prices"`
	RoundTripPrices            []interface{}     `json:"roundTripPrices"`
	ClassOfServiceName         string            `json:"class_of_service_name"`
	BPE                        bool              `json:"bpe"`
	HasConnection              bool              `json:"has_connection"`
	AvailableSeats             int               `json:"available_seats"`
	OriginalAvailability       int               `json:"original_availability"`
	ConnectionRoutes           []ConnectionRoute `json:"connection_routes"`
	CurrencyCode               *string           `json:"currency_code,omitempty"`
	Duration                   string            `json:"duration"`
	LogoCompany                *interface{}      `json:"logo_company,omitempty"`
	OptionalInsurance          float64           `json:"optional_insurance"`
	ReceiptType                string            `json:"receipt_type"`
	IDDailySchedule            int               `json:"id_daily_schedule"`
	IDSchedule                 int               `json:"id_schedule"`
	ScheduleControlNumber      string            `json:"schedule_control_number"`
	ControlNumberDailySchedule string            `json:"control_number_daily_schedule"`
	ContractModel              string            `json:"contractModel"`
	IsHighlightTrip            bool              `json:"is_highlight_trip"`
}

type Trecho struct {
	ExternalID               int     `json:"externalId"`
	Name                     string  `json:"name"`
	Latitude                 float64 `json:"latitude"`
	Longitude                float64 `json:"longitude"`
	TopOrigin                bool    `json:"topOrigin"`
	TopDestination           bool    `json:"topDestination"`
	Price                    float64 `json:"price"`
	CountryCode              string  `json:"countryCode"`
	Code                     string  `json:"code"`
	StateCode                string  `json:"stateCode"`
	CityName                 string  `json:"cityName"`
	LocationsGroupExternalID int     `json:"locationsGroupExternalId"`
	Description              string  `json:"description"`
	IsGroup                  bool    `json:"isGroup"`
	ID                       int     `json:"id"`
}

type TrechoDestino struct {
	ExternalID     int     `json:"external_id"`
	Name           string  `json:"name"`
	Latitude       string  `json:"latitude"`
	Longitude      string  `json:"longitude"`
	TopOrigin      bool    `json:"top_origin"`
	TopDestination bool    `json:"top_destination"`
	Description    string  `json:"description"`
	Price          *int    `json:"price,omitempty"`
	CountryCode    string  `json:"country_code"`
	Code           string  `json:"code"`
	StateCode      string  `json:"state_code"`
	CityName       string  `json:"city_name"`
	Distance       *int    `json:"distance,omitempty"`
	GroupName      *string `json:"group_name,omitempty"`
}

type APIResponse struct {
	Data []Trecho `json:"data"`
}

type CamposDestino struct {
	Origem    string `json:"origem"`
	Destino   string `json:"destino"`
	DataIda   string `json:"dataIda"`
	DataVolta string `json:"dataVolta"`
}

func buscarTrechos(trechos []Trecho, nomeDesejado string) []Trecho {
	var trechosEncontrados []Trecho

	// Itera sobre os trechos
	for _, trecho := range trechos {
		// Verifica se o nome do trecho contém a string desejada (ignorando diferenças de maiúsculas/minúsculas)
		if strings.Contains(removeDiacritics(strings.ToLower(trecho.Name)), removeDiacritics(strings.ToLower(nomeDesejado))) {
			// Adiciona o trecho encontrado ao array de trechos encontrados
			trechosEncontrados = append(trechosEncontrados, trecho)
		}
	}

	return trechosEncontrados
}

func buscarTrechosDestino(trechos []TrechoDestino, nomeDesejado string) []TrechoDestino {
	var trechosEncontrados []TrechoDestino

	// Itera sobre os trechos
	for _, trecho := range trechos {
		// Verifica se o nome do trecho contém a string desejada (ignorando diferenças de maiúsculas/minúsculas)
		if strings.Contains(strings.ToLower(trecho.Name), strings.ToLower(nomeDesejado)) {
			// Adiciona o trecho encontrado ao array de trechos encontrados
			trechosEncontrados = append(trechosEncontrados, trecho)
		}
	}

	return trechosEncontrados
}

// Pega todas as opções de origem da Gipsyy
func GetTrechos() ([]Trecho, error) {
	url := "https://gds.gipsyy.com.br/api/Locations"

	// Realiza uma solicitação GET à URL especificada
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Verifica se a resposta foi bem-sucedida (código de status 200)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao buscar os trechos, código de status: %d", response.StatusCode)
	}

	// Decodifica a resposta JSON em uma estrutura APIResponse
	var apiResp APIResponse
	err = json.NewDecoder(response.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}

	return apiResp.Data, nil
}

// Busca destinos compativeis com a origem
func GetTrechosDestinos(origem string) ([]TrechoDestino, error) {
	//https://gds.gipsyy.com.br/api/GipsyyWeb/Rails/GetAllRouteArrivalByOrigins/49495

	url := "https://gds.gipsyy.com.br/api/GipsyyWeb/Rails/GetAllRouteArrivalByOrigins/" + origem

	// Realiza uma solicitação GET à URL especificada
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Verifica se a resposta foi bem-sucedida (código de status 200)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao buscar os trechos, código de status: %d", response.StatusCode)
	}

	// Decodifica a resposta JSON em uma estrutura APIResponse
	var apiResp []TrechoDestino
	err = json.NewDecoder(response.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}

// Busca opções de viagem
func GetTrips(dataIda, origemId, destinoId string) ([]Trip, error) {
	fmt.Println("Buscando trips")
	//Origem > departure_location
	//Destino > arrival_location
	//Data ida > departure_date
	//Data Volta > arrival_date
	// Round trip > ida e volta ?
	// Canal de venda > salesChannel

	url := "https://gds.gipsyy.com.br/api/GipsyyWeb/Rails/Trips?arrival_location=" + destinoId + "&departure_location=" + origemId + "&departure_date=" + dataIda + "&salesChannel=1&round_trip=false"

	fmt.Println(url)
	// Realiza uma solicitação GET à URL especificada
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Verifica se a resposta foi bem-sucedida (código de status 200)
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro ao buscar os trechos, código de status: %d", response.StatusCode)
	}

	// Decodifica a resposta JSON em uma estrutura APIResponse
	var apiResp []Trip
	err = json.NewDecoder(response.Body).Decode(&apiResp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}

func GetTripInfo(dataIda, origemId, destinoId string) ([]models.TripInfo, error) {

	// Cria a requisição POST
	url := `https://backend.gipsyy.com.br/api/trips?`
	for _, campo := range []string{
		"locale=pt-BR",
		"format=json",
		"departure_location_id=" + origemId,
		"arrival_location_id=" + destinoId,
		"departure_at=" + dataIda,
	} {
		url += campo + "&"
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	// Realiza a requisição HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Converte o corpo da resposta em string e exibe
	var apiResp []models.TripInfo
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&apiResp)
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}

func NoGipssyCompra(cliente models.Customer, conversa string, meta models.MetaIds, proxPasso models.Node, workflow models.Workflow) error {
	db := database.GetDatabase()
	repo := repositories.NewCustomersRepository(db)
	clienteAtual, err := repo.FindId(cliente.ID.Hex())
	if err != nil {
		return err
	}

	fmt.Println("Dentro da func gipsyy compra , o passo atual é:", proxPasso.ID)
	/*if proxPasso.Funcao == "origem_destino" {
		_, err := OrigemDestino(*clienteAtual, meta, proxPasso.Comando, conversa, proxPasso)
		if err != nil {
			return err
		}
	}

	switch proxPasso.Funcao {
	case "origem":
		_, err := Origem(*clienteAtual, meta, proxPasso.Comando, conversa, proxPasso, workflow)
		if err != nil {
			return err
		}

	case "destino":
		_, err := Destino(*clienteAtual, meta, proxPasso.Comando, conversa, proxPasso, workflow)
		if err != nil {
			return err
		}

	case "dataIda":
		_, err := DataIda(*clienteAtual, meta, proxPasso.Comando, conversa, proxPasso, workflow)
		if err != nil {
			return err
		}

	case "buscaPassagem":
		_, err := buscaPassagem(*clienteAtual, meta, proxPasso.Comando, conversa, proxPasso, workflow)
		if err != nil {
			return err
		}

	case "buscaAssentos":
		fmt.Println("Chegando nos assentos")
	}

	/*repos := repositories.NovoRepositorioDeConversations(db)
	err = repos.AtualizaUltimoPasso(proximoPassoId, conversa)
	if err != nil {
		return err
	}

	return nil
}

/*func verificaCampo(valor string) int {
	//Verifica campo origem e onde estamos no processo
	if valor == "" {
		return 0
	} else {
		// Separar a string pelo hífen
		origemSplit := strings.Split(valor, "-")
		if _, err := strconv.Atoi(origemSplit[0]); err == nil {
			return 2
			// ja possuimos o id
		} else {
			fmt.Println("A primeira parte NÃO é um número.")
			return 1
			// temos um texto não o id
		}
	}
}

// função que identifica a origem que o cliente quer - precisamos adicionar tratativas caso não ache opções
func Origem(cliente models.Cliente, meta models.MetaIds, valor string, conversaId string, passoAtual models.WorkflowPasso, workflow models.Workflow) (bool, error) {
	db := database.GetDatabase()
	repos := repositories.NovoRepositorioDeConversations(db)
	// -- Pegando campos de usuario -- //
	fmt.Println("Dentro da func gipsyy origem , o passo atual é:", passoAtual.ID)

	err := repos.AtualizaUltimoPasso(passoAtual.ID, conversaId)
	if err != nil {
		return false, err
	}

	var jsonData CamposDestino
	var origem string
	var destino string
	var opcoesOrigem string
	var opcoesDestino string

	for _, campo := range cliente.Campos {
		if campo.ID == "customer_origem_destino" {
			// Decodifica a string JSON para uma estrutura ou mapa
			if err := json.Unmarshal([]byte(campo.Value), &jsonData); err != nil {
				fmt.Print("quebrando aqui 1")
				return false, err
			}

		} else if campo.ID == "customer_origem" {
			origem = campo.Value
		} else if campo.ID == "customer_destino" {
			destino = campo.Value
		} else if campo.ID == "customer_opcoes_origem" {
			opcoesOrigem = campo.Value
		} else if campo.ID == "customer_opcoes_destino" {
			opcoesDestino = campo.Value
		}

	}

	// -- FIM -- //

	// -- Pegando opções -- //
	trechos, err := GetTrechos()
	fmt.Println(trechos)

	fmt.Println("Campo Origem: " + origem)
	fmt.Println("Campo Destino" + destino)
	fmt.Println("Campo Opcoes de Origem : " + opcoesOrigem)
	fmt.Println("Campo Opcoes de Destino : " + opcoesDestino)
	fmt.Print(jsonData)

	if jsonData.Origem == "" {
		fmt.Println("Caso não tenha nem uma origem e nem destino")

		// Envia opções pro cliente escolher
		var proximoPasso models.WorkflowPasso
		proximoPasso.Conteudo = "Não consegui identificar a origem e destino da sua viagem, vamos começar de novo? \n 📍 *De onde você quer sair?* \n Você pode enviar assim: 'Quero viajar de Fortaleza para Sobral amanhã a noite e voltar dia 20 de leito.' \n Digite ou envie um áudio que eu te escuto! ☺️ \n Ou, se quiser, *digite V para voltar.* "
		err = enviarMensagemSimples(proximoPasso, *&cliente, meta)
		if err != nil {
			return false, err
		}

		err = repos.AtualizaUltimoPasso("002", conversaId)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if origem == "" && jsonData.Origem != "" && opcoesOrigem == "" {
		//cai aqui quando o cliente falou a origem mas não escolheu as opçoes
		opcoes := buscarTrechos(trechos, jsonData.Origem)
		if err != nil {
			return false, err
		}

		nomeOpcoes := ""
		nomeOpcoesInterno := ""
		contador := 1

		for _, opcao := range opcoes {
			nomeOpcoes += "*" + fmt.Sprint(contador) + " - " + opcao.Name + "* \n "
			nomeOpcoesInterno += fmt.Sprint(contador) + "-" + fmt.Sprint(opcao.ExternalID) + ";"
			contador++
		}

		// Envia opções pro cliente escolher
		var proximoPasso models.WorkflowPasso
		proximoPasso.Conteudo = "Selecione a sua origem: 📍 \n " + nomeOpcoes + " \n Confirma pra mim? Só digitar o número acima. \n Ou, se quiser, digite: \n • V para voltar."
		err = enviarMensagemSimples(proximoPasso, *&cliente, meta)
		if err != nil {
			return false, err
		}

		var campo models.Campo
		campo.Tipo = "string"
		campo.CampoNome = "opcoes_origem"
		campo.Conteudo = nomeOpcoesInterno
		campo.Entidade = "customer"

		_, err = SalvarDados(campo, cliente, models.Conversa{})
		if err != nil {
			return false, err
		}

		return true, nil
	} else if origem == "" && opcoesOrigem != "" && valor != "" {
		var IdOpcao string
		origemSplit := strings.Split(opcoesOrigem, ";")
		fmt.Print(origemSplit)

		for i := 0; i < len(origemSplit); i++ {
			Id := strings.Split(origemSplit[i], "-")

			if Id[0] == valor {
				fmt.Print("ID da opcao escolhida ------------------------------------- ")
				fmt.Print(Id)
				fmt.Print("ID da opcao escolhida ------------------------------------- ")
				IdOpcao = Id[1]
			}
		}

		var campo models.Campo
		campo.Tipo = "string"
		campo.CampoNome = "origem"
		campo.Conteudo = IdOpcao
		campo.Entidade = "customer"

		_, err = SalvarDados(campo, cliente, models.Conversa{})
		if err != nil {
			return false, err
		}

	}

	// pega campos atualizados
	repo := repositories.NewCustomersRepository(db)
	clienteAtual, err := repo.BuscarId(cliente.ID.Hex())
	if err != nil {
		return false, err
	}
	for _, campo := range clienteAtual.Campos {
		if campo.ID == "customer_origem_destino" {
			// Decodifica a string JSON para uma estrutura ou mapa
			if err := json.Unmarshal([]byte(campo.Value), &jsonData); err != nil {
				return false, err
			}

		} else if campo.ID == "customer_origem" {
			origem = campo.Value
		} else if campo.ID == "customer_destino" {
			destino = campo.Value
		} else if campo.ID == "customer_opcoes_origem" {
			opcoesOrigem = campo.Value
		} else if campo.ID == "customer_opcoes_destino" {
			opcoesDestino = campo.Value
		}

	}

	//Defini o proximo passo caso aqui esteja tudo certo
	if origem != "" {
		fmt.Println("origem ok, defino o proximo passo")
		str, err := PercorrePassosInternos(workflow, passoAtual, cliente, meta, true, conversaId)
		if err != nil {
			return true, err
		}
		fmt.Print(str)

	}
	return true, nil
}

// função que motra os destinos com base na origem e defini oque o cliente quer  - precisamos adicionar tratativas caso não ache opções
func Destino(cliente models.Cliente, meta models.MetaIds, valor string, conversaId string, passoAtual models.WorkflowPasso, workflow models.Workflow) (bool, error) {
	db := database.GetDatabase()
	repos := repositories.NovoRepositorioDeConversations(db)
	// -- Pegando campos de usuario -- //
	fmt.Println("Dentro da func gipsyy destino , o passo atual é:", passoAtual.ID)

	err := repos.AtualizaUltimoPasso(passoAtual.ID, conversaId)
	if err != nil {
		return false, err
	}

	var jsonData CamposDestino
	var origem string
	var destino string
	var opcoesOrigem string
	var opcoesDestino string

	for _, campo := range cliente.Campos {
		if campo.ID == "customer_origem_destino" {
			// Decodifica a string JSON para uma estrutura ou mapa
			if err := json.Unmarshal([]byte(campo.Value), &jsonData); err != nil {
				return false, err
			}

		} else if campo.ID == "customer_origem" {
			origem = campo.Value
		} else if campo.ID == "customer_destino" {
			destino = campo.Value
		} else if campo.ID == "customer_opcoes_origem" {
			opcoesOrigem = campo.Value
		} else if campo.ID == "customer_opcoes_destino" {
			opcoesDestino = campo.Value
		}

	}

	// -- FIM -- //

	// -- Pegando opções -- //
	//trechos, err := GetTrechos()

	fmt.Println("Campo Origem: " + origem)
	fmt.Println("Campo Destino" + destino)
	fmt.Println("Campo Opcoes de Origem : " + opcoesOrigem)
	fmt.Println("Campo Opcoes de Destino : " + opcoesDestino)
	fmt.Print(jsonData)

	if destino == "" && jsonData.Destino != "" && opcoesDestino == "" && origem != "" {
		//cai aqui quando o cliente falou a destino mas não escolheu as opçoes
		fmt.Println("Entrando no destino")
		trechos, err := GetTrechosDestinos(origem) //pega os destino que fazem rota com a origem
		if err != nil {
			return false, err
		}

		opcoes := buscarTrechosDestino(trechos, jsonData.Destino)
		if err != nil {
			return false, err
		}

		nomeOpcoes := ""
		nomeOpcoesInterno := ""
		contador := 1

		for _, opcao := range opcoes {
			nomeOpcoes += "*" + fmt.Sprint(contador) + " - " + opcao.Name + "* \n "
			nomeOpcoesInterno += fmt.Sprint(contador) + "-" + fmt.Sprint(opcao.ExternalID) + ";"
			contador++
		}

		// Envia opções pro cliente escolher
		var proximoPasso models.WorkflowPasso
		proximoPasso.Conteudo = "Selecione seu destino: 📍 \n " + nomeOpcoes + " \n Confirma pra mim? Só digitar o número acima. \n Ou, se quiser, digite: \n • V para voltar."
		err = enviarMensagemSimples(proximoPasso, *&cliente, meta)
		if err != nil {
			return false, err
		}

		/* Você vai sair de RIO DE JANEIRO (NOVO RIO) - RJ.📍

		Selecione o seu destino: 🏁

		1. BELO HORIZONTE - MG
		2. BELO HORIZONTE (SAVASSI) - MG
		3. BELO HORIZONTE (TERMINAL JK) - MG
		4. BELO HORIZONTE (GARAGEM CAICARAS) - MG

		Confirma pra mim? Só digitar o número acima.

		Ou, se quiser, digite:
		• Outro local;
		• V para voltar.


		var campo models.Campo
		campo.Tipo = "string"
		campo.CampoNome = "opcoes_destino"
		campo.Conteudo = nomeOpcoesInterno
		campo.Entidade = "customer"

		_, err = SalvarDados(campo, cliente, models.Conversa{})
		if err != nil {
			return false, err
		}
	} else if destino == "" && opcoesDestino != "" && valor != "" {
		var IdOpcao string
		origemSplit := strings.Split(opcoesDestino, ";")
		fmt.Print(origemSplit)

		for i := 0; i < len(origemSplit); i++ {
			Id := strings.Split(origemSplit[i], "-")

			if Id[0] == valor {
				fmt.Print("ID da opcao escolhida ------------------------------------- ")
				fmt.Print(Id)
				fmt.Print("ID da opcao escolhida ------------------------------------- ")
				IdOpcao = Id[1]
			}
		}

		var campo models.Campo
		campo.Tipo = "string"
		campo.CampoNome = "destino"
		campo.Conteudo = IdOpcao
		campo.Entidade = "customer"

		_, err = SalvarDados(campo, cliente, models.Conversa{})
		if err != nil {
			return false, err
		}
	}

	// pega campos atualizados
	repo := repositories.NewCustomersRepository(db)
	clienteAtual, err := repo.BuscarId(cliente.ID.Hex())
	if err != nil {
		return false, err
	}
	for _, campo := range clienteAtual.Campos {
		if campo.ID == "customer_origem_destino" {
			// Decodifica a string JSON para uma estrutura ou mapa
			if err := json.Unmarshal([]byte(campo.Value), &jsonData); err != nil {
				return false, err
			}

		} else if campo.ID == "customer_origem" {
			origem = campo.Value
		} else if campo.ID == "customer_destino" {
			destino = campo.Value
		} else if campo.ID == "customer_opcoes_origem" {
			opcoesOrigem = campo.Value
		} else if campo.ID == "customer_opcoes_destino" {
			opcoesDestino = campo.Value
		}

	}

	//Defini o proximo passo caso aqui esteja tudo certo
	if destino != "" {
		fmt.Println("destino ok, defino o proximo passo")
		str, err := PercorrePassosInternos(workflow, passoAtual, cliente, meta, true, conversaId)
		if err != nil {
			return true, err
		}
		fmt.Print(str)

	}
	return true, nil
}

// verifica data de ida
func DataIda(cliente models.Cliente, meta models.MetaIds, valor string, conversaId string, passoAtual models.WorkflowPasso, workflow models.Workflow) (bool, error) {
	db := database.GetDatabase()
	repos := repositories.NovoRepositorioDeConversations(db)
	// -- Pegando campos de usuario -- //
	fmt.Println("Dentro da func gipsyy Data , o passo atual é:", passoAtual.ID)

	err := repos.AtualizaUltimoPasso(passoAtual.ID, conversaId)
	if err != nil {
		return false, err
	}

	var jsonData CamposDestino
	var dataIda string

	for _, campo := range cliente.Campos {
		if campo.ID == "customer_origem_destino" {
			// Decodifica a string JSON para uma estrutura ou mapa
			if err := json.Unmarshal([]byte(campo.Value), &jsonData); err != nil {
				return false, err
			}

		} else if campo.ID == "customer_opcoes_dataIda" {
			dataIda = campo.Value
		}
	}

	fmt.Println(dataIda)
	fmt.Print(jsonData)
	ok := VerificarFormatoData(jsonData.DataIda)

	if jsonData.DataIda != "" && ok {
		/*ok := VerificarFormatoData(jsonData.DataIda)// não faz sentido porque a data já deve estar coerente
		if ok {
			//caso a data fornecida no audio ou escrita seja coerente, perguntamos dando a opção de botão, por que se não for vamos pedir para escrever
			var proximoPasso models.WorkflowPasso
			proximoPasso.Conteudo = "Selecione seu destino: 📍 \n \n Confirma pra mim? Só digitar o número acima. \n Ou, se quiser, digite: \n • V para voltar."
			err = enviarMensagemSimples(proximoPasso, *&cliente, meta)
			if err != nil {
				return false, err
			}
		}

		dataFormatada, err := converterDataFormato(jsonData.DataIda)

		var campo models.Campo
		campo.Tipo = "string"
		campo.CampoNome = "opcoes_dataIda"
		campo.Conteudo = dataFormatada
		campo.Entidade = "customer"

		_, err = SalvarDados(campo, cliente, models.Conversa{})
		if err != nil {
			return false, err
		}

	} else if jsonData.DataIda == "" && valor == "" {

		dataHojeBR, dataHojeISO, dataAmanhaBR, dataAmanhaISO, err := obterDataHoje()
		if err != nil {
			fmt.Println("Erro ao obter a data:", err)
			return false, err
		}
		// Envia opções pro cliente escolher
		var proximoPasso models.WorkflowPasso
		proximoPasso.Conteudo = "Maravilha! Qual é a data de ida? 🗓️ \n Você pode digitar a data (DD/MM) ou selecionar uma das sugestões abaixo."
		proximoPasso.Botoes = []models.Botao{
			{
				ID:     dataHojeISO,
				Titulo: "Hoje : " + dataHojeBR,
			},
			{
				ID:     dataAmanhaISO,
				Titulo: "Amanhã : " + dataAmanhaBR,
			},
		}
		_, err = enviarMensagemInterativa(proximoPasso, *&cliente, meta)
		if err != nil {
			return false, err
		}

	} else if jsonData.DataIda == "" && valor != "" {
		var valorFinal string
		ok := VerificarFormatoData(valor)

		if ok {
			valorFinal, err = converterDataFormato(valor)
			if err != nil {
				return false, err
			}
		} else if !ok {
			valorFinal = valor
		}

		var campo models.Campo
		campo.Tipo = "string"
		campo.CampoNome = "opcoes_dataIda"
		campo.Conteudo = valorFinal
		campo.Entidade = "customer"

		_, err = SalvarDados(campo, cliente, models.Conversa{})
		if err != nil {
			return false, err
		}

	}

	// pega campos atualizados
	repo := repositories.NewCustomersRepository(db)
	clienteAtual, err := repo.BuscarId(cliente.ID.Hex())
	if err != nil {
		return false, err
	}
	for _, campo := range clienteAtual.Campos {
		if campo.ID == "customer_origem_destino" {
			// Decodifica a string JSON para uma estrutura ou mapa
			if err := json.Unmarshal([]byte(campo.Value), &jsonData); err != nil {
				return false, err
			}
		} else if campo.ID == "customer_opcoes_dataIda" {
			dataIda = campo.Value
		}
	}

	//Defini o proximo passo caso aqui esteja tudo certo
	if dataIda != "" {
		fmt.Println("DataIda ok, defino o proximo passo")
		str, err := PercorrePassosInternos(workflow, passoAtual, cliente, meta, true, conversaId)
		if err != nil {
			return true, err
		}
		fmt.Print(str)

	}
	return true, nil
}

func buscaPassagem(cliente models.Cliente, meta models.MetaIds, valor string, conversaId string, passoAtual models.WorkflowPasso, workflow models.Workflow) (bool, error) {
	db := database.GetDatabase()
	repos := repositories.NovoRepositorioDeConversations(db)
	// -- Pegando campos de usuario -- //
	fmt.Println("Dentro da func gipsyy Busca passsagem , o passo atual é:", passoAtual.ID)

	err := repos.AtualizaUltimoPasso(passoAtual.ID, conversaId)
	if err != nil {
		return false, err
	}

	var jsonData CamposDestino
	var dataIda string
	var origem string
	var destino string
	var viagemEscolhida string

	for _, campo := range cliente.Campos {
		if campo.ID == "customer_origem_destino" {
			// Decodifica a string JSON para uma estrutura ou mapa
			if err := json.Unmarshal([]byte(campo.Value), &jsonData); err != nil {
				return false, err
			}

		} else if campo.ID == "customer_opcoes_dataIda" {
			dataIda = campo.Value
		} else if campo.ID == "customer_origem" {
			origem = campo.Value
		} else if campo.ID == "customer_destino" {
			destino = campo.Value
		} else if campo.ID == "customer_viagemEscolhida" {
			viagemEscolhida = campo.Value
		}

	}

	fmt.Println(dataIda)
	fmt.Println(origem)
	fmt.Println(destino)

	if dataIda != "" && origem != "" && destino != "" && viagemEscolhida == "" && valor == "" {
		fmt.Println("entra antes de buscar viagens")
		trips, err := GetTrips(dataIda, origem, destino)
		if err != nil {
			return false, err
		}

		if len(trips) == 0 {
			fmt.Println("Não encontrei passagem")

			// Envia opções pro cliente escolher
			var proximoPasso models.WorkflowPasso
			proximoPasso.Conteudo = "Não consegui encontrar passagens para a origem x destino e data selecionados, vamos tentar novamente? "
			err = enviarMensagemSimples(proximoPasso, *&cliente, meta)
			if err != nil {
				return false, err
			}

			campos := []string{"origem", "destino", "opcoes_dataIda"}
			err = deleteInfoVenda(campos, cliente)

			str, err := PercorrePassosInternos(workflow, passoAtual, cliente, meta, false, conversaId)
			if err != nil {
				return true, err
			}
			fmt.Print(str)

		}

		var contador = 1
		//percorremos as passagens e enviamos as opçoes
		for _, trip := range trips {
			fmt.Println("entra no for de mensagens")

			texto := fmt.Sprintf("%d. %s \n 📍 %s:%s - "+trip.DepartureLocation.Name+" \n 🏁 %s:%s - "+trip.ArrivalLocation.Name+" \n %s - R$ %.2f 💸 ", contador, trip.CompanyName, strings.Split(trip.DepartureTime, ":")[0], strings.Split(trip.DepartureTime, ":")[1], strings.Split(trip.ArrivalTime, ":")[0], strings.Split(trip.ArrivalTime, ":")[1], strings.Split(trip.ClassOfServiceName, ";")[0], trip.PriceValue)
			fmt.Println("passou pelo texto formatado")

			// Envia opções pro cliente escolher
			var proximoPasso models.WorkflowPasso
			proximoPasso.Conteudo = texto
			proximoPasso.Botoes = []models.Botao{
				{
					ID:     trip.ControlNumber,
					Titulo: "Selecionar",
				},
			}
			_, err = enviarMensagemInterativa(proximoPasso, *&cliente, meta)
			if err != nil {
				return false, err
			}

			contador++
		}
	}

	if valor != "" && viagemEscolhida == "" {
		var viagemNova string

		trips, err := GetTrips(dataIda, origem, destino)
		if err != nil {
			return false, err
		}

		for _, trip := range trips {
			if trip.ControlNumber == valor {
				viagemNovas, err := json.Marshal(trip)
				if err != nil {
					panic(err)
				}
				fmt.Sprintln(viagemNova)

				viagemNova = string(viagemNovas)
			}
		}

		fmt.Sprintln(viagemNova)

		var campo models.Campo
		campo.Tipo = "string"
		campo.CampoNome = "viagemEscolhida"
		campo.Conteudo = viagemNova
		campo.Entidade = "customer"

		_, err = SalvarDados(campo, cliente, models.Conversa{})
		if err != nil {
			return false, err
		}

		str, err := PercorrePassosInternos(workflow, passoAtual, cliente, meta, true, conversaId)
		if err != nil {
			return true, err
		}

		fmt.Print(str)
		return true, nil

	}

	return true, nil
}

/*func buscaAssentos(cliente models.Cliente, meta models.MetaIds, valor string, conversaId string, passoAtual models.WorkflowPasso, workflow models.Workflow) (bool, error) {
	db := database.GetDatabase()
	repos := repositories.NovoRepositorioDeConversations(db)
	// -- Pegando campos de usuario -- //
	fmt.Println("Dentro da func gipsyy Busca assento , o passo atual é:", passoAtual.ID)
	fmt.Println("Cria checkout id")

	err := repos.AtualizaUltimoPasso(passoAtual.ID, conversaId)
	if err != nil {
		return false, err
	}

	var jsonData CamposDestino
	var origem string
	var destino string
	var opcoesOrigem string
	var opcoesDestino string

	for _, campo := range cliente.Campos {
		if campo.ID == "customer_origem_destino" {
			// Decodifica a string JSON para uma estrutura ou mapa
			if err := json.Unmarshal([]byte(campo.Value), &jsonData); err != nil {
				return false, err
			}

		} else if campo.ID == "customer_origem" {
			origem = campo.Value
		} else if campo.ID == "customer_destino" {
			destino = campo.Value
		} else if campo.ID == "customer_opcoes_origem" {
			opcoesOrigem = campo.Value
		} else if campo.ID == "customer_opcoes_destino" {
			opcoesDestino = campo.Value
		} else if campo.ID == "customer_opcoes_destino" {
			opcoesDestino = campo.Value
		}

	}

	return true, nil
}

func deleteInfoVenda(campos []string, cliente models.Cliente) error {
	for _, nomecampo := range campos {
		var campo models.Campo
		campo.Tipo = "string"
		campo.CampoNome = nomecampo
		campo.Conteudo = ""
		campo.Entidade = "customer"

		_, err := SalvarDados(campo, cliente, models.Conversa{})
		if err != nil {
			return err
		}
	}

	return nil
}
*/
