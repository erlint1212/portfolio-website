package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"net/http"

	"github.com/erlint1212/portfolio/internal/messaging"
	"github.com/erlint1212/portfolio/internal/routing"
	"github.com/erlint1212/portfolio/internal/server"
)

func main() {
	const amqpURL = "amqp://guest:guest@localhost:5672/"
	const port = ":8002"

	client, err := messaging.NewClient(amqpURL)
	if err != nil {
		log.Printf("[WARNING] Could not connect to RabbitMQ: %v. Running in offline mode.", err)
	} else {
		defer client.Close()
	}

	var publisher *messaging.Publisher
	if client != nil {
		publisher, err = messaging.NewPublisher(client.Conn)
		if err != nil {
			log.Printf("[ERROR] Could not create publisher: %v", err)
		} else {
			defer publisher.Close()
		}
	}

	srv := server.NewServer(port, client, publisher)

	if client != nil {
		err := messaging.Subscribe(
			client.Conn,
			routing.ExchangePortfolioTopic,
			routing.GameLogSlug,
			routing.GameLogSlug+".*",
			routing.Durable,
			messaging.HandlerWriteLog(), 
			messaging.UnmarshalGob,
		)
		if err != nil {
			log.Printf("[WARNING] Failed to subscribe to %s: %v\n", routing.GameLogSlug, err)
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Serving on: http://localhost%s/", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	<- stop
	log.Println("\nShutting down server...")

	if client != nil {
		client.Close()
	}
	log.Println("Server exited properly")
}
