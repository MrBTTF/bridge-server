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

	ErrServerRoomIdInvalid  = errors.New("room_id parameter is invalid")
	ErrServerRoomIdNotFound = errors.New("Room ID not found")

	ErrServerUserIdInvalid  = errors.New("user_id parameter is invalid")
	ErrServerUserIdNotFound = errors.New("User ID not found")
	ErrServerUserNoSession  = errors.New("User has no session")
)

type Server struct {
	router         *chi.Mux
	sessionService core.SessionServicePort
	roomService    core.RoomServicePort
	authService    core.AuthServicePort
}

func New(
	sessionService core.SessionServicePort,
	roomService core.RoomServicePort,
	authService core.AuthServicePort,
	config config.Config,
) *Server {
	s := &Server{
		router:         chi.NewRouter(),
		sessionService: sessionService,
		authService:    authService,
		roomService:    roomService,
	}

	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	s.router.With(s.AuthMiddleware).Get("/session/{session_id}", s.sessionGet)
	s.router.With(s.AuthMiddleware).Post("/session/getByUser", s.sessionGetByUser)
	s.router.With(s.AuthMiddleware).Post("/session/create", s.sessionCreate)
	s.router.With(s.AuthMiddleware).Post("/session/lay", s.sessionLay)
	s.router.With(s.AuthMiddleware).Post("/session/pull", s.sessionPull)
	s.router.With(s.AuthMiddleware).Post("/session/nextTurn", s.sessionNextTurn)
	s.router.With(s.AuthMiddleware).Post("/session/close", s.sessionClose)

	s.router.With(s.AuthMiddleware).Get("/room/{room_id}", s.roomGet)
	s.router.With(s.AuthMiddleware).Post("/room/create", s.roomCreate)
	s.router.With(s.AuthMiddleware).Post("/room/list", s.roomList)
	s.router.With(s.AuthMiddleware).Post("/room/join", s.roomJoin)
	s.router.With(s.AuthMiddleware).Post("/room/delete", s.roomDelete)

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

// session/ godoc
// @Summary Get session by user id
// @Description Gets game session for user_id
// @Tags session
// @Produce  json
// @Param session_body body sessionGetByUserRequest true "body"
// @Success 200 {object} sessionGetByUserResponse
// @Failure 204 {object} sessionNoSessionErrorResponse
// @Failure 500 {object} ErrResponse
// @Router /session/getByUser [post]
func (s *Server) sessionGetByUser(w http.ResponseWriter, r *http.Request) {
	data := &sessionGetByUserRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, ErrServerUserIdInvalid)
		return
	}

	player, err := s.sessionService.GetPlayer(data.UserId)
	if err != nil {
		renderError(w, r, http.StatusNotFound, ErrServerUserIdNotFound, err)
		return
	}

	if player.SessionId == "" {
		renderError(w, r, http.StatusNoContent, ErrServerUserNoSession, ErrServerUserNoSession)
		return
	}

	render.Render(w, r, &sessionGetByUserResponse{
		SessionId: player.SessionId,
	})
}

// session/create godoc
// @Summary Creates session
// @Description Creates a game session and returns its id
// @Tags session
// @Accept   json
// @Produce  json
// @Param session_body body sessionCreateRequest true "body"
// @Success 200 {object} sessionCreateResponse
// @Failure 500 {object} ErrResponse
// @Router /session/create [post]
func (s *Server) sessionCreate(w http.ResponseWriter, r *http.Request) {
	data := &sessionCreateRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	session_id, err := s.sessionService.Create(data.RoomId, nil)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, ErrServerInternal, err)
		return
	}

	err = s.roomService.Close(data.RoomId)
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

// session/close godoc
// @Summary Closes session
// @Description Deletes a session and its players
// @Tags session
// @Accept   json
// @Produce  json
// @Param body body sessionCloseRequest true "Body"
// @Success 200 {object} DefaultResponse
// @Failure 500 {object} ErrResponse
// @Router /session/close [post]
func (s *Server) sessionClose(w http.ResponseWriter, r *http.Request) {
	data := &sessionCloseRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	err := s.sessionService.DeleteSession(data.SessionId)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err, err)
		return
	}
	render.Render(w, r, &DefaultResponse{})
}

