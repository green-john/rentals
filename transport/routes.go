package transport

import (
	"encoding/json"
	"log"
	"net/http"
	"rentals"
)

func (s *Server) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&userData)
		if err != nil {
			respond(w, http.StatusInternalServerError, "Internal Server error")
			log.Printf("[ERROR] %v", err)
			return
		}

		token, err := s.authn.Login(userData.Username, userData.Password)
		if err != nil {
			respond(w, http.StatusUnauthorized, "Not allowed")
			return
		}

		var returnToken struct {
			Token string `json:"token"`
		}
		returnToken.Token = token

		respond(w, http.StatusOK, returnToken)
	}
}

func (s *Server) profileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This must exist otherwise the middleware would have rejected it
		token := r.Header["Authorization"][0]
		user := s.authn.Verify(token)

		if user == nil {
			respond(w, http.StatusUnauthorized, "Not allowed")
			return
		}
		respond(w, http.StatusOK, user)
	})
}

func (s *Server) newClientHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var newClient struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&newClient); err != nil {
			respond(w, http.StatusInternalServerError, "Internal Server error")
			log.Printf("[ERROR] %v", err)
			return
		}

		user, err := s.userService.Create(rentals.UserCreateInput{
			Username: newClient.Username,
			Password: newClient.Password,
			Role:     "client",
		})

		if err != nil {
			respond(w, http.StatusInternalServerError, err.Error())
			log.Printf("[ERROR] %v", err)
			return
		}

		respond(w, http.StatusCreated, user)
	})
}
