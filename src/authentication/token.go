package authentication

import (
	"autflow_back/models"
	"autflow_back/src/config"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	//jwt "github.com/dgrijalva/jwt-go"
)

// Returns signed token with user permissions - Retorna token assinado com as permissoes do usuario
func CreateToken(usuario models.User) (string, error) {
	permissoes := jwt.MapClaims{}
	permissoes["authorized"] = true
	permissoes["exp"] = time.Now().Add(time.Hour * 6).Unix() // expiration time
	permissoes["userId"] = usuario.ID
	permissoes["profile"] = usuario.Profile
	// secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissoes)
	return token.SignedString([]byte(config.SecretKey)) //secret
}

// Checks if the token is valid
func ValidateToken(r *http.Request) error {
	tokenString := extractToken(r)
	token, erro := jwt.Parse(tokenString, returnVerificationKey)
	if erro != nil {
		return erro
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("Token Invalido")
}

// Function to extract a token.
func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}

	return ""
}

func returnVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Método de assinatura inesperado! %v", token.Header["alg"])
	}

	return config.SecretKey, nil
}

func ExtractPermissions(r *http.Request) (string, error) {
	tokenString := extractToken(r)
	token, erro := jwt.Parse(tokenString, returnVerificationKey)
	if erro != nil {
		return "", erro
	}

	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		permissoesUsuario := permissoes["profile"].(string)
		return permissoesUsuario, nil
	}

	return "", errors.New("Token invalido")
}

func ExtractIdToken(r *http.Request) (string, error) {
	tokenString := extractToken(r)
	token, erro := jwt.Parse(tokenString, returnVerificationKey)
	if erro != nil {
		return "", erro
	}

	if permissoes, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		permissoesUsuario := permissoes["userId"].(string)
		return permissoesUsuario, nil
	}

	return "", errors.New("Token invalido")
}

// Verifica o perfil do usuario para saber se possui permissão ou não
func HasPermission(r *http.Request, profiles []string) bool {

	fmt.Println(profiles)

	userPermissions, erro := ExtractPermissions(r)
	if erro != nil {
		return false
	}

	//permissoes := []string{"admin", "editor", "usuario", "suporte"}
	for _, p := range profiles {
		if p == userPermissions {
			return true
		}
	}
	return false
}
