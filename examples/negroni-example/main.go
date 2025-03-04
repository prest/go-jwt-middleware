package main

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	jwtmiddleware "github.com/prest/go-jwt-middleware"
	"github.com/urfave/negroni"
)

func main() {

	StartServer()

}

func StartServer() {
	r := mux.NewRouter()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r.HandleFunc("/ping", PingHandler)
	r.Handle("/secured/ping", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(SecuredPingHandler)),
	))
	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
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
