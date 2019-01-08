package rentals

import (
	"encoding/json"
	"log"
	"net/http"
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
			ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			log.Printf("[ERROR] %v", err)
			return
		}

		token, err := s.authN.Login(userData.Username, userData.Password)
		if err != nil {
			ErrorResponse(w, http.StatusUnauthorized, "Not allowed")
			return
		}

		var returnToken struct {
			Token string `json:"token"`
		}
		returnToken.Token = token

		jsonRes, err := json.Marshal(returnToken)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			log.Printf("[ERROR] %v", err)
			return
		}

		w.Write(jsonRes)
	}
}
