package main


import (
	"os"
	"fmt"
	"log"
	"time"
	"gocapt.cha/canvas"
	"encoding/json"
	"math/rand"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/google/uuid"
	"net/http"
)


func SolveCaptcha(w http.ResponseWriter, r *http.Request) {
	var captchaRequest canvas.CaptchaRequest
	decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&captchaRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	captcha := captchaCache[captchaRequest.Key]
	delete(captchaCache, captchaRequest.Key)
	if captcha == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Captcha not found")
		return
	}
	if captcha.Index != captchaRequest.Index {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Wrong element")
		return
	}
	if !captchaRequest.Position.In(captcha.Solution) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Not in bounds")
		return
	}
}

func GetCaptcha(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	captcha := canvas.MakeCaptcha()
	captcha.Key = uuid.NewString()
	captchaCache[captcha.Key] = &captcha.Solution
	encoded, err := json.Marshal(captcha)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	_, err = w.Write(encoded)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

var captchaCache = make(map[string]*canvas.CaptchaSolution)

func main() {
	
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
		
	rand.Seed(time.Now().UnixNano())
	r := mux.NewRouter()
    r.HandleFunc("/captcha/get", GetCaptcha)
    r.HandleFunc("/captcha/solve", SolveCaptcha)
	http.ListenAndServe("127.0.0.1:8080", c(r))
}