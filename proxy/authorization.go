package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func (s *server) Authorize(next http.HandlerFunc) http.HandlerFunc {

	var (
		privateKey          = "TODO:change this"
		bypassAuthorization = true
	)

	return func(w http.ResponseWriter, r *http.Request) {
		if !bypassAuthorization && !IsAuthorized(r, privateKey) {
			http.Error(w, "Error: 405", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

func IsAuthorized(r *http.Request, privateKey string) bool {
	// TODO: Need to compare dongle id from request url
	//       to dongle id from token

	type ClaimsV1 struct {
		DongleId string `json:dongleId,omitempty`
		jwt.StandardClaims
	}

	r.Header.Set("DongleId", "")

	if r.Header["Authorization"] == nil {
		return false
	}

	authorizationHeader := r.Header.Get("Authorization")

	if !strings.HasPrefix(authorizationHeader, "JWT ") {
		return false
	}

	encoded := strings.SplitAfter(authorizationHeader, " ")[1]
	token, err := jwt.ParseWithClaims(encoded, &ClaimsV1{}, func(token *jwt.Token) (interface{}, error) {
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
		r.Header.Set("DongleId", claims.DongleId)
		return true
	}
}
