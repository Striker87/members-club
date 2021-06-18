package members_club

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Striker87/members_club/router"
	"github.com/Striker87/members_club/storage"
	"github.com/spf13/viper"
)

func TestIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Fatal(err)
	}

	store := make(map[string]storage.User)
	route := router.Set()

	srv := Run(viper.GetString("port"), route, store)

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.index)

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, http.StatusOK)
	}
}

func TestAddMemberHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "http://localhost:8080/add_member", strings.NewReader(`{"name":"test name","email":"test@test.com"}`))
	if err != nil {
		t.Fatal(err)
	}

	store := make(map[string]storage.User)
	route := router.Set()

	srv := Run(viper.GetString("port"), route, store)

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(srv.addMemberHandler)

	handler.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, http.StatusOK)
	}

	headerType := resp.Header().Get("Content-type")
	wantHeaderType := "application/json"

	if headerType != wantHeaderType {
		t.Errorf("content type header does not match: got %v want %v", headerType, wantHeaderType)
	}
}
