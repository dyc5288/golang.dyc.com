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

/* aes加密的key */
var aes_key = "34jbfhg3gnhs90ds1gj1vhjfcsdf4sdv"
var md5_key = "dfgs435vh345"
var md6_key = "fgh52dgfd34f"
var SCONFIG = map[string]string{
	"name":    "ICP",
	"version": "2.1",
	"passwd":  "run1234!@#",
}
var memcache_prefix = "y"
var bcp_rows = 4500
var ICP_DATA = "/www/web/web/icp/data/notsync"
var RDATA = map[string]map[string]map[string]string{
	"register": {
		"ICP_CODE":             {"name": "ICP编码", "length": "14", "index": "1", "default": "44190013100300"},
		"DATA_TYPE":            {"name": "协议类型标识", "length": "10", "index": "2", "default": "PASSPORT"},
		"USER_ID":              {"name": "用户ID", "length": "30", "index": "3", "unique": "1"},
		"USER_NAME":            {"name": "用户账号", "length": "30", "index": "4"},
		"PASSWORD":             {"name": "密码", "length": "128", "index": "5"},
		"NICK_NAME":            {"name": "昵称", "length": "60", "index": "6"},
		"STATUS":               {"name": "状态", "length": "10", "index": "7", "default": "新建"},
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
		"REGISTER_TIME":        {"name": "注册时间", "length": "14", "index": "18", "format": "date"},
		"LAST_LOGIN_TIME":      {"name": "最后登录时间", "length": "14", "index": "19", "format": "date"},
		"LAST_CHANGE_PASSWORD": {"name": "最后更改密码时间", "length": "14", "index": "20"},
		"LAST_MODIFY_TIME":     {"name": "最后修改资料时间", "length": "14", "index": "21", "format": "date"},
		"REGISTER_IP":          {"name": "注册/修改IP", "length": "128", "index": "22"},
		"LAST_LOGIN_IP":        {"name": "最后登录IP", "length": "128", "index": "23"},
		"REGISTER_MAC":         {"name": "注册/修改Mac", "length": "17", "index": "24"},
		"LAST_LOGIN_MAC":       {"name": "最后登录MAC", "length": "17", "index": "25"},
		"PROVINCE":             {"name": "省份", "length": "10", "index": "26"},
		"CITY":                 {"name": "城市", "length": "10", "index": "27"},
		"ADDRESS":              {"name": "地址", "length": "100", "index": "28"}},
	"login": {
		"ICP_CODE":    {"name": "ICP编码", "length": "14", "index": "1", "default": "44190013100300"},
		"DATA_TYPE":   {"name": "协议类型标识", "length": "10", "index": "2", "default": "PASSPORT"},
		"SRC_IP":      {"name": "源IP", "length": "128", "index": "3"},
		"SRC_PORT":    {"name": "源端口", "length": "128", "index": "4"},
		"DST_IP":      {"name": "目的IP", "length": "128", "index": "5"},
		"DST_PORT":    {"name": "目的端口", "length": "128", "index": "6"},
		"USER_ID":     {"name": "用户ID", "length": "30", "index": "7"},
		"USER_NAME":   {"name": "用户账号", "length": "30", "index": "8"},
		"NICK_NAME":   {"name": "昵称", "length": "30", "index": "9"},
		"PASSWORD":    {"name": "用户密码", "length": "128", "index": "10"},
		"MAC_ADDRESS": {"name": "源MAC 地址", "length": "64", "index": "11"},
		"INNER_IP":    {"name": "内网IP 地址", "length": "32", "index": "12"},
		"ACTION_TIME": {"name": "发生时间", "length": "60", "index": "13", "format": "date"},
		"ACTION":      {"name": "动作", "length": "21", "index": "14"}}}

/* IV */
var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

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
	http.HandleFunc("/report/", report)
}

/* 首页 */
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
}

/* 上报 */
func report(w http.ResponseWriter, r *http.Request) {
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

/* 获取文件名 */
func get_filename(ip string, btype string, ptype string) string {
	cache_key := btype + "_" + ip
	file_data := GM(cache_key)
	fmt.Println(file_data)

	if file_data != "" {
		sjson := json_decode(file_data)
		file_arr, _ := sjson.Map()
		filename, ok := file_arr["name"]
		file_num, ok2 := file_arr["num"]

		if ok && ok2 {
			fname := filename.(string)

			if _, err := os.Stat(fname); err == nil {
				fnum, _ := strconv.Atoi(file_num.(string))

				if fnum <= bcp_rows {
					save_filename(cache_key, fname, fnum+1)
					return fname
				}
			}
		}
	}

	filename := create_file(ip, btype, ptype)

	fmt.Println("te:" + filename)
	if filename == "" {
		panic("创建文件失败")
	}

	header := SCONFIG["name"] + "\t" + SCONFIG["version"] + "\t" + ptype
	write_data(filename, header)
	save_filename(cache_key, filename, 1)
	return filename
}

/* 保存缓存 */
func save_filename(cache_key string, filename string, fnum int) {
	//fmt.Println(fnum)
	cache := make(map[string]string)
	cache["name"] = filename
	cache["num"] = strconv.Itoa(fnum)
	SM(cache_key, json_encode(cache))
}

func create_file(ip string, btype string, ptype string) string {
	ctime := strconv.FormatInt(time.Now().Unix(), 10)
	file_path := ICP_DATA + "/data/" + ip + "/" + btype + "/" + ctime + "/" + ptype
	err := os.MkdirAll(file_path, 0777)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	dir := [...]string{"data", ip, btype, ctime, ptype}
	tmp_dir := ICP_DATA

	for index := 0; index < len(dir); index++ {
		tmp_dir = tmp_dir + "/" + dir[index]
		err = os.Chmod(tmp_dir, 0777)

		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand := fmt.Sprintf("%06d", r.Intn(1000000))
	file_name := btype + "_" + ptype + "_" + ctime + "_" + rand + ".bcp"
	filename := file_path + "/" + file_name
	return filename
}

/* 格式化 */
func format_data(btype string, data map[string]interface{}, check_queue *map[string]string) string {
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
				fmt.Println(k)
				fmt.Println(v)
				panic("data的格式err")
			}

			if len(val) > le {
				val = Substr(val, 0, le)
			}

			format, _ := c["format"]

			switch format {
			case "date":
				if len(val) != 0 {
					t, err := strconv.ParseInt(val, 10, 64)

					if err != nil {
						fmt.Println(err)
						fmt.Println(val)
						panic(k + "的属性不是date类型")
					}

					val = time.Unix(t, 0).Format("20060102150405")
				}
			default:
			}

			uq, _ := c["unique"]

			if len(uq) != 0 {
				check_unique(btype, val, check_queue)
			}

			val = reg.ReplaceAllString(val, " ")
		} else {

			def, _ := c["default"]

			if len(def) != 0 {
				val = def
			}
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

/* 检查唯一性 */
func check_unique(btype string, value string, check_queue *map[string]string) {
	cache_key := btype + "_" + value
	cache := GM(cache_key)

	if cache != "" {
		panic("已存在" + value)
	}

	(*check_queue)[cache_key] = "1"
	//SM(cache_key, "1")
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
	mc_cache := memcache.New("10.11.2.7:11211")
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
	SM("test", "test")
	d := GM("test")
	fmt.Println(d)
	fmt.Println("start listen 8809:")
	err := http.ListenAndServe(":8889", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
