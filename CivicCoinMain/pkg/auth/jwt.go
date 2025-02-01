package auth

import (
	"fmt"
	"time"

	"CivicCoinMain/pkg/utils"

	"github.com/golang-jwt/jwt/v4"
)

// Create token jwk signed with private key and encripted with public key
func GenerateJWT(nodeId string, privateKeyPath string) (string, error) {

	// Leer la clave privada desde el archivo PEM
	privateKey, err := utils.ReadPrivateKey(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("no se pudo leer la clave privada: %v", err)
	}

	// Crear el token JWT con claims
	claims := jwt.MapClaims{
		"sub": nodeId,                             // Identificador del usuario
		"exp": time.Now().Add(time.Minute).Unix(), // Fecha de expiración
		"iat": time.Now().Unix(),                  // Fecha de emisión
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Firmar el token con la clave privada
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("no se pudo firmar el token: %v", err)
	}

	return signedToken, nil
}
