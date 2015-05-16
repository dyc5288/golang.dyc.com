package main

import (
	_ "bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	_ "database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/astaxie/goredis"
	"github.com/bitly/go-simplejson"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/bradleyg/go-address"
	_ "helper"
	_ "html/template"
	_ "io"
	_ "io/ioutil"
	"log"
	"math/rand"
	_ "net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	http.HandleFunc("/imgload/", report)
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

	if r.Method != "POST" {
		panic("请求不正确")
	}

	data := r.FormValue("data")
	btype := r.FormValue("btype")
	ptype := r.FormValue("ptype")

	if len(data) == 0 && len(btype) == 0 {
		panic("参数错误")
	}

	rptype := "PASSPORT"

	if len(ptype) != 0 {
		rptype = ptype
	}

	rdata := data
	rbtype := btype
	conf, _ := RDATA[rbtype]
	ip := get_client_ip(r)

	if len(conf) == 0 {
		panic("不存在该btype")
	}

	check_queue := map[string]string{}
	filename := get_filename(ip, rbtype, rptype)
	//fmt.Println(filename)

	if _, err := os.Stat(filename); err != nil {
		panic(filename + "不存在")
	}

	sjson := json_decode(rdata)
	map_data, _ := sjson.Map()
	content := format_data(rbtype, map_data, (&check_queue))
	write_data(filename, content)

	for k, _ := range check_queue {
		SM(k, "1")
	}
	//fmt.Println("write success")
}

/* 写数据 */
func write_data(filename string, content string) int {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	res, err := file.WriteString(content + "\n")

	if err != nil {
		panic(err)
	}

	return res
}

/* 字符串截取函数 */
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

/* 获取IP地址 */
func get_client_ip(r *http.Request) string {
	address, err := goaddress.Get(r, nil)

	if err != nil {
		panic(err)
	}

	return address
}

/* md5函数封装 */
func gomd5(str string) string {
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

/* memcache初始化172.16.0.118 */
func init_memcache() *memcache.Client {
	mc_cache := memcache.New("10.10.2.7:11211")
	return mc_cache
}

/* 设置memcache缓存 */
func SM(key string, value string) bool {
	mc_cache := init_memcache()
	key = key + memcache_prefix
	err := mc_cache.Set(&memcache.Item{Key: key, Value: []byte(value)})

	if err != nil {
		return false
	}

	return true
}

/* 获取memcache缓存 */
func GM(key string) string {
	mc_cache := init_memcache()
	key = key + memcache_prefix
	res, err := mc_cache.Get(key)

	if err != nil {
		fmt.Println(mc_cache)
		fmt.Println(key)
		fmt.Println(res)
		fmt.Println(err)
		return ""
	}

	return string(res.Value)
}

/* 删除memcache缓存 */
func DM(key string) bool {
	mc_cache := init_memcache()
	key = key + memcache_prefix
	err := mc_cache.Delete(key)

	if err != nil {
		return false
	}

	return true
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
