package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/(save|view)/([a-zA-Z0-9]+)$")

func main() {
	start()

	p1 := &Page{Title: "TestPage2", Body: "This is a sample Page"}
	err := p1.save()
	if err != nil {
		fmt.Println(err)
		return
	}

	p2, err := loadPage(p1.Title)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(p2.Body))

	// http.HandleFunc("/", handler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Printf("view %s\n", title)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p := &Page{Title: title, Body: string(body)}

	err = p.save()
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

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)

		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

type Page struct {
	Title string
	Body  string
}

func loadPage(title string) (*Page, error) {
	var page Page

	page.Title = title
	row := db.QueryRow("SELECT * FROM article WHERE title = ?", title)
	if err := row.Scan(&page.Title, &page.Body); err != nil {
		if err == sql.ErrNoRows {
			return &page, fmt.Errorf("pageByTitle %s, no such page", title)
		} else {
			return &page, fmt.Errorf("pageByTitle %s, %v", title, err)
		}
	}

	return &page, nil
}

func (p *Page) save() error {
	_, err := loadPage(p.Title)

	var sqlScript string
	var result sql.Result
	if err == nil {
		sqlScript = "UPDATE article set body = ? WHERE title = ?"
		result, err = db.Exec(sqlScript, p.Body, p.Title)
	} else {
		sqlScript = "INSERT INTO article (title, body) VALUES (?, ?)"
		result, err = db.Exec(sqlScript, p.Title, p.Body)
	}

	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}
