package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/* json返回对象 */
type MyReturn struct {
	State   bool     `json:"state"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

/* 路由 */
func route() {
	http.HandleFunc("/", index)
	http.HandleFunc("/log/", log_request)
}

/* 首页 */
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
}

/* 日志查询 */
func log_request(w http.ResponseWriter, r *http.Request) {
	var res = MyReturn{State: false, Message: ""}
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		w.Header().Set("Content-type", "text/json; charset=utf-8")
		if err := recover(); err != nil {
			res.Message = fmt.Errorf("%v", err).Error()
			fmt.Fprintln(w, json_encode(res))
			return
		} else {
			res.State = true
			fmt.Fprintln(w, json_encode(res))
			return
		}
	}()
	r.ParseForm()
	path := get_params(r, "path", "")

	if path == "" {
		panic("参数错误")
	}

	keyword := get_params(r, "keyword", "")
	start, err1 := strconv.Atoi(get_params(r, "start", "0"))
	limit, err2 := strconv.Atoi(get_params(r, "limit", "20"))

	if err1 != nil || err2 != nil {
		panic("start或limit不是数字")
	}

	result := search_log(path, keyword, start, limit)
	res.Data = result
	//fmt.Println(result)
}

/* 搜索 */
func search_log(file string, keyword string, start int, limit int) []string {
	if _, err := os.Stat(file); err != nil {
		panic(file + "文件不存在")
	}

	res := []string{}
	f, err := os.Open(file)
	defer f.Close()
	line_num := 0

	if nil == err {
		buff := bufio.NewReader(f)

		for {
			line, err := buff.ReadString('\n')

			if err != nil || io.EOF == err {
				break
			}

			if keyword != "" && strings.Contains(line, keyword) == false {
				continue
			}

			//fmt.Println("[" + line + "]")
			line_num++

			if line_num >= start && line_num < start+limit {
				res = append(res, line)
			}
		}
	}

	return res
}

/* 接收参数 */
func get_params(r *http.Request, name string, default_value string) string {
	res, _ := r.Form[name]
	result := default_value

	if len(res) != 0 {
		result = res[0]
	}

	return result
}

/* MyReturn的json编码 */
func json_encode(obj MyReturn) string {
	body, err := json.Marshal(obj)
	checkErr(err)
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
	fmt.Println("start listen 8810:")
	err := http.ListenAndServe(":8810", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
