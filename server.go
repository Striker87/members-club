package members

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/Striker87/members_club/storage"
	"github.com/gorilla/mux"
	"github.com/jordan-wright/unindexed"
)

type Server struct {
	router     *mux.Router
	HttpServer *http.Server
	templates  *template.Template
	store      map[string]storage.User
}

func Run(port string, router *mux.Router, store map[string]storage.User) Server {
	server := Server{
		router:    router,
		templates: template.Must(template.ParseGlob("./templates/*.html")),
		store:     store,
		HttpServer: &http.Server{
			Handler:        router,
			Addr:           ":" + port,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20, // 1 mb
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		},
	}
	server.initRouter()

	return server
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}

func (s *Server) initRouter() {
	s.router.PathPrefix("/assets/css/").Handler(http.StripPrefix("/assets/css/", http.FileServer(unindexed.Dir("./assets/css"))))
	s.router.PathPrefix("/assets/js/").Handler(http.StripPrefix("/assets/js/", http.FileServer(unindexed.Dir("./assets/js"))))

	s.router.HandleFunc("/", s.index).Methods("GET")
	s.router.HandleFunc("/add_member", s.addMemberHandler).Methods("POST")
	s.router.NotFoundHandler = http.HandlerFunc(s.notFound)
}

func (s *Server) executeTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	if err := s.templates.ExecuteTemplate(w, tmpl, data); err != nil {
		log.Fatal(err)
	}
}
