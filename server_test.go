package members_club

import (
	"github.com/Striker87/members_club/router"
	"github.com/Striker87/members_club/storage"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"testing"
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

	// Проверяем HTTP код
	if resp.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resp.Code, http.StatusOK)
	}
}
