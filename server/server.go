package server

import (
	"autflow_back/repositories"
	"autflow_back/server/controllers"
	"autflow_back/server/routes"
	"autflow_back/services"
	"autflow_back/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	e           *echo.Echo
	mongoClient *mongo.Client
	redis       *redis.Client
	logger      utils.Logger
}

func NewServer(mongoClient *mongo.Client, redis *redis.Client, logger utils.Logger, e *echo.Echo) *Server {
	return &Server{
		e:           e,
		mongoClient: mongoClient,
		redis:       redis,
		logger:      logger,
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
	guanabaraRepository := repositories.NewGuanabaraRepository(s.redis)
	openaiRepository := repositories.NewOpenAiRepository(s.redis)
	whatsappRepository := repositories.NewWhatsappRepository(s.redis)

	// Services
	workflowService := services.NewWorkflow(workflowRepository, metaRepository, customerRepository, sessionRepository, conversationsRepository, s.logger, guanabaraRepository, openaiRepository, whatsappRepository)
	accountMetaService := services.NewMeta(metaRepository, s.logger)
	conversationService := services.NewConversation(conversationsRepository, s.logger)
	customerService := services.NewCustomer(workflowRepository, customerRepository, s.logger)
	userService := services.NewUser(userRepository, s.logger)
	loginService := services.NewLogin(userRepository, s.logger)
	sessionService := services.NewSession(sessionRepository, s.logger)

	//Controllers
	metaController := controllers.NewMetaController(accountMetaService)
	conversationsController := controllers.NewConversationController(conversationService)
	workflowController := controllers.NewWorkflowsController(workflowService)
	customerController := controllers.NewCustomerController(customerService)
	userController := controllers.NewUserController(userService)
	loginController := controllers.NewLoginController(loginService)
	webhookController := controllers.NewWebhookController(workflowService, accountMetaService)
	sessionController := controllers.NewSessionController(sessionService)

	//Start routes
	routes.RegisterMetaRoutes(s.e, metaController)
	routes.RegisterConversationsRoutes(s.e, conversationsController)
	routes.RegisterWorkflowsRoutes(s.e, workflowController)
	routes.RegisterCustomerRoutes(s.e, customerController)
	routes.RegisterUsersRoutes(s.e, userController)
	routes.RegisterLoginRoutes(s.e, loginController)
	routes.RegisterWebhookRoutes(s.e, webhookController)
	routes.RegisterSessionRoutes(s.e, sessionController)

	// Middlewares
	s.e.Use(middleware.CORS())
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.RequestID())
	s.e.Use(middleware.Gzip())
	s.e.Use(middleware.Secure())

	// Start server
	err := s.e.Start(":" + viper.GetString("PORT"))
	if err != nil {
		return err
	}

	return nil
}
