package main

// http://www.ttbiji.com/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/astaxie/goredis"
	_ "github.com/bitly/go-simplejson" // for json get
	"helper"
	"html/template"
	"io"
	"log"
	_ "net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/* aes加密的key */
var aes_key = "34jbfhg3gnhs90ds1gj1vhjfcsdf4sdv"
var md5_key = "dfgs435vh345"
var md6_key = "fgh52dgfd34f"
var redis_cli goredis.Client

/* IV */
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

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
	http.HandleFunc("/login/", login)
	http.HandleFunc("/register/", register)
	http.HandleFunc("/logout/", logout)
	http.HandleFunc("/note/", note)
	http.HandleFunc("/captcha/", captcha)
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
		username := get_currinfo(r)

		if username != "" {
			fmt.Fprintln(w, "<script type='text/javascript'>top.window.location.href='/note/';</script>")
			return
		}

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

	set_currinfo(w, username)
	fmt.Println("username:", template.HTMLEscapeString(username))
	fmt.Println("password:", template.HTMLEscapeString(password))
	//template.HTMLEscape(w, []byte(username))
	res.State = true
	fmt.Fprintln(w, json_encode(res))
}

/* 登录首页 */
func login(w http.ResponseWriter, r *http.Request) {
	username := get_currinfo(r)
	fmt.Println("lusername:" + username)

	if username != "" {
		fmt.Fprintln(w, "<script type='text/javascript'>top.window.location.href='/note/';</script>")
		return
	}

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

/* 退出登录 */
func logout(w http.ResponseWriter, r *http.Request) {
	set_currinfo(w, "")
	fmt.Fprintln(w, "<script type='text/javascript'>top.window.location.href='/login/';</script>")
}

/* 注册提交 */
func register_submit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/json; charset=utf-8")
	var res = MyReturn{State: false, Message: ""}
	r.ParseForm()
	username := r.Form["username"][0]
	password := r.Form["password"][0]
	//email := r.Form["email"][0]
	//mobile := r.Form["mobile"][0]
	//captcha := r.Form["captcha"][0]

	if len(username) < 6 || len(username) > 24 {
		res.Message = "用户名必须为6~24位字符串"
		fmt.Fprintln(w, json_encode(res))
		return
	}

	if len(password) < 6 {
		res.Message = "用户名必须为6位以上"
		fmt.Fprintln(w, json_encode(res))
		return
	}

	db, err := sql.Open("mysql", "root:d54321@/gonote?charset=utf8")
	checkErr(err)
	var user = get_user(db, username)

	if user.user_id != 0 {
		res.Message = "用户已存在，请输入其他用户名"
		fmt.Fprintln(w, json_encode(res))
		return
	}
}

/* 注册首页 */
func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	if r.Method == "GET" {
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println("token:", token)

		t, _ := template.ParseFiles("template/register.gtpl")
		t.Execute(w, token)
	} else {
		register_submit(w, r)
	}
}

/* 记事本接口 */
func note(w http.ResponseWriter, r *http.Request) {
	username := get_currinfo(r)
	fmt.Println("nusername:" + username)
	if username == "" {
		fmt.Fprintln(w, "<script type='text/javascript'>top.window.location.href='/login/';</script>")
		return
	}

	fmt.Println("path", r.URL.Path)
	if r.Method == "GET" {
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println("token:", token)

		t, _ := template.ParseFiles("template/note.gtpl")
		t.Execute(w, token)
	} else {
		register_submit(w, r)
	}
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

/* 获取当前用户信息 */
func get_currinfo(r *http.Request) string {
	var username, _ = r.Cookie("username")
	var userkey, _ = r.Cookie("userkey")

	if username == nil || userkey == nil {
		return ""
	}

	fmt.Println("username:" + username.Value)
	fmt.Println("userkey:" + userkey.Value)
	ukey := cookie_arithmetic(username.Value)

	if ukey != userkey.Value {
		return ""
	}

	return username.Value
}

/* 设置当前用户信息 */
func set_currinfo(w http.ResponseWriter, username string) {
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	cookie := http.Cookie{Name: "username", Path: "/", Value: username, Expires: expiration, HttpOnly: true}
	http.SetCookie(w, &cookie)
	key := cookie_arithmetic(username)
	fmt.Println("login:" + username)
	fmt.Println("key:" + key)
	cookie1 := http.Cookie{Name: "userkey", Path: "/", Value: key, Expires: expiration, HttpOnly: true}
	http.SetCookie(w, &cookie1)
}

/* 设置当前用户信息 */
func cookie_arithmetic(str string) string {
	res := mymd5(mymd5(str+md5_key) + md6_key)
	return res
}

/* md5函数封装 */
func mymd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

/* aes加密 */
func aes_encode(str string) string {
	plaintext := []byte(str)
	c, err := aes.NewCipher([]byte(aes_key))

	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(aes_key), err)
		os.Exit(-1)
	}

	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	fmt.Printf("%s=>%x\n", plaintext, ciphertext)
	return string(ciphertext)
}

/* aes解密密 */
func aes_decode(str string) string {
	ciphertext := []byte(str)
	c, err := aes.NewCipher([]byte(aes_key))

	if err != nil {
		fmt.Printf("Error: NewCipher(%d bytes) = %s", len(aes_key), err)
		os.Exit(-1)
	}

	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(ciphertext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	fmt.Printf("%x=>%s\n", ciphertext, plaintextCopy)
	return string(plaintextCopy)
}

/* MyReturn的json编码 */
func json_encode(obj MyReturn) string {
	body, err := json.Marshal(obj)

	if err != nil {
		panic(err.Error())
	}

	return string(body)
}

/* 设置redis缓存 */
func SR(key string, value string) {
	redis_cli.Set(key, []byte(value))
}

/* 获取redis缓存 */
func GR(key string) string {
	res, _ := redis_cli.Get(key)
	return string(res)
}

/* 删除redis缓存 */
func DR(key string) {
	redis_cli.Del(key)
}

/* 异常流程 */
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/* 验证码 */
func captcha(w http.ResponseWriter, req *http.Request) {
	d := make([]byte, 4)
	s := helper.NewLen(4)
	ss := ""
	d = []byte(s)
	for v := range d {
		d[v] %= 10
		ss += strconv.FormatInt(int64(d[v]), 32)
	}
	w.Header().Set("Content-Type", "image/png")
	SR("C_100", string(d))
	helper.NewImage(d, 100, 40).WriteTo(w)
	fmt.Println(ss)
}

/* 主线程 */
func main() {
	route()
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
