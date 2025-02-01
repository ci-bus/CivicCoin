package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"CivicCoinNode/configs"
	"CivicCoinNode/pkg/auth"
	"CivicCoinNode/pkg/utils"

	"github.com/gorilla/websocket"
)

var cfg *configs.Configs
var err error

func init() {
	utils.ConfigureLogs()
	cfg, err = configs.LoadConfigs()
	if err != nil {
		log.Fatalf("Error loading configs: %v", err)
	}
}

func main() {

	for {
		fmt.Println("=== Menu ===")
		fmt.Println("0. Exit")
		fmt.Println("1. Connect")
		fmt.Print("Select option: ")

		choice := readInput()

		switch choice {
		case 0:
			fmt.Println("Good bye!")
			return
		case 1:
			connectToMainNode()
		default:
			fmt.Println("Invalid option.")
		}

		fmt.Println()
	}
}

func createLoginToken() string {
	privatePath := "keys/me/" + cfg.Keys.Me + "_private.pem"
	token, error := auth.GenerateJWT(cfg.Keys.Me, privatePath)
	if error != nil {
		fmt.Println(error.Error())
	} else {
		fmt.Println("\nToken: " + token)
	}
	return token
}

func connectToMainNode() (*websocket.Conn, error) {
	loginToken := createLoginToken()
	// URL
	u := url.URL{Scheme: "ws", Host: cfg.MainAddress, Path: "/", RawQuery: "token=" + loginToken}
	// WebSocket connect
	log.Printf("Connecting with %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Error connect: %v", err)
		return nil, err
	}
	log.Println("Connected with main node!")
	// Recibir el token JWT del servidor
	_, tokenBytes, err := conn.ReadMessage()
	if err != nil {
		log.Fatal("Error recibiendo el token JWT:", err)
	}
	tokenString := string(tokenBytes)
	log.Println("Token JWT recibido:", tokenString)
	return conn, nil
}

func readInput() int {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	choice, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Invalid input. Enter a number.")
		return -1
	}

	return choice
}

func readStringInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
