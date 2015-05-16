package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"log"
	"net/http"
)

/* json返回对象 */
type MyReturn struct {
	State   bool   `json:"state"`
	Message string `json:"err_msg"`
	Code    string `json:"err_code"`
}

type FILE_DATA struct {
	Name string
	Num  int
}

/* 路由 */
func route() {
	http.HandleFunc("/", index)
	http.HandleFunc("/imgload/", imgload)
}

/* 首页 */
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
}

/* 图片跳转 */
func imgload(w http.ResponseWriter, r *http.Request) {
	var res = MyReturn{State: false, Message: "", Code: "0"}
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			res.Message = fmt.Errorf("%v", err).Error()
			fmt.Println(res)
			fmt.Fprintln(w, json_encode(res))
			return
		} else {
			res.State = true
			fmt.Fprintln(w, json_encode(res))
			return
		}
	}()
	w.Header().Set("Content-type", "text/json; charset=utf-8")
	r.ParseForm()
	hash := get_param(r, "r")
	size := get_param(r, "i")
	uid := get_param(r, "u")
	sign := get_param(r, "s")
	encode := get_param(r, "e")
	static := get_param(r, "t")
	ttype := get_param(r, "p")
	fmt.Println(hash, size, uid, sign, encode, static, ttype)

	url := ""
	ttl := 0

	if hash != "" && uid != "" && sign != "" {
		real_uid := uid

		if encode != "" {
			real_uid = "get_dec_s"
		}

		check_sign := "make_imgload_sign"

		if check_sign != "" && check_sign == sign {

		}
	}

	if url == "" {
		url = "http://q.115.com/static/images/404.gif"
	}

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

/* 获取参数 */
func get_param(r *http.Request, column string) string {
	if len(r.Form[column]) > 0 {
		return r.Form[column][0]
	}

	return ""
}

/* MyReturn的json编码 */
func json_encode(obj interface{}) string {
	body, err := json.Marshal(obj)
	checkErr(err)
	return string(body)
}

/* MyReturn的json编码 */
func json_decode(str string) *simplejson.Json {
	js, err := simplejson.NewJson([]byte(str))
	checkErr(err)
	return js
}

/* 异常流程 */
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/* 主线程 */
func main() {
	fmt.Println("start:")
	route()
	err := http.ListenAndServe(":8889", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
