package items

import (
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
}

func (a *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		writeError(w, http.StatusBadRequest, "no credentials provided")
		return
	}
	// in admin and secret are placeholders until there is an actual user database
	if creds.Username != "admin" || creds.Password != "secret" {
		writeError(w, http.StatusUnauthorized, "user password combination not known")
		return
	}
	token, err := GenerateToken(creds.Username)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "token couldn't be generated")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
