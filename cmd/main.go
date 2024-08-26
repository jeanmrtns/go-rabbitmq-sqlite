package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rabbitmq/database"
	"rabbitmq/types"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func handleMessage(msgs <-chan amqp.Delivery, db *sql.DB) error {
	for payload := range msgs {
		var msg types.Message

		err := json.Unmarshal(payload.Body, &msg)

		if err != nil {
			log.Println("Failed to unmarshal message")
			return err
		}

		log.Println(msg.Message)
		_, err = db.Exec("INSERT INTO messages (message) VALUES (?)", msg.Message)

		if err != nil {
			log.Println("Failed to insert message into database")
			return err
		}
	}

	return nil
}

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	rabbitmqPassword := os.Getenv("RABBITMQ_PASSWORD")
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := os.Getenv("RABBITMQ_PORT")
	rabbitmqUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort)

	conn, err := amqp.Dial(rabbitmqUrl)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"golang-queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	db, err := database.GetSQLiteInstance()
	if err != nil {
		log.Fatalf("Failed to get database instance: %s", err)
	}
	defer db.Close()

	forever := make(chan bool)
	go handleMessage(msgs, db)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
