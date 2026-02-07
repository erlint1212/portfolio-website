package main

import (
	"github.com/erlint1212/portfolio/internal/server"
	"log"
)

func main() {
	srv := server.NewServer()

	log.Printf("Serving on: http://localhost%s/\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
