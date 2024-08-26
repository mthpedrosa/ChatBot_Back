package controllers

// import (
// 	"autflow_back/models/dto"
// 	"autflow_back/requests"
// 	"autflow_back/services"
// 	"autflow_back/src/authentication"
// 	"autflow_back/src/responses"
// 	"errors"
// 	"net/http"
// 	"strings"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/labstack/echo/v4"
// )

// type Workflow struct {
// 	workflowService *services.Workflow
// }

// func NewWorkflowsController(workflowService *services.Workflow) *Workflow {
// 	return &Workflow{
// 		workflowService: workflowService,
// 	}
// }

// func (r *Workflow) Insert(c echo.Context) error {
// 	// Check request body using Bind
// 	createWorkflowRequest := new(requests.CreateWorkflowRequest)

// 	if err := c.Bind(createWorkflowRequest); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{
// 			"message": err.Error(),
// 		})
// 	}

// 	if err := validate.Struct(createWorkflowRequest); err != nil {
// 		validationErrors := err.(validator.ValidationErrors)

// 		errorsMessages := []string{}
// 		for _, err := range validationErrors {
// 			errorsMessages = append(errorsMessages, err.Error())
// 		}
// 		return c.JSON(http.StatusBadRequest, echo.Map{
// 			"message": "Validation errors",
// 			"errors":  errorsMessages,
// 		})
// 	}

// 	// ID of the user creating the user
// 	creatorUser, erro := authentication.ExtractIdToken(c.Request())
// 	if erro != nil {
// 		return responses.Erro(c, http.StatusBadRequest, erro)
// 	}

// 	dt := &dto.CreateWorkflowDTO{
// 		PhoneMetaId: createWorkflowRequest.PhoneMetaId,
// 		Name:        createWorkflowRequest.Name,
// 		Description: createWorkflowRequest.Description,
// 		Nodes:       createWorkflowRequest.Nodes,
// 		Active:      createWorkflowRequest.Active,
// 		FirstNode:   createWorkflowRequest.FirstNode,
// 		LastNode:    createWorkflowRequest.LastNode,
// 		CreatedBy:   creatorUser,
// 	}

// 	if erro = dt.Prepare("cadastro"); erro != nil {
// 		return responses.Erro(c, http.StatusBadRequest, erro)
// 	}

// 	idCriado, erro := r.workflowService.Create(c.Request().Context(), dt)
// 	if erro != nil {
// 		return responses.Erro(c, http.StatusBadRequest, erro)
// 	}

// 	return responses.JSON(c, http.StatusCreated, idCriado)
// }

// func (r *Workflow) Find(c echo.Context) error {
// 	urlParts := strings.Split(c.Request().URL.String(), "?")
// 	var query string
// 	if len(urlParts) > 1 {
// 		query = urlParts[1]
// 	}

// 	workflows, erro := r.workflowService.Find(c.Request().Context(), query)
// 	if erro != nil {
// 		return responses.Erro(c, http.StatusInternalServerError, erro)
// 	}

// 	// Map goals to a MetaDTO list
// 	workflowsDTO := make([]dto.WorkflowListDTO, len(workflows))
// 	for i, workflow := range workflows {
// 		workflowsDTO[i] = dto.WorkflowListDTO{
// 			ID:          workflow.ID,
// 			Name:        workflow.Name,
// 			Description: workflow.Description,
// 			Active:      workflow.Active,
// 			PhoneMetaId: workflow.PhoneMetaId,
// 		}
// 	}

// 	return responses.JSON(c, http.StatusOK, workflowsDTO)
// }

// func (r *Workflow) FindId(c echo.Context) error {
// 	id := c.Param("id")
// 	if id == "" {
// 		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
// 	}

// 	workflow, erro := r.workflowService.FindId(c.Request().Context(), id)
// 	if erro != nil {
// 		return responses.Erro(c, http.StatusInternalServerError, erro)
// 	}

// 	workflowDTO := dto.WorkflowDetailDTO{
// 		ID:          workflow.ID,
// 		PhoneMetaId: workflow.PhoneMetaId,
// 		Name:        workflow.Name,
// 		Description: workflow.Description,
// 		Nodes:       workflow.Nodes,
// 		CreatedAt:   workflow.CreatedAt,
// 		UpdateAt:    workflow.UpdateAt,
// 		CreatedBy:   workflow.CreatedBy,
// 		Active:      workflow.Active,
// 		FirstNode:   workflow.FirstNode,
// 		LastNode:    workflow.LastNode,
// 	}

// 	return responses.JSON(c, http.StatusOK, workflowDTO)
// }

// func (r *Workflow) Edit(c echo.Context) error {
// 	// Check request body using Bind
// 	createWorkflowRequest := new(requests.CreateWorkflowRequest)

// 	if err := c.Bind(createWorkflowRequest); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{
// 			"message": err.Error(),
// 		})
// 	}

// 	if err := validate.Struct(createWorkflowRequest); err != nil {
// 		validationErrors := err.(validator.ValidationErrors)

// 		errorsMessages := []string{}
// 		for _, err := range validationErrors {
// 			errorsMessages = append(errorsMessages, err.Error())
// 		}
// 		return c.JSON(http.StatusBadRequest, echo.Map{
// 			"message": "Validation errors",
// 			"errors":  errorsMessages,
// 		})
// 	}

// 	dt := &dto.CreateWorkflowDTO{
// 		PhoneMetaId: createWorkflowRequest.PhoneMetaId,
// 		Name:        createWorkflowRequest.Name,
// 		Description: createWorkflowRequest.Description,
// 		Nodes:       createWorkflowRequest.Nodes,
// 		Active:      createWorkflowRequest.Active,
// 		FirstNode:   createWorkflowRequest.FirstNode,
// 		LastNode:    createWorkflowRequest.LastNode,
// 	}

// 	// Extract the ID from the request
// 	id := c.Param("id")
// 	if id == "" {
// 		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
// 	}

// 	erro := r.workflowService.Edit(c.Request().Context(), id, dt)
// 	if erro != nil {
// 		return responses.Erro(c, http.StatusInternalServerError, erro)

// 	}

// 	return responses.JSON(c, http.StatusOK, "Workflow editada com sucesso")
// }

// func (r *Workflow) Delete(c echo.Context) error {

// 	// Extract the ID from the request
// 	id := c.Param("id")
// 	if id == "" {
// 		return responses.Erro(c, http.StatusInternalServerError, errors.New("Parâmetro não fornecido na solicitação"))
// 	}

// 	// call the function delete
// 	erro := r.workflowService.Delete(c.Request().Context(), id)
// 	if erro != nil {
// 		return responses.Erro(c, http.StatusInternalServerError, erro)

// 	}

// 	return responses.JSON(c, http.StatusOK, "Workflow deletado com sucesso")
// }
