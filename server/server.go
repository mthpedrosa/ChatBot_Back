package server

import (
	"autflow_back/repositories"
	"autflow_back/server/controllers"
	"autflow_back/server/routes"
	"autflow_back/services"
	"autflow_back/utils"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	e           *echo.Echo
	mongoClient *mongo.Client
	//redis       *redis.Client
	logger utils.Logger
}

func NewServer(mongoClient *mongo.Client, logger utils.Logger, e *echo.Echo) *Server {
	return &Server{
		e:           e,
		mongoClient: mongoClient,
		//redis:       redis,
		logger: logger,
	}
}

func (s *Server) Start() error {
	// Repositories
	metaRepository := repositories.NewMetaRepository(s.mongoClient)
	conversationsRepository := repositories.NewConversationsRepository(s.mongoClient)
	customerRepository := repositories.NewCustomersRepository(s.mongoClient)
	sessionRepository := repositories.NewSessionsRepository(s.mongoClient)
	userRepository := repositories.NewUsersRepository(s.mongoClient)
	workflowRepository := repositories.NewWorkflowsRepository(s.mongoClient)
	//guanabaraRepository := repositories.NewGuanabaraRepository(s.redis)
	openaiRepository := repositories.NewOpenAiRepository()
	openaiMongoRepository := repositories.NewOpenAiMongoRepository(s.mongoClient)
	whatsappRepository := repositories.NewWhatsappRepository()
	userPlanRepository := repositories.NewUserPlanRepository(s.mongoClient)

	// Services
	//workflowService := services.NewWorkflow(workflowRepository, metaRepository, customerRepository, sessionRepository, conversationsRepository, s.logger, openaiRepository, whatsappRepository)
	userPlanService := services.NewUserPlanService(userPlanRepository)
	messageHandlerService := services.NewMessageHandler(metaRepository, customerRepository, sessionRepository, conversationsRepository, s.logger, openaiRepository, whatsappRepository, userPlanRepository, userPlanService)
	accountMetaService := services.NewMeta(metaRepository, s.logger)
	conversationService := services.NewConversation(conversationsRepository, s.logger)
	customerService := services.NewCustomer(workflowRepository, customerRepository, s.logger)
	userService := services.NewUser(userRepository, s.logger)
	loginService := services.NewLogin(userRepository, s.logger)
	sessionService := services.NewSession(sessionRepository, s.logger)
	openaiService := services.NewOpenAi(openaiRepository, openaiMongoRepository, s.logger)
	reportsService := services.NewReports(metaRepository, customerRepository, sessionRepository, conversationsRepository, s.logger, openaiRepository, whatsappRepository)

	//Controllers
	metaController := controllers.NewMetaController(accountMetaService)
	conversationsController := controllers.NewConversationController(conversationService)
	//workflowController := controllers.NewWorkflowsController(workflowService)
	customerController := controllers.NewCustomerController(customerService)
	userController := controllers.NewUserController(userService)
	loginController := controllers.NewLoginController(loginService)
	webhookController := controllers.NewWebhookController(messageHandlerService, accountMetaService, userPlanService)
	sessionController := controllers.NewSessionController(sessionService)
	openaiController := controllers.NewOpenAiController(openaiService)
	reportsController := controllers.NewReportsController(reportsService)
	userPlanController := controllers.NewUserPlanController(userPlanService)

	//Start routes
	routes.RegisterMetaRoutes(s.e, metaController)
	routes.RegisterConversationsRoutes(s.e, conversationsController)
	//routes.RegisterWorkflowsRoutes(s.e, workflowController)
	routes.RegisterCustomerRoutes(s.e, customerController)
	routes.RegisterUsersRoutes(s.e, userController)
	routes.RegisterLoginRoutes(s.e, loginController)
	routes.RegisterWebhookRoutes(s.e, webhookController)
	routes.RegisterSessionRoutes(s.e, sessionController)
	routes.RegisterOpenAiRoutes(s.e, openaiController)
	routes.RegisterReportsRoutes(s.e, reportsController)
	routes.RegisterUserPlanRoutes(s.e, userPlanController)

	// Middlewares
	s.e.Use(middleware.CORS())
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.RequestID())
	s.e.Use(middleware.Gzip())
	s.e.Use(middleware.Secure())

	// Inicia o servidor Echo diretamente na porta desejada
	log.Println("Iniciando o servidor Echo na porta 8080...")
	if err := s.e.Start(":8080"); err != nil {
		return err
	}

	return nil
}
