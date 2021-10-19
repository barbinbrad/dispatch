package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func (s *server) AuthenticateHeader(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil {
			http.Error(w, "No authorization header", http.StatusBadRequest)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authorizationHeader, "JWT ") {
			http.Error(w, "Bad authorization header", http.StatusBadRequest)
			return
		}

		jwtToken := strings.SplitAfter(authorizationHeader, " ")[1]

		if !IsAuthenticated(r, jwtToken) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func (server *server) AuthenticateCookie(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "No jwt cookie", http.StatusForbidden)
		}

		jwtToken := cookie.Value

		if !IsAuthenticated(r, jwtToken) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func IsAuthenticated(r *http.Request, jwtToken string) bool {
	// TODO: refactor to extract dongle id parameter from
	//		 request, and lookup the public key for that
	//		 identity in the database (cache) using RS256
	//		 public/private key encryption

	var (
		privateKey          = "TODO: this should be looked up per dongle id"
		bypassAuthorization = true
	)

	type ClaimsV1 struct {
		Identity string `json:identity,omitempty`
		jwt.StandardClaims
	}

	if bypassAuthorization {
		return true
	}

	r.Header.Set("Identity", "")

	token, err := jwt.ParseWithClaims(jwtToken, &ClaimsV1{}, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if err != nil {
		return false
	}

	claims, ok := token.Claims.(*ClaimsV1)

	if !ok || !token.Valid {
		return false
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return false
	} else {
		r.Header.Set("Identity", claims.Identity)
		return true
	}
}

func IsAuthorized(r *http.Request, dongleId string) bool {
	identity := r.Header.Get("Identity")
	return identity == dongleId
}