// room/ godoc
// @Summary Get room
// @Description Gets room for room_id
// @Tags room
// @Produce  json
// @Param room_id path string true "ID of room"
// @Param token query string true "token"
// @Param user_id query string true "user_id"
// @Success 200 {object} roomGetResponse
// @Failure 500 {object} ErrResponse
// @Router /room/{room_id} [get]
func (s *Server) roomGet(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "room_id")
	if roomId == "" {
		renderError(w, r, http.StatusBadRequest, ErrServerRoomIdInvalid, ErrServerRoomIdInvalid)
		return
	}

	room, err := s.roomService.Get(roomId)
	if err != nil {
		renderError(w, r, http.StatusNotFound, ErrServerRoomIdNotFound, err)
		return
	}
	users, err := s.roomService.GetUsers(room.Id)
	if err != nil {
		renderError(w, r, http.StatusNotFound, ErrServerRoomIdNotFound, err)
		return
	}
	response := NewRoomResponse(&room, users)

	render.Render(w, r, &roomGetResponse{
		Room: *response,
	})
}

// room/create godoc
// @Summary Creates room
// @Description Creates a room and returns its id
// @Tags room
// @Accept   json
// @Produce  json
// @Param room_body body roomCreateRequest true "body"
// @Success 200 {object} roomCreateResponse
// @Failure 500 {object} ErrResponse
// @Router /room/create [post]
func (s *Server) roomCreate(w http.ResponseWriter, r *http.Request) {
	data := &roomCreateRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	room_id, err := s.roomService.Create(data.HostId)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, ErrServerInternal, err)
		return
	}
	render.Render(w, r, &roomCreateResponse{
		RoomId: room_id,
	})
}

// room/join godoc
// @Summary Joins a user to room
// @Description Joins a user to room by user id
// @Tags room
// @Accept   json
// @Produce  json
// @Param body body roomJoinRequest true "Body"
// @Success 200 {object} DefaultResponse
// @Failure 500 {object} ErrResponse
// @Router /room/join [post]
func (s *Server) roomJoin(w http.ResponseWriter, r *http.Request) {
	data := &roomJoinRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	err := s.roomService.Join(data.RoomId, data.UserId)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err, err)
		return
	}
	render.Render(w, r, &DefaultResponse{})
}

// room/list godoc
// @Summary List rooms
// @Description List open or closed rooms
// @Tags room
// @Accept   json
// @Produce  json
// @Param body body roomListRequest true "Body"
// @Success 200 {object} roomListResponse
// @Failure 500 {object} ErrResponse
// @Router /room/list [post]
func (s *Server) roomList(w http.ResponseWriter, r *http.Request) {
	data := &roomListRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	rooms, err := s.roomService.List(data.Open)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, err, err)
		return
	}
	roomUsers := make([]RoomUser, 0, len(rooms))
	for _, room := range rooms {
		users, err := s.roomService.GetUsers(room.Id)
		if err != nil {
			renderError(w, r, http.StatusNotFound, ErrServerRoomIdNotFound, err)
			return
		}
		roomUsers = append(roomUsers, RoomUser{
			room:  room,
			users: users,
		})
	}
	render.Render(w, r, NewRoomListResponse(roomUsers))
}

// room/delete godoc
// @Summary Deletes room
// @Description Deletes room
// @Tags room
// @Accept   json
// @Produce  json
// @Param body body roomDeleteRequest true "Body"
// @Success 200 {object} DefaultResponse
// @Failure 500 {object} ErrResponse
// @Router /room/delete [post]
func (s *Server) roomDelete(w http.ResponseWriter, r *http.Request) {
	data := &roomDeleteRequest{}

	if err := render.Bind(r, data); err != nil {
		renderError(w, r, http.StatusBadRequest, ErrServerBadRequest, err)
		return
	}
	fmt.Println("roomDelete")
	fmt.Println(data)
	err := s.roomService.Delete(data.RoomId)
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
	fmt.Printf("%+v\n", data)
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
