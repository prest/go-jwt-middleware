package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/golang-jwt/jwt"
	jwtmiddleware "github.com/prest/go-jwt-middleware"
)

func main() {

	StartServer()

}

func StartServer() {
	m := martini.Classic()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	m.Get("/ping", PingHandler)
	m.Get("/secured/ping", func(w http.ResponseWriter, r *http.Request) {
		jwtMiddleware.CheckJWT(w, r)
	}, SecuredPingHandler)

	m.Run()
}

type Response struct {
	Text string `json:"text"`
}

func respondJSON(text string, w http.ResponseWriter) {
	response := Response{text}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON("All good. You don't need to be authenticated to call this", w)
}

func SecuredPingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON("All good. You only get this message if you're authenticated", w)
}
