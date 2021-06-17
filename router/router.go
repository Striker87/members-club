package router

import "github.com/gorilla/mux"

func Set() *mux.Router {
	return mux.NewRouter()
}
