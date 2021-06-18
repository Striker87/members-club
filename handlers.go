package members_club

import (
	"encoding/json"
	validations2 "github.com/Striker87/members_club/pkg/validations"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Striker87/members_club/storage"
)

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

	if !validations2.IsValidName(user.Name) {
		newErrorResponse(w, "members name must contains only English letters, dots and spaces", http.StatusBadRequest)
		return
	}

	if !validations2.IsEmailValid(user.Email) {
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

func (s *Server) index(w http.ResponseWriter, _ *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "index.html", s.store); err != nil {
		log.Fatal(err)
	}
}

func (s Server) notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	s.executeTemplate(w, "404.html", nil)
}
