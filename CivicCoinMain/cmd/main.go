package main

/*
import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Definir la configuración del WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir conexiones de cualquier origen
	},
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	// Actualizar la conexión HTTP a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	fmt.Println("New WebSocket connection established!")

	// Leer mensajes del cliente
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Mostrar el mensaje recibido
		fmt.Printf("Message received: %s\n", p)

		// Enviar un mensaje de vuelta al cliente
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	// Ruta para manejar las conexiones WebSocket
	http.HandleFunc("/ws", handleWebSocketConnection)

	// Iniciar el servidor HTTP
	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe failed:", err)
	}
}
*/

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"CivicCoinMain/configs"
	"CivicCoinMain/pkg/db/redis"
	"CivicCoinMain/pkg/models"
	"CivicCoinMain/pkg/nodes"
	"CivicCoinMain/pkg/server"
	"CivicCoinMain/pkg/utils"
)

var cfg *models.Configs
var err error
var stopWebSocketNodes chan bool

func init() {
	utils.ConfigureLogs()
	cfg, err = configs.LoadConfigs()
	if err != nil {
		log.Fatalf("Error loading configs: %v", err)
	}
	stopWebSocketNodes = make(chan bool)
	err = redis.Init(cfg.Redis.Addr, cfg.Redis.Pass, cfg.Redis.Db)
	if err != nil {
		log.Println("Error init Redis:", err)
		return
	}
}

func main() {

	for {
		fmt.Println("=== Menu ===")
		fmt.Println("0. Exit")
		fmt.Println("1. Create new keys")
		fmt.Println("2. Start WebSocket nodes")
		fmt.Println("3. Stop WebSocket nodes")
		fmt.Println("4. Connected nodes")
		fmt.Print("Select option: ")

		choice := readInput()

		switch choice {
		case 0:
			fmt.Println("Good bye!")
			return
		case 1:
			createKeys()
		case 2:
			go startWebSocketServerNodes()
		case 3:
			stopWebSocketServerNodes()
		case 4:
			getConnectedNodes()
		default:
			fmt.Println("Invalid option.")
		}

		fmt.Println()
	}
}

func startWebSocketServerNodes() {
	go server.StartWebSocketServerNodes(cfg.Websocket.Address, stopWebSocketNodes)
}
func stopWebSocketServerNodes() {
	stopWebSocketNodes <- true
}

func createKeys() {
	fmt.Println("Input key name or empty to create hash:")
	keysName := readStringInput()
	if keysName == "" {
		timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
		keysName = utils.GenerateHash(timestamp)
	}
	fmt.Printf("Creating keys '%s'...\n", keysName)
	privateKey, publicKey, error := utils.GenerateKeys(2048)
	if error != nil {
		log.Fatalf("Error creating keys: %v", error)
	}
	error = utils.SavePrivateKey(privateKey, "keys/"+keysName)
	if error != nil {
		log.Fatalf("Error saving private key: %v", error)
	} else {
		fmt.Printf("Private key saved!\n")
	}
	error = utils.SavePublicKey(publicKey, "keys/"+keysName)
	if error != nil {
		log.Fatalf("Error saving public key: %v", error)
	} else {
		fmt.Printf("Public key saved!\n")
	}
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

func getConnectedNodes() {
	// Obtener todos los nodos
	nodes, err := nodes.GetAllNodes()
	if err != nil {
		fmt.Println("Error obteniendo todos los nodos:", err)
		return
	}
	fmt.Println("Todos los nodos:")
	for _, n := range nodes {
		fmt.Printf("- %+v\n", n)
	}
}
