package server

import (
	"log"
	"net/http"
	"github.com/a-h/templ"
	"github.com/erlint1212/portfolio/internal/views"
	"github.com/erlint1212/portfolio/internal/messaging"
)

type Server struct {
	Addr string
	client *messaging.Client
}

func  addGodotHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlerUserStartGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.client != nil {
			err := s.client.Publish("game_events", "User started RPG Game")
			if err != nil {
				log.Printf("Failed to publish event: %v", err)
			}
		}

		templ.Handler(views.GameView("/assets/games/rpg/index.html")).ServeHTTP(w, r)
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	home_component := views.Home("Apprentice")
	projects_component := views.ProjectList()
	
	mux.Handle("/", templ.Handler(home_component))
	mux.Handle("/projects", templ.Handler(projects_component))

	file_server := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", addGodotHeaders(file_server)))

	mux.Handle("/games/rpg", s.handlerUserStartGame())

	return mux
}

func NewServer() (*Server, *messaging.Client) {
	const port = "8000"
	const amqpURL = "amqp://guest:guest@localhost:5672/"

	client, err := messaging.NewClient(amqpURL)
	if err != nil {
		log.Printf("[WARNING] Could not connect to RabbitMQ: %v", err)
	}
	
	srv := &Server{
		Addr:		":" + port,
		client:		client,
	}

	return srv, client

}

func (s *Server) ListenAndServe() error {
	srv := &http.Server{
		Addr: s.Addr,
		Handler: s.RegisterRoutes(),
	}
	return srv.ListenAndServe()
}
