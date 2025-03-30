package domain

import (
	"log"
	"github.com/gorilla/websocket"
)

type Session struct {
	Conn         *websocket.Conn
	SessionID    string
	ClientID  string  
	closeHandler func(sessionID string)
	sessions     map[string]*Session
}

func NewSession(conn *websocket.Conn, sessionID, clientID string, sessions map[string]*Session) *Session {
	return &Session{
		Conn:      conn,
		SessionID: sessionID,
		ClientID:  clientID,
		sessions:  sessions,
	}
}


func (s *Session) SetCloseHandler(handler func(sessionID string)) {
	s.closeHandler = handler
}

func (s *Session) StartHandling(removeSession func(sessionID string)) {
	s.closeHandler = removeSession
	go s.readPump()
}

func (s *Session) readPump() {
	defer func() {
		s.Conn.Close()
		if s.closeHandler != nil {
			s.closeHandler(s.SessionID)
		}
	}()

	for {
		_, message, err := s.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}

		log.Printf("Received message from client %s: %s", s.SessionID, message)
		s.broadcast(websocket.TextMessage, message)
	}
}

func (s *Session) broadcast(messageType int, payload []byte) {
	err := s.Conn.WriteMessage(messageType, payload)
	if err != nil {
		log.Printf("Broadcast error to session %s: %v", s.SessionID, err)
	}
}

func (s *Session) SendMessage(messageType int, payload []byte) {
	err := s.Conn.WriteMessage(messageType, payload)
	if err != nil {
		log.Printf("Error sending message to session %s: %v", s.SessionID, err)
	}
}
