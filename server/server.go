package server

import (
	"autflow_back/repositories"
	"autflow_back/server/controllers"
	"autflow_back/server/routes"
	"autflow_back/services"
	"autflow_back/utils"
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

type Server struct {
	e           *echo.Echo
	mongoClient *mongo.Client
	redis       *redis.Client
	logger      utils.Logger
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

	// Services
	//workflowService := services.NewWorkflow(workflowRepository, metaRepository, customerRepository, sessionRepository, conversationsRepository, s.logger, openaiRepository, whatsappRepository)
	messageHandlerService := services.NewMessageHandler(metaRepository, customerRepository, sessionRepository, conversationsRepository, s.logger, openaiRepository, whatsappRepository)
	accountMetaService := services.NewMeta(metaRepository, s.logger)
	conversationService := services.NewConversation(conversationsRepository, s.logger)
	customerService := services.NewCustomer(workflowRepository, customerRepository, s.logger)
	userService := services.NewUser(userRepository, s.logger)
	loginService := services.NewLogin(userRepository, s.logger)
	sessionService := services.NewSession(sessionRepository, s.logger)
	openaiService := services.NewOpenAi(openaiRepository, openaiMongoRepository, s.logger)

	//Controllers
	metaController := controllers.NewMetaController(accountMetaService)
	conversationsController := controllers.NewConversationController(conversationService)
	//workflowController := controllers.NewWorkflowsController(workflowService)
	customerController := controllers.NewCustomerController(customerService)
	userController := controllers.NewUserController(userService)
	loginController := controllers.NewLoginController(loginService)
	webhookController := controllers.NewWebhookController(messageHandlerService, accountMetaService)
	sessionController := controllers.NewSessionController(sessionService)
	openaiController := controllers.NewOpenAiController(openaiService)

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

	// Middlewares
	s.e.Use(middleware.CORS())
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
	s.e.Use(middleware.RequestID())
	s.e.Use(middleware.Gzip())
	s.e.Use(middleware.Secure())

	// Configuração do Ngrok
	ctx := context.Background()
	listener, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(
			config.WithDomain("talented-starling-quietly.ngrok-free.app"),
		),
		ngrok.WithAuthtokenFromEnv(),
	)
	if err != nil {
		return err
	}

	// Inicia o servidor Echo usando o listener do Ngrok
	s.e.Listener = listener

	log.Println("Ngrok tunnel established at:", listener.URL())

	// Agora inicie o servidor Echo usando o listener do Ngrok
	if err := s.e.StartServer(s.e.Server); err != nil {
		return err
	}

	return nil
}
