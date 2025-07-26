package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"web_socket/src/application"
	infrastructure "web_socket/src/infraestructure"
)

func main() {

	wsService := application.NewWebsocketService()

	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not loaded: %v", err)
	}

	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
	}))

	infrastructure.Routes(engine, wsService)

	rabbitConsumer, err := infrastructure.NewRabbitMQConsumer(wsService)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}

	go rabbitConsumer.StartConsuming()

	if err := engine.Run(":8084"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
