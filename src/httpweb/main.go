package main

// http://www.ttbiji.com/

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func route() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/note", note)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	if strings.HasPrefix(r.URL.Path, "/static") {
		file := "static" + r.URL.Path[len("/static"):]
		http.ServeFile(w, r, file)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {

	}
}

func note(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	fmt.Fprintf(w, "hello note")
}

func main() {
	route()
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
