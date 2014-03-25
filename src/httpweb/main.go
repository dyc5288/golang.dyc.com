package main

// http://www.ttbiji.com/

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println("token:", token)

		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]

		if len(username) == 0 {
			fmt.Fprintln(w, "username not empty")
			return
		}

		if len(password) == 0 {
			fmt.Fprintln(w, "password not empty")
			return
		}

		fmt.Println("username:", template.HTMLEscapeString(username))
		fmt.Println("password:", template.HTMLEscapeString(password))
		template.HTMLEscape(w, []byte(username))
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
