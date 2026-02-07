package main

import (
	"github.com/erlint1212/portfolio/internal/server"
	"log"
)

func main() {
	mux := server.RegisterRoutes()
	srv := server.NewServer(mux)

	log.Printf("Serving on: http://localhost%s/\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
