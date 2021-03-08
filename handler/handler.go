package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lateralusd/shorty/db"
)

type Env struct {
	DB     *db.Config
	Scheme string
}

type Handler struct {
	*Env
	H func(e *Env, w http.ResponseWriter, r *http.Request)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.H(h.Env, w, r)
}

/* IndexPath will redirect user to shorted link or return http.StatusNotFound */
func IndexPath(env *Env, w http.ResponseWriter, r *http.Request) {
	shorted := strings.TrimLeft(r.URL.Path, "/")
	url := env.DB.GetShorted(shorted)

	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "URL is not shorted")
	} else {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

/* ShortyPath will create shorted link */
func ShortyPath(env *Env, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintln(w, "Error fetching form data")
		return
	}
	url := r.FormValue("url")
	if url == "" {
		fmt.Fprintln(w, "You need to pass url value")
		return
	}

	shorted, err := env.DB.SaveLink(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Something went wrong")
	}

	host := r.Host
	link := env.Scheme + host + "/" + shorted

	fmt.Fprintln(w, link)
}
