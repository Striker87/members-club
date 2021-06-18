package members

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"time"

	"github.com/Striker87/members-club/storage"
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

func (s Server) notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	s.executeTemplate(w, "404.html", nil)
}

func (s *Server) index(w http.ResponseWriter, _ *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "index.html", s.store); err != nil {
		log.Fatal(err)
	}
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var regName = regexp.MustCompile(`^[a-zA-Z\.\s]+$`)

func isValidName(name string) bool {
	return regName.Match([]byte(name))
}

func (s *Server) addMemberHandler(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user storage.User
	if err := json.Unmarshal(jsonData, &user); err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !isValidName(user.Name) {
		newErrorResponse(w, "members name must contains only English letters, dots and spaces", http.StatusBadRequest)
		return
	}

	if !isEmailValid(user.Email) {
		newErrorResponse(w, "wrong email", http.StatusBadRequest)
		return
	}

	if err := user.Add(s.store); err != nil {
		newErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(statusResponse{"ok"}); err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) executeTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	if err := s.templates.ExecuteTemplate(w, tmpl, data); err != nil {
		log.Fatal(err)
	}
}

type statusResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	jsonError, _ := json.Marshal(errorResponse{message})
	http.Error(w, string(jsonError), statusCode)
}
