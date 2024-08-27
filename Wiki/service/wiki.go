package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"example.com/wiki/page"
)

var validPath = regexp.MustCompile("^/(save|view)/([a-zA-Z0-9]+)$")

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	title := m[2]

	p, err := page.LoadPage(title)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// fmt.Printf("view %s\n", title)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}
	title := m[2]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p := &page.Page{Title: title, Body: string(body)}

	err = p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("save %v\n", p)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
