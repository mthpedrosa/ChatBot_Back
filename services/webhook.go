package services

import (
	"autflow_back/repositories"
	"autflow_back/utils"
)

/*type MetaIds struct {
	IdTelefone   string
	TokenConexao string
}*/

type Webhook struct {
	metaRepository     *repositories.Metas
	workflowRepository *repositories.Workflows
	logger             utils.Logger
}

func NewWebhook(metaRepository *repositories.Metas, workflowRepository *repositories.Workflows, logger utils.Logger) *Webhook {
	return &Webhook{
		metaRepository:     metaRepository,
		workflowRepository: workflowRepository,
		logger:             logger,
	}
}

/*func (r *Webhook) WebhookRun(ctx context.Context, payload models.WebhookPayload, webhookId string) error {

	//return
	var payload models.WebhookPayload

	// Extrair o ID do webhook da solicitação
	parametros := mux.Vars(c.Request())
	webhookId := parametros["webhookId"]

	// Decodifica o corpo da solicitação (payload JSON) em 'payload'
	if err := json.NewDecoder(c.Request().Body).Decode(&payload); err != nil {
		return responses.Erro(c, http.StatusBadRequest, errors.New("Erro ao decodificar o payload do webhook"))
	}

	//ignora retorno de status, porque não usamos artualmente - vamos usar nas mensagens
	if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 && len(payload.Entry[0].Changes[0].Value.Statuses) > 0 {
		return errors.New("")
	}

	// Check account_meta
	//repoMeta := repositories.NewMetaRepository(r.db)
	meta, erro := repoMeta.Find(ctx, "webhook="+webhookId)
	if erro != nil {
		return responses.Erro(c, http.StatusInternalServerError, erro)
	}

	workflow, err := r.Workflow.IdentifyWorkflow(models.WebhookPayload(payload), meta[0])
	if err != nil {
		return responses.Erro(c, http.StatusBadRequest, err)
	}

	err = r.Workflow.RunWorkflow(models.WebhookPayload(payload), meta[0], workflow)
	if err != nil {
		return responses.Erro(c, http.StatusBadRequest, err)
	}

	//return
	/*strin, erro := services.PercorreWorkflow(*workflow, models.WebhookPayload(payload), *meta)
	if erro != nil {
		http.Error(w, erro.Error(), http.StatusBadRequest)
		fmt.Print(strin)
	}

	return nil
}
*/
