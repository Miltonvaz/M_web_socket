package main

import (
	"log"
	"web_socket/src/application"
	infrastructure "web_socket/src/infraestructure"

	"github.com/gin-gonic/gin"
)

func main() {
	

	wsService := application.NewWebsocketService()

	engine := gin.Default()
	infrastructure.Routes(engine, wsService)
	rabbitConsumer, err := infrastructure.NewRabbitMQConsumer(wsService)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}

	go rabbitConsumer.StartConsuming()

	if err := engine.Run(":8083"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
