package utils

import (
	"autflow_back/models"
	"strconv"
	"time"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GetExtension(contentType string) string {
	switch contentType {
	case "audio/ogg":
		return "ogg"
	case "audio/mpeg":
		return "mp3"
	default:
		return ""
	}
}

// removeDiacritics remove os acentos de uma string
func RemoveDiacritics(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}))
	result, _, _ := transform.String(t, s)
	return result
}

func OtherFields(otherFields []models.Fields, nameField string) (string, error) {
	var value string

	for _, fields := range otherFields {
		if fields.Name == nameField {
			value = fields.Value
		}
	}

	return value, nil
}

func FormatDate(inputDate string) (string, error) {
	parsedDate, err := time.Parse("2006-01-02", inputDate)
	if err != nil {
		return "", err
	}

	formattedDate := parsedDate.Format("02/01/2006")
	return formattedDate, nil
}

func ExtractAndConvertToInt(str string) (int, error) {
	var onlyNumbersString string
	for _, char := range str {
		if unicode.IsDigit(char) {
			onlyNumbersString += string(char)
		}
	}

	return strconv.Atoi(onlyNumbersString)
}

func AddStringIfNotExists(str string, slice *[]string) {
	for _, v := range *slice {
		if v == str {
			return
		}
	}

	*slice = append(*slice, str)
}

/*
// refatorar para não duplicar o campo
func SalvarDados(campo models.Campo, cliente models.Cliente, conversa models.Conversa) (string, error) {
	db := database.GetDatabase()
	var idCampo string

	switch campo.Entidade {
	case "customer":
		{
			repo := repositories.NovoRepositorioCustomers(db)
			idCampo, err := repo.AdicionarCampo(campo.Tipo, campo.CampoNome, campo.Conteudo, cliente.ID.Hex())
			if err != nil {
				return "", err
			}
			return idCampo, nil
		}
	default:
		{
		}
	}

	return idCampo, nil
}

func PercorrePassosInternos(workflow models.Workflow, passo models.WorkflowPasso, cliente models.Cliente, meta models.MetaIds, success bool, conversaId string) (string, error) {

	//Identifica se deu erro e qual passo a seguir
	var idProxPasso string
	if success == true {
		fmt.Print(passo.SaidaSuccess)
		idProxPasso = passo.SaidaSuccess
	} else {
		fmt.Print(passo.SaidaError)
		idProxPasso = passo.SaidaError
	}

	//idenfica o proximo passo
	var proximoPasso models.WorkflowPasso
	for _, passo := range workflow.Passos {
		if passo.ID == idProxPasso {
			proximoPasso = passo
		}
	}

	switch proximoPasso.Tipo {
	case "mensagem_simples":
		fmt.Print("Mensagem Simples")
		err := enviarMensagemSimples(proximoPasso, cliente, meta)
		if err != nil {
			return "", nil
		}
	case "mensagem_interativa":
		fmt.Print("Mensagem Interativa")
		_, err := enviarMensagemInterativa(proximoPasso, cliente, meta)
		if err != nil {
			return "", nil
		}
	case "salvar_dados":
		if passo.Tipo == "chatgpt_text" {
			proximoPasso.Campo.Conteudo = passo.Comando
		} else if passo.Tipo == "chatgpt_audio_text" {
			proximoPasso.Campo.Conteudo = passo.Comando
		}
		_, err := SalvarDados(proximoPasso.Campo, cliente, models.Conversa{})
		if err != nil {
			//_, err = PercorrePassosInternos(worflow, prox, *cliente, ok, false)
			return "", err
		}

		db := database.GetDatabase()
		repo := repositories.NovoRepositorioDeConversations(db)
		err = repo.AtualizaUltimoPasso(proximoPasso.ID, conversaId)
		if err != nil {
			return "", err
		}

		str, err := PercorrePassosInternos(workflow, proximoPasso, cliente, meta, success, conversaId)
		if err != nil {
			return str, err
		}
		return "", nil
	case "chatgpt_text":
		proximoPasso.Comando = proximoPasso.Comando + " " + passo.Comando
		retornoChat, err := GetChatGPTResponse(proximoPasso.Comando)
		if err != nil {
			return "", err
		}
		proximoPasso.Comando = retornoChat

		str, err := PercorrePassosInternos(workflow, proximoPasso, cliente, meta, success, conversaId)
		if err != nil {
			return str, err
		}
		return "", nil
	case "gipsyy_venda":
		err := NoGipssyCompra(cliente, conversaId, meta, proximoPasso, workflow)
		fmt.Print(err)
		if err != nil {
			return "", err
		}

		return "", nil

	default:
		{
		}
	}

	db := database.GetDatabase()
	repo := repositories.NovoRepositorioDeConversations(db)
	err := repo.AtualizaUltimoPasso(proximoPasso.ID, conversaId)
	if err != nil {
		return "", err
	}

	return proximoPasso.ID, nil
}

///////  UTEIS  ///////

func removerQuebrasDeLinha(texto string) string {
	return strings.Replace(texto, "\n", "", -1)
}

func removerEspacosEmBranco(texto string) string {
	return strings.ReplaceAll(texto, " ", "")
}

func removerEspacosEmBrancoInicioFim(texto string) string {
	return strings.TrimSpace(texto)
}



func converterParaBase64(nomeArquivo string) (string, error) {
	// Ler o arquivo de áudio
	arquivo, err := ioutil.ReadFile("arquivos_temporarios/" + nomeArquivo)
	if err != nil {
		return "", err
	}

	// Converter o conteúdo para base64
	base64String := base64.StdEncoding.EncodeToString(arquivo)
	return base64String, nil
}

func VerificarFormatoData(str string) bool {
	// Expressão regular para verificar o formato DD/MM
	regex := `^(0[1-9]|[12][0-9]|3[01])/(0[1-9]|1[0-2])$`

	// Compila a expressão regular
	r := regexp.MustCompile(regex)

	// Verifica se a string corresponde ao padrão de data
	return r.MatchString(str)
}

func converterDataFormato(data string) (string, error) {
	// Obtém o ano atual
	anoAtual := time.Now().Year()

	// Formata a string atual para o formato desejado (YYYY)
	dataCompleta := fmt.Sprintf("%d-%s", anoAtual, data)

	// Faz o parsing da string para validar a data
	dataParseada, err := time.Parse("2006-01-02", dataCompleta)
	if err != nil {
		return "", err
	}

	// Retorna a data formatada
	return dataParseada.Format("2006-01-02"), nil
}

func obterDataHoje() (string, string, string, string, error) {
	// Obtém a data atual
	horaAtual := time.Now()

	// Formata a data de hoje no formato DD/MM
	dataFormatoBR := horaAtual.Format("02/01")

	// Formata a data de hoje no formato YYYY-MM-DD
	dataFormatoISO := horaAtual.Format("2006-01-02")

	// Obtém a data de amanhã
	amanha := horaAtual.AddDate(0, 0, 1)

	// Formata a data de amanhã no formato DD/MM
	dataAmanhaBR := amanha.Format("02/01")

	// Formata a data de amanhã no formato YYYY-MM-DD
	dataAmanhaISO := amanha.Format("2006-01-02")

	return dataFormatoBR, dataFormatoISO, dataAmanhaBR, dataAmanhaISO, nil
}

func removerCaracteresInvalidos(texto string) string {
	// Expressão regular para encontrar caracteres inválidos
	reg := regexp.MustCompile("[^\x00-\x7F]+")

	// Remove caracteres inválidos do texto usando a expressão regular
	textoLimpo := reg.ReplaceAllString(texto, "")

	return textoLimpo
}
*/
