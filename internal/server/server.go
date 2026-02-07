package server

import (
	"net/http"
	"github.com/a-h/templ"
	"github.com/erlint1212/portfolio/internal/views"
)

type Server struct {
	port int
}

func NewServer() *http.Server {
	const port = "8000"

	mux := http.NewServeMux()

	component := views.Hello("Apprentice")
	
	mux.Handle("/", templ.Handler(component))
	
	srv := &http.Server{
		Addr:		":" + port,
		Handler:	mux,
	}

	return srv

}
