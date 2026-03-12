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

	templ.Handler(views.GameView("/assets/games/color-shooter/index.html")).ServeHTTP(w, r)
}

func (s *Server) handlerProjects(w http.ResponseWriter, r *http.Request) {
	myProjects := []routing.Project{
		{
			Title:       "High-Performance Portfolio",
			Category:    "/infrastructure/web",
			SourceURL:   "https://github.com/erlint1212/portfolio-website",
			Description: "A distributed system serving this website. Built to demonstrate <strong>infrastructure-as-code</strong> and network optimization. Features a sub-14kB initial TCP payload for maximum speed.",
			Tags:        []string{"Go", "Templ", "Tailwind", "HTMX", "RabbitMQ"},
		},
		{
			Title:       "Event Bus Architecture",
			Category:    "/backend/messaging",
			IsLive:      true,
			Description: "Implementation of the AMQP protocol using RabbitMQ. Handles decoupling of services via Topic Exchanges. Includes a durable queue setup for message persistence.",
			Tags:        []string{"RabbitMQ", "Event Sourcing", "Concurrency"},
		},
		{
			Title:       "2D Color-Switch Shooter",
			Category:    "/games/color-shooter",
			GamePath:    "/games/color-shooter",
			Description: "A fast-paced 2D platformer and shooter. Players must rapidly switch between color masks (Red, Green, Blue) to strategically match enemy colors, reflect laser attacks, and survive against various enemy archetypes like shotguns and snipers.",
			Tags:        []string{"Godot 4", "GDScript", "2D Action"},
		},
		{
			Title:       "Fault-Tolerant AI Translation Pipeline",
			Category:    "/projects/ai-translation",
			SourceURL:   "https://github.com/erlint1212/ai-transealtion-novel-to-anki-tts",
			Description: `A highly resilient ETL pipeline that transforms unstructured text into structured learning datasets. It features a custom <strong>VRAM Safety Bridge</strong> to seamlessly hot-swap 14B LLM and 1.7B TTS models on standard 12GB hardware, granular JSON micro-caching for immediate crash recovery, and a contextual <strong>RAG-lite injection system</strong> for translation accuracy.`,
			Tags:        []string{"Python", "Qwen-2.5", "Ollama"},
		},
		{
			Title:       "Audiobook ETL Pipeline",
			Category:    "/projects/audiobook-creator",
			SourceURL:   "https://github.com/erlint1212/audiobook_creator",
			Description: `An <strong>end-to-end data pipeline</strong> that scrapes unstructured web text and fully automates the creation of audiobooks. The transformation process includes text cleaning, machine translation via the Gemini API, local TTS generation, audio conversion to OPUS, and automated metadata tagging for seamless library integration.`,
			Tags:        []string{"Python", "Web Scraping", "APIs"},
		},
		{
			Title:       "Bare-Metal K3s Portfolio",
			Category:    "/projects/bare-metal-portfolio",
			SourceURL:   "https://github.com/erlint1212/portfolio-website",
			Description: `A self-hosted fullstack web application running on a dedicated physical server configured declaratively with <strong>NixOS</strong>. The infrastructure is orchestrated using K3s (Kubernetes), utilizing an event-driven <strong>RabbitMQ</strong> architecture for real-time interaction logging, and is fully monitored by a Prometheus and Grafana stack.`,
			Tags:        []string{"Go", "Kubernetes", "NixOS"},
		},
	}

	if r.Header.Get("HX-Request") == "true" {
		templ.Handler(views.ProjectListHTMX(myProjects)).ServeHTTP(w, r)
	} else {
		templ.Handler(views.ProjectsPage(myProjects)).ServeHTTP(w, r)
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	home_component := views.Home("Portfolio Page")

	mux.Handle("/", templ.Handler(home_component))
	mux.HandleFunc("/projects", s.handlerProjects)

	file_server := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", addGodotHeaders(file_server)))

	mux.HandleFunc("/games/color-shooter", s.handlerUserStartGame)

	return gzipMiddleware(mux)
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
