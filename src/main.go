package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gocapt.cha/canvas"
)

func RecoverFrom() {
	if r := recover(); r != nil {
		log.Println("Recovered from panic:", r)
	}
}

func SolveCaptcha(w http.ResponseWriter, r *http.Request) {
	solution, err := canvas.SolutionFromJson(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	rs := solution.Validate()
	if rs != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, rs)
		return
	}
}

func GetCaptcha(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	captcha := canvas.Make()
	encodedCaptcha, err := captcha.ToJson()
	if err != nil {
		panic(err)
	}
	_, err = w.Write(encodedCaptcha)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

func main() {
	defer RecoverFrom()

	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Panicln("Couldn't open log file", err)
	}
	defer file.Close()
	log.SetOutput(file)

	c := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://127.0.0.1:8000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	log.Println("Started")

	r := mux.NewRouter()
	r.HandleFunc("/captcha/get", GetCaptcha)
	r.HandleFunc("/captcha/solve", SolveCaptcha)
	http.ListenAndServe("127.0.0.1:8080", c(r))
}
