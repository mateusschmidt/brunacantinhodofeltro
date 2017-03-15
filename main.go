package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mateusschmidt/brunacantinhodofeltro/api/tweet"
)

type Page struct {
	Title string
	Body  []byte
}

func loadPage(title string) (*Page, error) {
	filename := "views/" + title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	if title == "" {
		title = "index"
	}
	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p.Body)
}

//Criação de rotas para API
func registerRouterAPI(router *mux.Router) {
	tweet.RegisterRouter(router)
}

//Criação de escuta para renderização arquivos front
func registerListenAndServe(router *mux.Router) {
	router.HandleFunc("/", viewHandler).Methods("GET")
	router.HandleFunc("/{id}", viewHandler).Methods("GET")
}

func main() {
	router := mux.NewRouter()
	registerListenAndServe(router)
	registerRouterAPI(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}
