package application

import (
	"encoding/json"
	"log"
	"sync"
	"web_socket/src/domain"

	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	clients map[string]map[string]*domain.Session
	mu      sync.Mutex
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		clients: make(map[string]map[string]*domain.Session),
	}
}

func (ws *WebsocketService) Register(clientID, sessionID string, conn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.clients[clientID] == nil {
		ws.clients[clientID] = make(map[string]*domain.Session)
	}

	ws.clients[clientID][sessionID] = domain.NewSession(conn, sessionID, clientID, ws.clients[clientID])

	log.Printf("Usuario %s con sesión %s registrado\n", clientID, sessionID)
}

func (ws *WebsocketService) Remove(clientID, sessionID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if sessions, exists := ws.clients[clientID]; exists {
		if _, found := sessions[sessionID]; found {
			delete(sessions, sessionID)
			log.Printf("Sesión %s de usuario %s eliminada\n", sessionID, clientID)
		}
		if len(sessions) == 0 {
			delete(ws.clients, clientID)
		}
	}
}
func (ws *WebsocketService) SendMessage(clientID, message string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if sessions, exists := ws.clients[clientID]; exists {
		for sessionID, session := range sessions {

			messageObj := map[string]string{"mensaje": message}
			messageJSON, err := json.Marshal(messageObj)
			if err != nil {
				log.Printf("Error al serializar mensaje para usuario %s, sesión %s: %v\n", clientID, sessionID, err)
				continue
			}
			log.Printf("Enviando mensaje al usuario %s, sesión %s: %s", clientID, sessionID, string(messageJSON))

			err = session.Conn.WriteMessage(websocket.TextMessage, messageJSON)
			if err != nil {
				log.Printf("Error enviando mensaje a sesión %s de usuario %s: %v\n", sessionID, clientID, err)
				session.Conn.Close()
				delete(sessions, sessionID) 
			}
		}
	} else {
		log.Printf("No se encontraron sesiones activas para el usuario %s", clientID)
	}
}
