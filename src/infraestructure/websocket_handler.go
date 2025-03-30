package infrastructure

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"web_socket/src/application"
	"github.com/gin-gonic/gin"
)

type WebsocketHandler struct {
	wsService *application.WebsocketService
}

func NewWebsocketHandler(wsService *application.WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{wsService: wsService}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebsocketHandler) Upgrade(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error en la conexión WebSocket:", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("Error cerrando la conexión WebSocket:", err)
		}
	}()

	clientID := c.DefaultQuery("user_id", "")
	sessionID := c.DefaultQuery("user_id", "")

	if clientID == "" || sessionID == "" {
		log.Println("Error: client_id o session_id faltante")
		conn.Close()
		return
	}

	h.wsService.Register(clientID, sessionID, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Usuario %s con sesión %s desconectado\n", clientID, sessionID)
			h.wsService.Remove(clientID, sessionID)
			break
		}
	}
}
