package main

// http://www.ttbiji.com/

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/Go-SQL-Driver/MySQL"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type user_obj struct {
	user_id  int
	username string
	password string
	created  string
}

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
			fmt.Fprintln(w, "用户名不能为空")
			return
		}

		if len(password) == 0 {
			fmt.Fprintln(w, "密码不能为空")
			return
		}

		db, err := sql.Open("mysql", "dyc5288:d54321@/gonote?charset=utf8")
		checkErr(err)
		var user = get_user(db, username)

		if user.user_id == 0 {
			fmt.Fprintln(w, "用户不存在")
			return
		}

		if mymd5(password) != user.password {
			fmt.Fprintln(w, "密码错误")
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

// 插入数据
func add_user(db sql.DB, username string, password string) int64 {
	stmt, err := db.Prepare("INSERT userinfo SET username=?,password=?,created=?")
	checkErr(err)

	password = mymd5(password)
	res, err := stmt.Exec(username, password)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}

// 获取用户信息
func get_user(db *sql.DB, username string) user_obj {
	rows, err := db.Query("SELECT * FROM userinfo where username='" + username + "'")
	checkErr(err)

	defer rows.Close()

	for rows.Next() {
		var user_id int
		var username string
		var password string
		var created string
		err = rows.Scan(&user_id, &username, &password, &created)
		checkErr(err)
		return user_obj{user_id, username, password, created}
	}

	return user_obj{0, "", "", ""}
}

func mymd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

/* 异常流程 */
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	route()
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
