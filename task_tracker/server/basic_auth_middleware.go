package server

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

var namepass = map[[32]byte][32]byte{}

func init() {
	namepass[sha256.Sum256([]byte("alice"))] = sha256.Sum256([]byte("pa55word"))
}

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			if userAndPassExists(usernameHash, passwordHash) {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Unauthenticated", http.StatusUnauthorized)
		return
	})
}

// Note: probably not the best implementation because username hash should also be compared with subtle.
func userAndPassExists(unHash, pHash [32]byte) bool {
	if passHash, userFound := namepass[unHash]; userFound {
		// ConstantTimeCompare is used against timing attack.
		return subtle.ConstantTimeCompare(passHash[:], pHash[:]) == 1
	}
	return false
}
