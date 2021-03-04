package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lateralusd/shorty/db"
)

type Env struct {
	DB *db.Config
}

type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request)
}

func main() {
	var err error
	database, err := db.InitDatabase()
	if err != nil {
		log.Fatal("Coult not initialize database", err)
	}

	env := &Env{
		DB: database,
	}

	http.Handle("/", Handler{env, IndexPath})
	http.Handle("/shorty", Handler{env, ShortyPath})

	fmt.Println("[*] Starting the server")
	http.ListenAndServe(":8080", nil)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.H(h.Env, w, r)
}

func IndexPath(env *Env, w http.ResponseWriter, r *http.Request) {
	shorted := strings.TrimLeft(r.URL.Path, "/")
	url, err := env.DB.GetShorted(shorted)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "URL is not shorted")
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func ShortyPath(env *Env, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, "Error fetching form data")
		return
	}
	url := r.FormValue("url")
	shorted, err := env.DB.SaveLink(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Something went wrong")
	}

	/* TODO: Do not hardcode link */
	link := "https://localhost:8080/" + shorted
	fmt.Fprintln(w, link)
}
