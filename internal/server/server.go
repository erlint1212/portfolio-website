package server

import (
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/erlint1212/portfolio/internal/messaging"
	"github.com/erlint1212/portfolio/internal/routing"
	"github.com/erlint1212/portfolio/internal/views"
)

type Server struct {
	Addr      string
	Client    *messaging.Client
	Publisher *messaging.Publisher 
}

func addGodotHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlerUserStartGame(w http.ResponseWriter, r *http.Request) {
	if s.Publisher != nil {
		msg := "User started RPG Game"
		gl := routing.GameLog{
			CurrentTime: time.Now(),
			Message:     msg,
		}
		err := s.Publisher.PublishGameLog(r.Context(), gl)
		if err != nil {
			log.Printf("Failed to publish event: %v", err)
		}
	} else {
		log.Println("[WARNING] Publisher is nil, skipping log.")
	}

	templ.Handler(views.GameView("/assets/games/rpg/index.html")).ServeHTTP(w, r)
}

func (s *Server) handlerProjects(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("HX-Request") == "true" {
        templ.Handler(views.ProjectsList()).ServeHTTP(w, r)
    } else {
        templ.Handler(views.ProjectsPage()).ServeHTTP(w, r)
    }
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	home_component := views.Home("Portfolio Page")

	mux.Handle("/", templ.Handler(home_component))
	mux.HandleFunc("/projects", s.handlerProjects)

	file_server := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", addGodotHeaders(file_server)))

	mux.HandleFunc("/games/rpg", s.handlerUserStartGame)

	return mux
}

func NewServer(addr string, client *messaging.Client, pub *messaging.Publisher) *Server {
	return &Server{
		Addr:      addr,
		Client:    client,
		Publisher: pub,
	}
}

func (s *Server) ListenAndServe() error {
	srv := &http.Server{
		Addr:    s.Addr,
		Handler: s.RegisterRoutes(),
	}
	return srv.ListenAndServe()
}
