package server

import (
	"fmt"
	"log"
	"net/http"

	"CivicCoinMain/configs"
	"CivicCoinMain/pkg/auth"
	"CivicCoinMain/pkg/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

type NodeClaims struct {
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}

type NodeConnected struct {
	Id   string
	Addr string
}

var cfg *configs.Configs
var nodes *[]NodeConnected

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Aceptar la conexión desde cualquier origen
		return true
	},
}

func GetNodes() *[]NodeConnected {
	return nodes
}

func isValidToken(tokenString string) (bool, string) {

	// Parsear el token sin verificar
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &NodeClaims{})

	// Si hubo un error al parsear el token
	if err != nil {
		log.Fatal("Error al parsear el token: ", err)
	}

	// Acceder a las claims
	if claims, ok := token.Claims.(*NodeClaims); ok {
		// Imprimir todas las claims
		log.Printf("exp: %v\n", claims.Exp)
		log.Printf("iat: %v\n", claims.Iat)
		log.Printf("sub: %v\n", claims.Sub)

		// Check valid id
		if !utils.Contains(cfg.Keys.Nodes, claims.Sub) {
			log.Println("Invalid node id:", claims.Sub)
			return false, ""
		}

		// Check jwt signal
		publicKey, err := utils.ReadPublicKey("keys/nodes/" + claims.Sub + "_public.pem")
		if err != nil {
			log.Printf("no se ha podido leer la clave publica: %v\n", err)
			return false, ""
		}
		// Validar el token
		_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Asegúrate de que el método de firma es el esperado
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return publicKey, nil
		})

		if err != nil {
			log.Printf("Error al validar el token: %v\n", err)
			return false, ""
		}

		// Generated token
		log.Println("valid token", token)
		return true, claims.Sub

	} else {
		log.Fatal("No se pudieron obtener las claims del token")
	}

	return false, ""
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("Recibiendo petición")
	// Obtener el token de la URL
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Println("No autorizado")
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Verificar el token
	valid, Id := isValidToken(token)
	if !valid {
		log.Println("No autorizado")
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Actualizar HTTP a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Devuelve un token nuevo firmado por el nodo principal
	loggedToken, err := auth.GenerateJWT(Id, "keys/me/"+cfg.Keys.Me+"_private.pem")
	if err != nil {
		log.Println(err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, []byte(loggedToken))
	if err != nil {
		log.Println("Error enviando el token JWT:", err)
		return
	}

	log.Println("Token JWT enviado al nodo secundario:", loggedToken)

	// Leer y escribir mensajes en WebSocket
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error leyendo mensaje:", err)
			break
		}

		// Responder al cliente
		if err := conn.WriteMessage(websocket.TextMessage, []byte("Mensaje recibido: "+string(msg))); err != nil {
			log.Println("Error escribiendo mensaje:", err)
			break
		}
	}

	// Setup node
	Addr := r.RemoteAddr
	*nodes = append(*nodes, NodeConnected{Id, Addr})
}

// Función que lanza el servidor WebSocket en segundo plano
func StartWebSocketServerNodes(address string, stop chan bool) {
	nodes = &[]NodeConnected{}
	cfg = configs.GetConfig()
	http.HandleFunc("/", handleWebSocket)
	log.Println("WebSocket server listening " + address)
	// WebSocket goroutine
	go func() {
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Println("Init server error:", err)
		}
	}()
	// Wait stop signal
	<-stop
	log.Println("WebSocket server stopped!")
}
