package main

// http://www.ttbiji.com/

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
	"github.com/astaxie/goredis"
	"github.com/bitly/go-simplejson" // for json get
	_ "helper"
	_ "html/template"
	_ "io"
	_ "io/ioutil"
	"log"
	_ "net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	_ "time"
)

/* aes加密的key */
var aes_key = "34jbfhg3gnhs90ds1gj1vhjfcsdf4sdv"
var md5_key = "dfgs435vh345"
var md6_key = "fgh52dgfd34f"
var redis_cli goredis.Client
var bcp_rows = 4500
var RDATA = map[string]map[string]map[string]string{
	"register": {
		"ICP_CODE":             {"name": "ICP编码", "length": "14", "index": "1"},
		"DATA_TYPE":            {"name": "协议类型标识", "length": "10", "index": "2"},
		"USER_ID":              {"name": "用户ID", "length": "30", "index": "3", "unique": "1"},
		"USER_NAME":            {"name": "用户账号", "length": "30", "index": "4"},
		"PASSWORD":             {"name": "密码", "length": "128", "index": "5"},
		"NICK_NAME":            {"name": "昵称", "length": "60", "index": "6"},
		"STATUS":               {"name": "状态", "length": "10", "index": "7"},
		"REAL_NAME":            {"name": "姓名", "length": "30", "index": "8"},
		"SEX":                  {"name": "性别", "length": "10", "index": "9"},
		"BIRTHDAY":             {"name": "生日", "length": "14", "index": "10"},
		"CONTACT_TEL":          {"name": "联系电话", "length": "64", "index": "11"},
		"CERTIFICATE_TYPE":     {"name": "证件类型", "length": "32", "index": "12"},
		"CERTIFICATE_CODE":     {"name": "证件号码", "length": "60", "index": "13"},
		"BIND_TEL":             {"name": "绑定手机号码", "length": "21", "index": "14"},
		"BIND_QQ":              {"name": "绑定QQ号码", "length": "15", "index": "15"},
		"BIND_MSN":             {"name": "绑定MSN号码", "length": "50", "index": "16"},
		"EMAIL":                {"name": "电子邮件地址", "length": "50", "index": "17"},
		"REGISTER_TIME":        {"name": "注册时间", "length": "14", "index": "18"},
		"LAST_LOGIN_TIME":      {"name": "最后登录时间", "length": "14", "index": "19"},
		"LAST_CHANGE_PASSWORD": {"name": "最后更改密码时间", "length": "14", "index": "20"},
		"LAST_MODIFY_TIME":     {"name": "最后修改资料时间", "length": "14", "index": "21"},
		"REGISTER_IP":          {"name": "注册/修改IP", "length": "128", "index": "22"},
		"REGISTER_PORT":        {"name": "注册/修改端口", "length": "128", "index": "23"},
		"REGISTER_MAC":         {"name": "注册/修改Mac", "length": "17", "index": "24"},
		"REGISTER_BIOS_ID":     {"name": "注册硬件编码", "length": "120", "index": "25"},
		"PROVINCE":             {"name": "省份", "length": "10", "index": "26"},
		"CITY":                 {"name": "城市", "length": "10", "index": "27"},
		"ADDRESS":              {"name": "地址", "length": "100", "index": "28"}}}

/* IV */
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

/* json返回对象 */
type MyReturn struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
}

/* 路由 */
func route() {
	http.HandleFunc("/", index)
	http.HandleFunc("/report/", report)
}

/* 首页 */
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
}

/* 上报 */
func report(w http.ResponseWriter, r *http.Request) {
	var res = MyReturn{State: false, Message: ""}
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
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
	w.Header().Set("Content-type", "text/json; charset=utf-8")
	r.ParseForm()
	data, _ := r.Form["data"]
	btype, _ := r.Form["btype"]

	if len(data) == 0 && len(btype) == 0 {
		panic("参数错误")
	}

	rdata := data[0]
	rbtype := btype[0]
	conf, _ := RDATA[rbtype]

	if len(conf) == 0 {
		panic("不存在该btype")
	}

	filename := get_filename()
	sjson := json_decode(rdata)
	map_data, _ := sjson.Map()
	content := format_data(rbtype, map_data)
	//fmt.Println(content)
	write_data(filename, content)
}

/* 写数据 */
func write_data(filename string, content string) int {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0x644)
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

func get_filename() string {
	return "E:/git/golang.dyc.com/src/icp/test2.txt"
}

/* 格式化 */
func format_data(btype string, data map[string]interface{}) string {
	conf, _ := RDATA[btype]

	if len(conf) == 0 {
		panic("不存在该btype")
	}

	res := []string{}
	rindex := []int{}
	reg, err := regexp.Compile("[\r\t\n]")

	if err != nil {
		panic("正则错误")
	}

	for k, c := range conf {
		index, _ := c["index"]

		if len(index) == 0 {
			panic(k + "的length或index属性不存在")
		}

		ind, err := strconv.Atoi(index)

		if err != nil {
			panic(k + "的index属性不是整形数")
		}

		v, ok := data[k]
		val := ""

		if ok {
			l, _ := c["length"]

			if len(l) == 0 {
				panic(k + "的length或index属性不存在")
			}

			le, err := strconv.Atoi(l)

			if err != nil {
				panic(k + "的length属性不是整形数")
			}

			switch v.(type) {
			case string:
				val = v.(string)
			case json.Number:
				val = (string)(v.(json.Number))
			default:
				panic("data的格式err")
			}

			if len(val) > le {
				val = Substr(val, 0, le)
			}

			val = reg.ReplaceAllString(val, " ")
		}

		rindex = append(rindex, ind-1)
		res = append(res, val)
	}

	zres := []string{}

	for i := 0; i < len(rindex); i++ {
		for j := 0; j < len(rindex); j++ {
			if i == rindex[j] {
				zres = append(zres, res[j])
			}
		}
	}

	return strings.Join(zres, "\t")
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
	checkErr(err)
	return string(body)
}

/* MyReturn的json编码 */
func json_decode(str string) *simplejson.Json {
	js, err := simplejson.NewJson([]byte(str))
	checkErr(err)
	return js
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

/* 主线程 */
func main() {
	route()
	err := http.ListenAndServe(":8889", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("start:")
	}
}
