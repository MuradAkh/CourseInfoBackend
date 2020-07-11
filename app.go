package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var ctx = context.Background()

type App struct {
	Router *mux.Router
	DB     *redis.Client
}

func (a *App) Initialize(password string, address string) {

	a.DB = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,  // use default DB
	})

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/course/{code}", a.course).Methods("GET")
	a.Router.HandleFunc("/textbook/isbn/{code}", a.textbooksIsbn).Methods("GET")
	a.Router.HandleFunc("/textbook/code/{code}", a.textbooksCode).Methods("GET")
	a.Router.HandleFunc("/textbook/title/{code}", a.textbooksTitle).Methods("GET")
}

func respondWithJSON(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}

func (a *App) course(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	fmt.Print("requested code:  " + code)
	val2, err := a.DB.Get(ctx, code).Result()
	if err == redis.Nil {
		fmt.Print("fetching from nikel: " + code)
		resp, err := http.Get("https://nikel.ml/api/courses?code=" + code)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		err = a.DB.Set(ctx, code, body, 48 * time.Hour).Err()
		if err != nil {
			panic(err)
		}
		respondWithJSON(w, 200, body)
	} else if err != nil {
		panic(err)
	} else {
		respondWithJSON(w, 200, []byte(val2))
	}

}

func (a *App) textbooksCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]
	_ = code
	w.WriteHeader(400)
}

func (a *App) textbooksTitle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	fmt.Print("fetching from nikel: " + code)
	resp, err := http.Get("https://nikel.ml/api/textbooks?title=" + code)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = a.DB.Set(ctx, code, body, 48 * time.Hour).Err()
	if err != nil {
		panic(err)
	}
	respondWithJSON(w, 200, body)
}

func (a *App) textbooksIsbn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	fmt.Print("fetching from nikel: " + code)
	resp, err := http.Get("https://nikel.ml/api/testbooks?isbn=" + code)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = a.DB.Set(ctx, code, body, 48 * time.Hour).Err()
	if err != nil {
		panic(err)
	}
	respondWithJSON(w, 200, body)

}