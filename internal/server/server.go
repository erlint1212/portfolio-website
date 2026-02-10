package server

import (
	"net/http"
	"github.com/a-h/templ"
	"github.com/erlint1212/portfolio/internal/views"
)

type Server struct {
	port int
}

func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	home_component := views.Home("Apprentice")
	projects_component := views.ProjectList()
	
	mux.Handle("/", templ.Handler(home_component))
	mux.Handle("/projects", templ.Handler(projects_component))

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
