package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	Port            = 0
	SecretKey       []byte
	ChatGPTAPI      string
	PermissionsUser []string
)

// Load environment variables.
func Load() {
	var erro error

	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	Port, erro = strconv.Atoi(os.Getenv("API_PORT"))
	if erro != nil {
		Port = 9000
	}

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
	ChatGPTAPI = os.Getenv("CHATGPT_KEY")
	PermissionsUser = strings.Split(os.Getenv("PERMISSIONS_USERS"), ":")
}
