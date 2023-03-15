package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mrbttf/bridge-server/pkg/config"
	"github.com/mrbttf/bridge-server/pkg/core"
	"github.com/mrbttf/bridge-server/pkg/log"
	"github.com/mrbttf/bridge-server/pkg/repositories"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/mrbttf/bridge-server/docs"
)

var (
	ErrServerBadRequest = errors.New("Bad request occured")
	ErrServerInternal   = errors.New("Internal server error")
	ErrServerForbidden  = errors.New("Forbidden")

	ErrServerSessionIdInvalid  = errors.New("session_id parameter is invalid")
	ErrServerSessionIdNotFound = errors.New("Session ID not found")
)

type Server struct {
	router         *chi.Mux
	sessionService core.SessionServicePort
	authService    core.AuthServicePort
}

func New(
	sessionService core.SessionServicePort,
	authService core.AuthServicePort,
	config config.Config,
) *Server {
	s := &Server{
		router:         chi.NewRouter(),
		sessionService: sessionService,
		authService:    authService,
	}

	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.With(s.AuthMiddleware).Get("/session/{session_id}", s.sessionGet)
	s.router.With(s.AuthMiddleware).Post("/session/create", s.sessionCreate)
	s.router.With(s.AuthMiddleware).Post("/session/lay", s.sessionLay)
	s.router.With(s.AuthMiddleware).Post("/session/pull", s.sessionPull)
	s.router.With(s.AuthMiddleware).Post("/session/nextTurn", s.sessionNextTurn)

	s.router.Post("/auth/register", s.authRegister)
	s.router.Post("/auth/login", s.authLogin)
	s.router.Post("/auth/logout", s.authLogout)

	s.router.Get("/health", s.health)
	s.router.Get("/docs/*", httpSwagger.WrapHandler)

	return s
}

func (s *Server) Run(addr string) error {
	log.Info("Server running on ", addr)
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

// session/ godoc
// @Summary Get session
// @Description Gets game session for session_id
// @Tags session
// @Produce  json
// @Param session_id path string true "ID of session"
// @Param token query string true "token"
// @Param user_id query string true "user_id"
// @Success 200 {object} sessionGetResponse
// @Failure 500 {object} ErrResponse
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
// @Param session_body body sessionCreateRequest true "Ids of players in the session"
// @Success 200 {object} sessionCreateResponse
// @Failure 500 {object} ErrResponse
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
// @Failure 500 {object} ErrResponse
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
// @Failure 500 {object} ErrResponse
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
// @Failure 500 {object} ErrResponse
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

// auth/register godoc
// @Summary Registers user
// @Description Registers user
// @Tags auth
// @Accept   json
// @Produce  json
// @Param register_body body authRegisterRequest true "Body"
// @Success 200 {object} authRegisterResponse
// @Failure 500 {object} ErrResponse
// @Router /auth/register [post]
func (s *Server) authRegister(w http.ResponseWriter, r *http.Request) {
	data := &authRegisterRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	err := s.authService.Register(
		data.Email,
		data.Password,
		data.Nickname,
	)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, ErrServerInternal, err)
		return
	}
	render.Render(w, r, &authRegisterResponse{})
}

// auth/login godoc
// @Summary Logs user in
// @Description Logs user in and returns user data
// @Tags auth
// @Accept   json
// @Produce  json
// @Param login_body body authLoginRequest true "Body"
// @Success 200 {object} authLoginResponse
// @Failure 500 {object} ErrResponse
// @Router /auth/login [post]
func (s *Server) authLogin(w http.ResponseWriter, r *http.Request) {
	data := &authLoginRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusForbidden, ErrServerForbidden, err)
		return
	}
	user, err := s.authService.Login(
		data.Email,
		data.Password,
	)
	if err != nil {
		renderError(w, r, http.StatusForbidden, ErrServerForbidden, err)
		return
	}
	render.Render(w, r, &authLoginResponse{
		User: *NewUserResponse(&user),
	})
}

// auth/logout godoc
// @Summary Logs user in
// @Description Logs user in and returns token
// @Tags auth
// @Accept   json
// @Produce  json
// @Param logout_body body authLogoutRequest true "Body"
// @Success 200 {object} authLogoutResponse
// @Failure 500 {object} ErrResponse
// @Router /auth/logout [post]
func (s *Server) authLogout(w http.ResponseWriter, r *http.Request) {
	data := &authLogoutRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusForbidden, ErrServerForbidden, err)
		return
	}
	err := s.authService.Logout(
		data.Email,
		data.Token,
	)
	if err != nil {
		renderError(w, r, http.StatusForbidden, ErrServerForbidden, err)
		return
	}
	render.Render(w, r, &authLogoutResponse{})
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
