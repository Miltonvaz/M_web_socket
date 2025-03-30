package infrastructure

import (
	"web_socket/src/application"

	"github.com/gin-gonic/gin"
)

func Routes(engine *gin.Engine, wsService *application.WebsocketService) {
	wsHandler := NewWebsocketHandler(wsService)
	wsGroup := engine.Group("ws")
	wsGroup.GET("handshake", wsHandler.Upgrade) 
}
