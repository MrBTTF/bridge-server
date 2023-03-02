package server

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := &AuthRequest{}

		bodyBytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close() //  must close
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := render.Bind(r, data); err != nil {
			renderError(w, r, http.StatusBadRequest, ErrServerForbidden, err)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		err := s.authService.ValidateToken(data.UserId, data.Token)
		if err != nil {
			renderError(w, r, http.StatusBadRequest, ErrServerForbidden, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
