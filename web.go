package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Frontend struct {
	listeningAddr string
}

func NewFrontend() *Frontend {
	return &Frontend{}
}

func (f *Frontend) Init(router *mux.Router) {

	router.HandleFunc("/login", makeHTTPHandleFunc(f.createHandleFunc("login.html")))
	router.HandleFunc("/signup", makeHTTPHandleFunc(f.createHandleFunc("signup.html")))
	router.HandleFunc("/", makeHTTPHandleFunc(f.createHandleFunc("index.html")))

	router.HandleFunc("/main.js", makeHTTPHandleFunc(f.createHandleFunc("main.js")))
	router.HandleFunc("/main.css", makeHTTPHandleFunc(f.createHandleFunc("main.css")))
}

func (f *Frontend) createHandleFunc(file string) func(http.ResponseWriter, *http.Request) error {

	return func(w http.ResponseWriter, r *http.Request) error {
		http.ServeFile(w, r, "/home/moritz/Programming/Go/gobank/www/frontend/static/"+file)
		return nil
	}
}
