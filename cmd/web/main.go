package main

import (
	"github.com/erlint1212/portfolio/internal/server"
	"log"
)

func main() {
	srv, client := server.NewServer()
	defer func() {
		err := client.Close()
		if err != nil {
			log.Printf("[WARNING] Failed to close connection: %v", err)
		}
	}()

	log.Printf("Serving on: http://localhost%s/\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
