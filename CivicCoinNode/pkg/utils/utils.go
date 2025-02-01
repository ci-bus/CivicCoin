package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Genrate keys RSA
func GenerateKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, error := rsa.GenerateKey(rand.Reader, bits)
	return privateKey, &privateKey.PublicKey, error
}

// Save private key PEM file
func SavePrivateKey(privateKey *rsa.PrivateKey, pathFile string) error {
	file, error := os.Create(pathFile + "_private.pem")
	if error != nil {
		return error
	}
	defer file.Close()
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	error = pem.Encode(file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})
	return error
}

// Save public key PEM file
func SavePublicKey(publicKey *rsa.PublicKey, pathFile string) error {
	file, error := os.Create(pathFile + "_public.pem")
	if error != nil {
		return error
	}
	defer file.Close()
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	error = pem.Encode(file, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyBytes})
	return error
}

// Generate hash
func GenerateHash(input string) string {
	// Create a new SHA-256 hash
	hash := sha256.New()
	// Write the input data to the hash
	hash.Write([]byte(input))
	// Get the hashed result
	hashBytes := hash.Sum(nil)
	// Convert the hash to a hexadecimal string and return it
	return fmt.Sprintf("%x", hashBytes)
}

func ReadPublicKey(publicMainKeyPath string) (*rsa.PublicKey, error) {
	// Lee la clave pública desde el archivo
	publicKeyData, err := os.ReadFile(publicMainKeyPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer la clave pública: %v", err)
	}

	// Decodifica el bloque PEM
	block, _ := pem.Decode(publicKeyData)
	if block == nil {
		return nil, fmt.Errorf("no se pudo decodificar el bloque PEM de la clave pública")
	}

	// Intenta parsear como PKCS#1 (clave RSA clásica)
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return publicKey, nil
	}

	// Si no es PKCS#1, intenta como PKCS#8
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("no se pudo parsear la clave pública: %v", err)
	}

	// Asegúrate de que sea una clave RSA
	rsaKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("la clave pública no es de tipo RSA")
	}

	return rsaKey, nil
}

func ReadPrivateKey(path string) (*rsa.PrivateKey, error) {
	// Cargar y leer la clave privada desde el archivo
	privateKeyData, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error al leer la clave privada: %v", err)
		return nil, err
	}

	// Decodificar la clave privada PEM
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		log.Fatalf("No se pudo encontrar el bloque PEM de la clave privada")
		return nil, err
	}

	// Parsear la clave privada RSA
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Error al parsear la clave privada: %v", err)
		return nil, err
	}

	return privateKey, nil
}

// Search string into array
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func ConfigureLogs() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logs/main.log",
		MaxSize:    128, // MB
		MaxBackups: 10,
		MaxAge:     60, // Days
		Compress:   false,
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
	log.SetOutput(multiWriter)
}
