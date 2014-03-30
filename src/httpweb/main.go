package main

// http://www.ttbiji.com/

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/bitly/go-simplejson" // for json get
	"html/template"
	"io"
	"log"
	_ "net"
	"net/http"
	_ "os"
	"strconv"
	"strings"
	"time"
)

/* json返回对象 */
type MyReturn struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}

/* 用户对象 */
type user_obj struct {
	user_id  int
	username string
	password string
	created  string
}

/* 路由 */
func route() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/note", note)
}

/* 其他 */
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	if strings.HasPrefix(r.URL.Path, "/static") {
		file := "static" + r.URL.Path[len("/static"):]
		http.ServeFile(w, r, file)
		return
	} else if strings.HasPrefix(r.URL.Path, "/images") {
		file := "images" + r.URL.Path[len("/images"):]
		http.ServeFile(w, r, file)
		return
	} else {
		token := ""
		t, _ := template.ParseFiles("template/index.gtpl")
		t.Execute(w, token)
	}
}

/* 登录提交 */
func login_submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/json; charset=utf-8")
	var res = MyReturn{State: false, Message: ""}
	r.ParseForm()
	username := r.Form["username"][0]
	password := r.Form["password"][0]

	if len(username) == 0 {
		res.Message = "亲爱的，用户名不能为空"
		fmt.Fprintln(w, json_encode(res))
		return
	}

	if len(password) == 0 {
		res.Message = "亲爱的，密码不能为空"
		fmt.Fprintln(w, json_encode(res))
		return
	}

	db, err := sql.Open("mysql", "root:d54321@/gonote?charset=utf8")
	checkErr(err)
	var user = get_user(db, username)

	if user.user_id == 0 {
		res.Message = "亲爱的，用户不存在"
		fmt.Fprintln(w, json_encode(res))
		return
	}

	if mymd5(password) != user.password {
		res.Message = "亲爱的，密码错误"
		fmt.Fprintln(w, json_encode(res))
		return
	}

	fmt.Println("username:", template.HTMLEscapeString(username))
	fmt.Println("password:", template.HTMLEscapeString(password))
	template.HTMLEscape(w, []byte(username))
	fmt.Fprintln(w, json_encode(res))
}

/* 登录首页 */
func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	if r.Method == "GET" {
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println("token:", token)

		t, _ := template.ParseFiles("template/login.gtpl")
		t.Execute(w, token)
	} else {
		login_submit(w, r)
	}
}

/* 记事本接口 */
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

/* 插入数据 */
func add_user(db sql.DB, username string, password string) int64 {
	stmt, err := db.Prepare("INSERT userinfo SET username=?,password=?,created=?")
	checkErr(err)

	password = mymd5(password)
	res, err := stmt.Exec(username, password)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}

/* 获取用户信息 */
func get_user(db *sql.DB, username string) user_obj {
	var sql_str = "SELECT user_id, username, password, created FROM userinfo where username='" + username + "'"
	fmt.Println(sql_str)
	rows, err := db.Query(sql_str)
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

/* md5函数封装 */
func mymd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

/* MyReturn的json编码 */
func json_encode(obj MyReturn) string {
	body, err := json.Marshal(obj)

	if err != nil {
		panic(err.Error())
	}

	return string(body)
}

/* 异常流程 */
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/* 主线程 */
func main() {
	route()
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
