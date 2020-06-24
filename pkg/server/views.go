package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mrbttf/bridge-server/pkg/db"
	"github.com/mrbttf/bridge-server/pkg/game"
)

type app struct {
	db db.DB
}

func renderJson(data interface{}, w http.ResponseWriter) error {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

func errorResponse(w http.ResponseWriter, httpError int, desc string) {
	w.WriteHeader(httpError)
	renderJson(map[string]string{"result": "error", "desc": desc}, w)
}

func (a *app) sessionCreate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct {
		SessionName string
		HostID      string
	}
	err := decoder.Decode(&req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.SessionName) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid sessionName")
		return
	}
	if len(req.HostID) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid hostId")
		return
	}

	sessionID, err := a.db.CreateSession(req.SessionName, req.HostID)

	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = renderJson(map[string]string{
		"result":    "ok",
		"sessionId": sessionID,
	}, w)

	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (a *app) sessionJoin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct {
		SessionID  string
		PlayerName string
	}
	err := decoder.Decode(&req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.SessionID) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid sessionId")
		return
	}
	if len(req.PlayerName) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid playerName")
		return
	}

	err = a.db.JoinSession(req.SessionID, req.PlayerName)

	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = renderJson(map[string]string{
		"result": "ok",
	}, w)

	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (a *app) sessionStart(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct {
		SessionID string
	}
	err := decoder.Decode(&req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.SessionID) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid sessionId")
		return
	}

	session, err := a.db.GetSession(req.SessionID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	session = game.InitSession(session)
	err = a.db.SaveSession(session)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]interface{}{
		"result":  "ok",
		"players": session.Players,
	}
	err = renderJson(result, w)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *app) sessionEndTurn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct {
		SessionID string
		Player    string
	}
	err := decoder.Decode(&req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.SessionID) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid sessionId")
		return
	}

	session, err := a.db.GetSession(req.SessionID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	session, err = game.EndTurn(session, req.Player) 
	if err != nil {
		errorResponse(w, http.StatusNotAcceptable, err.Error())
		return
	}

	err = a.db.SaveSession(session)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]interface{}{
		"result":  "ok",
	}
	err = renderJson(result, w)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *app) cardLay(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct {
		SessionID string
		Player    string
		Card      string
		SuitOrdered      string
	}
	err := decoder.Decode(&req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.SessionID) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid sessionId")
		return
	}
	if len(req.Player) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid player")
		return
	}
	if len(req.Card) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid card")
		return
	}

	session, err := a.db.GetSession(req.SessionID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	session, err = game.LayCard(session, req.Player, req.Card, req.SuitOrdered)
	if err != nil {
		errorResponse(w, http.StatusNotAcceptable, err.Error())
		return
	}

	err = a.db.SaveSession(session)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = renderJson(map[string]string{
		"result": "ok",
	}, w)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *app) cardUnlay(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct {
		SessionID string
		Player    string
	}
	err := decoder.Decode(&req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.SessionID) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid sessionId")
		return
	}
	if len(req.Player) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid player")
		return
	}

	session, err := a.db.GetSession(req.SessionID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	session, err = game.UnlayCard(session, req.Player)
	if err != nil {
		errorResponse(w, http.StatusNotAcceptable, err.Error())
		return
	}

	err = a.db.SaveSession(session)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = renderJson(map[string]string{
		"result": "ok",
	}, w)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}



func (a *app) cardOrderSuit(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct {
		SessionID string
		Player    string
		Rank      string
	}
	err := decoder.Decode(&req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.SessionID) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid sessionId")
		return
	}
	if len(req.Player) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid player")
		return
	}
	if len(req.Rank) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid rank")
		return
	}

	session, err := a.db.GetSession(req.SessionID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	session, err = game.OrderSuit(session, req.Player, req.Rank)
	if err != nil {
		errorResponse(w, http.StatusNotAcceptable, err.Error())
		return
	}

	err = a.db.SaveSession(session)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = renderJson(map[string]string{
		"result": "ok",
	}, w)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *app) deckPull(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req struct {
		SessionID string
		Player    string
	}
	err := decoder.Decode(&req)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.SessionID) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid sessionId")
		return
	}
	if len(req.Player) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid player")
		return
	}

	session, err := a.db.GetSession(req.SessionID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	session, err = game.PullDeck(session, req.Player)
	if err != nil {
		errorResponse(w, http.StatusNotAcceptable, err.Error())
		return
	}

	err = a.db.SaveSession(session)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = renderJson(map[string]string{
		"result": "ok",
	}, w)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func Start(port string) error {
	db, err := db.New("db")
	if err != nil {
		return err
	}
	defer db.Close()

	app := &app{
		db: *db,
	}

	r := mux.NewRouter()

	session := r.PathPrefix("/session").Subrouter()
	session.HandleFunc("/create", app.sessionCreate).Methods("POST")
	session.HandleFunc("/join", app.sessionJoin).Methods("POST")
	session.HandleFunc("/start", app.sessionStart).Methods("POST")
	session.HandleFunc("/endTurn", app.sessionEndTurn).Methods("POST")

	card := r.PathPrefix("/card").Subrouter()
	card.HandleFunc("/lay", app.cardLay).Methods("POST")
	card.HandleFunc("/unlay", app.cardUnlay).Methods("POST")

	deck := r.PathPrefix("/deck").Subrouter()
	deck.HandleFunc("/pull", app.deckPull).Methods("POST")

	fmt.Printf("server started at port :%s\n", port)
	return http.ListenAndServe(":"+port, r)
}
