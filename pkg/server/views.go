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
		HostPlayer  string
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
	if len(req.HostPlayer) == 0 {
		errorResponse(w, http.StatusBadRequest, "invalid hostId")
		return
	}

	sessionID, err := a.db.CreateSession(req.SessionName, req.HostPlayer)

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

	session, err = game.InitSession(session)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = a.db.SaveSession(session)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	result := map[string]interface{}{
		"result":  "ok",
		"session": session,
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
		"result": "ok",
	}
	err = renderJson(result, w)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *app) sessionWaitForTurn(w http.ResponseWriter, r *http.Request) {
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

	session, err = game.WaitForTurn(session, req.Player)
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
		"result": "ok",
	}
	err = renderJson(result, w)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *app) sessionPlayers(w http.ResponseWriter, r *http.Request) {
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

	playersNames := []string{}
	for _, player := range session.Players {
		playersNames = append(playersNames, player.Name)
	}

	result := map[string]interface{}{
		"result":  "ok",
		"players": session.Players,
	}
	err = renderJson(result, w)

	if err != nil {
		fmt.Println(err)
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}
func (a *app) sessionGet(w http.ResponseWriter, r *http.Request) {
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

	result := map[string]interface{}{
		"result":  "ok",
		"session": session,
	}

	err = renderJson(result, w)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (a *app) sessionList(w http.ResponseWriter, r *http.Request) {

	sessions, err := a.db.ListSessions()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sessionsForJSON := []interface{}{}
	for _, session := range sessions {
		sessionsForJSON = append(sessionsForJSON, struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			PlayersCount int    `json:"playersCount"`
		}{
			ID:           session.ID,
			Name:         session.Name,
			PlayersCount: len(session.Players),
		})
	}
	result := map[string]interface{}{
		"result":   "ok",
		"sessions": sessionsForJSON,
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
		SessionID   string
		Player      string
		Card        string
		SuitOrdered string
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
	session.HandleFunc("/waitForTurn", app.sessionWaitForTurn).Methods("POST")
	session.HandleFunc("/players", app.sessionPlayers).Methods("POST")
	session.HandleFunc("/list", app.sessionList).Methods("GET")
	session.HandleFunc("/get", app.sessionGet).Methods("POST")

	card := r.PathPrefix("/card").Subrouter()
	card.HandleFunc("/lay", app.cardLay).Methods("POST")
	card.HandleFunc("/unlay", app.cardUnlay).Methods("POST")

	deck := r.PathPrefix("/deck").Subrouter()
	deck.HandleFunc("/pull", app.deckPull).Methods("POST")

	fmt.Printf("server started at port :%s\n", port)
	return http.ListenAndServe(":"+port, r)
}
