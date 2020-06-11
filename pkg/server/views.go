package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrbttf/bridge-server/pkg/db"
)

type app struct {
	db db.DB
}

func (a *app) sessionCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "session create handler")
}

func (a *app) cardLay(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "card lay handler")
}

func Start(port string) error {
	app := &app{}

	r := mux.NewRouter()

	session := r.PathPrefix("/session").Subrouter()
	session.HandleFunc("/create", app.sessionCreate).Methods("GET")

	card := r.PathPrefix("/card").Subrouter()
	card.HandleFunc("/lay", app.cardLay).Methods("GET")

	fmt.Printf("server started at port :%s", port)
	return http.ListenAndServe(":"+port, r)
}
