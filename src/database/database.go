package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

/*func OpenConnection() {
	// Crie um contexto com timeout para a conexão
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Obtenha a string de conexão do arquivo de configuração (certifique-se de que sua função StringConexaoBanco retorne a string correta)
	mongoURI := config.DatabaseConnectionString

	// Configura as opções do cliente
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Crie o cliente MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Verifique se a conexão está ativa
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// Se a conexão estiver ativa, atribua o banco de dados
	db = client.Database("autflow")
}*/

// Função para obter a instância do banco de dados
func GetDatabase() *mongo.Database {
	return db
}
