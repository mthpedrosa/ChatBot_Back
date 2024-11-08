package main

import (
	"autflow_back/server"
	"autflow_back/src/config"
	"autflow_back/utils"
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	config.Load()

	// Carregar as variáveis do arquivo .env
	// envFile := ".env"
	// if os.Getenv("ENV") == "local" {
	// 	envFile = ".env.local"
	// }

	// err := godotenv.Load(envFile)
	// if err != nil {
	// 	log.Fatalf("Error loading %s file: %v", envFile, err)
	// }

	// Configura o Viper para carregar variáveis do ambiente
	viper.AutomaticEnv()

	log.Println("Starting microservice")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a new Logger
	logger := utils.NewApiLogger(&utils.Config{
		Encoding: "json",
		Env:      "dev",
	})
	logger.InitLogger("info")

	// Create MongoDB Client
	clientOptions := options.Client().ApplyURI(viper.GetString("MONGO_URL"))

	// Connect to MongoDB
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Panicf("Error creating mongoDB: %s", err)
	}

	// Check the connection
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Panicf("Error pinging mongoDB: %s", err)
	}

	logger.Info("Connected to MongoDB!")

	// Redis
	/*redisClient := adapters.CreateRedisClient(ctx, adapters.RedisConfig{
		DB:  0,
		URL: viper.GetString("REDIS_URL"),
	})
	if redisClient == nil {
		logger.Panicf("Error creating redis client: %s", err)
	}

	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Panicf("Error pinging redis client: %s", err)
	}

	logger.Info("Connected to Redis!")*/

	// echo
	e := echo.New()

	app := server.NewServer(mongoClient, logger, e)
	err = app.Start()
	if err != nil {
		logger.Panicf("Error starting server: %s", err)
	}
}
