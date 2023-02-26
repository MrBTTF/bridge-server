package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mrbttf/bridge-server/pkg/core/services/session"
	"github.com/mrbttf/bridge-server/pkg/log"
	"github.com/mrbttf/bridge-server/pkg/repositories"

	_ "github.com/mrbttf/bridge-server/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	ErrServerBadRequest = errors.New("Bad request occured")
	ErrServerInternal   = errors.New("Internal server error")

	ErrServerSessionIdInvalid  = errors.New("session_id parameter is invalid")
	ErrServerSessionIdNotFound = errors.New("Session ID not found")
)

type Server struct {
	router         *chi.Mux
	sessionService *session.SessionService
}

func New(sessionService *session.SessionService) *Server {
	s := &Server{
		router:         chi.NewRouter(),
		sessionService: sessionService,
	}

	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.Get("/session/{session_id}", s.sessionGet)
	s.router.Post("/session/create", s.sessionCreate)
	s.router.Post("/session/lay", s.sessionLay)
	s.router.Post("/session/pull", s.sessionPull)
	s.router.Post("/session/nextTurn", s.sessionNextTurn)

	s.router.Get("/docs/*", httpSwagger.WrapHandler)
	return s
}

func (s *Server) Run(addr string) error {
	log.Info("Server running on ", addr)
	return http.ListenAndServe(addr, s.router)
}

// session/ godoc
// @Summary Get session
// @Description Gets game session for session_id
// @Tags session
// @Produce  json
// @Param session_id path string true "ID of session"
// @Success 200 {object} sessionGetResponse
// @Router /session/{session_id} [get]
func (s *Server) sessionGet(w http.ResponseWriter, r *http.Request) {
	sessionId := chi.URLParam(r, "session_id")
	if sessionId == "" {
		renderError(w, r, http.StatusBadRequest, ErrServerSessionIdInvalid, ErrServerSessionIdInvalid)
		return
	}

	session, err := s.sessionService.GetSession(sessionId)
	if err != nil {
		renderError(w, r, http.StatusNotFound, ErrServerSessionIdNotFound, err)
		return
	}
	player, err := s.sessionService.GetPlayer(session.CurrentPlayer)
	if err != nil {
		renderError(w, r, http.StatusNotFound, ErrServerSessionIdNotFound, err)
		return
	}
	response := NewSessionResponse(&session, &player)

	render.Render(w, r, &sessionGetResponse{
		Session: *response,
	})
}

// session/create godoc
// @Summary Creates session
// @Description Creates a game session and returns its id
// @Tags session
// @Accept   json
// @Produce  json
// @Param player_ids body sessionCreateRequest true "Ids of players in the session"
// @Success 200 {object} sessionCreateResponse
// // @Failure 400 {object} failureResponse
// // @Failure 500 {object} failureResponse
// @Router /session/create [post]
func (s *Server) sessionCreate(w http.ResponseWriter, r *http.Request) {
	data := &sessionCreateRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	session_id, err := s.sessionService.Create(data.PlayerIds, nil)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, ErrServerInternal, err)
		return
	}
	render.Render(w, r, &sessionCreateResponse{
		SessionID: session_id,
	})
}

// session/lay godoc
// @Summary Lays a card
// @Description Lays a card for player and session id
// @Tags session
// @Accept   json
// @Produce  json
// @Param body body sessionLayRequest true "Body"
// @Success 200 {object} DefaultResponse
// @Router /session/lay [post]
func (s *Server) sessionLay(w http.ResponseWriter, r *http.Request) {
	data := &sessionLayRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	card := repositories.StringToCard(data.Card)
	err := s.sessionService.Lay(data.SessionId, data.PlayerId, card)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err, err)
		return
	}
	render.Render(w, r, &DefaultResponse{})
}

// session/pull godoc
// @Summary Pulls a card
// @Description Pulls a card for player and session id
// @Tags session
// @Accept   json
// @Produce  json
// @Param body body sessionPullRequest true "Body"
// @Success 200 {object} DefaultResponse
// @Router /session/pull [post]
func (s *Server) sessionPull(w http.ResponseWriter, r *http.Request) {
	data := &sessionPullRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	err := s.sessionService.Pull(data.SessionId, data.PlayerId)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err, err)
		return
	}
	render.Render(w, r, &DefaultResponse{})
}



// session/nextTurn godoc
// @Summary Next turn
// @Description Ends turn for player id and passes turn to next player
// @Tags session
// @Accept   json
// @Produce  json
// @Param body body sessionNextTurnRequest true "Body"
// @Success 200 {object} DefaultResponse
// @Router /session/nextTurn [post]
func (s *Server) sessionNextTurn(w http.ResponseWriter, r *http.Request) {
	data := &sessionNextTurnRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	err := s.sessionService.NextTurn(data.SessionId, data.PlayerId)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err, err)
		return
	}
	render.Render(w, r, &DefaultResponse{})
}

func renderError(w http.ResponseWriter, r *http.Request, code int, message error, err error) {
	if err != nil {
		log.Error(err)
	} else {
		log.Error(message)
	}
	render.Render(w, r, ErrResponse{
		Message: message.Error(),
		Code:    code,
	})
}
