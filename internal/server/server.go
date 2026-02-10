package server

import (
	"net/http"
	"github.com/a-h/templ"
	"github.com/erlint1212/portfolio/internal/views"
)

type Server struct {
	port int
}

func addGodotHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		next.ServeHTTP(w, r)
	})
}

func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	home_component := views.Home("Apprentice")
	projects_component := views.ProjectList()
	
	mux.Handle("/", templ.Handler(home_component))
	mux.Handle("/projects", templ.Handler(projects_component))

	file_server := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", addGodotHeaders(file_server)))

	mux.Handle("/games/rpg", templ.Handler(views.GameView("/assets/games/rpg/index.html")))

	return mux
}

func NewServer(mux http.Handler) *http.Server {
	const port = "8000"
	
	srv := &http.Server{
		Addr:		":" + port,
		Handler:	mux,
	}

	return srv

}
