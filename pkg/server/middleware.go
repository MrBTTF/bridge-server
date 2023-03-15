package server

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrServerUserIdNotFound = errors.New("user_id not found")
	ErrServerTokenNotFound = errors.New("token not found")
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data *AuthRequest

		var err error
		if r.Method == "GET" {
			data, err = getQueryParams(r)
			if err != nil {
				renderError(w, r, http.StatusBadRequest, ErrServerForbidden, err)
				return
			}
		}  else {
			data, err = getBodyParams(r)
			if err != nil {
				renderError(w, r, http.StatusBadRequest, ErrServerForbidden, err)
				return
			}
		}

		err = s.authService.ValidateToken(data.UserId, data.Token)
		if err != nil {
			renderError(w, r, http.StatusBadRequest, ErrServerForbidden, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getQueryParams(r *http.Request) (*AuthRequest, error) {
	var userId, token []string
	var ok bool
	q := r.URL.Query()

	if userId, ok = q["user_id"]; !ok {
		return nil, ErrServerUserIdNotFound
	}
	if token, ok = q["token"]; !ok {
		return nil, ErrServerTokenNotFound
	}
	return &AuthRequest{
		UserId: userId[0],
		Token: token[0],
	}, nil
}
func getBodyParams(r *http.Request) (*AuthRequest, error) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	data := &AuthRequest{}
	if err := render.Bind(r, data); err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return data, nil
}