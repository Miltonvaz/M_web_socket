package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"web_socket/src/application"
)

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	ws      *application.WebsocketService
}

type Message struct {
	UserID  string `json:"user_id"`
	Mensaje string `json:"mensaje"`
}

func NewRabbitMQConsumer(ws *application.WebsocketService) (*RabbitMQConsumer, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	queueName := os.Getenv("RABBITMQ_QUEUE_NAME")
	if rabbitURL == "" || queueName == "" {
		log.Fatal("RabbitMQ URL or queue name is not set in .env file")
		return nil, fmt.Errorf("RabbitMQ URL or queue name is not set in .env file")
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
		queue:   queueName,
		ws:      ws,
	}, nil
}

func (r *RabbitMQConsumer) StartConsuming() {
	msgs, err := r.channel.Consume(r.queue, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Error al consumir mensajes: %v", err)
	}

	for msg := range msgs {
		var message Message
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Printf("Error al deserializar el mensaje: %v", err)
			continue
		}

		if message.UserID == "" {
			log.Println("Error: user_id faltante en el mensaje")
			continue
		}

		if message.Mensaje == "" {
			log.Println("Error: mensaje faltante en el mensaje")
			continue
		}

		log.Printf("Mensaje recibido para el usuario %s: %s", message.UserID, message.Mensaje)
		r.ws.SendMessage(message.UserID, message.Mensaje)
	}
}
