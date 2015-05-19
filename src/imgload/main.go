package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/couchbase/go-couchbase"
	"log"
	"math"
	"net/http"
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

var DEBUG_LEVEl = false

var _super2dec_arr = map[string]int{
	"H": 0, "2": 1, "t": 2, "O": 3, "u": 4, "z": 5, "b": 6, "F": 7, "P": 8,
	"V": 9, "E": 10, "3": 11, "8": 12, "5": 13, "x": 14, "f": 15, "X": 16,
	"c": 17, "k": 18, "Z": 19, "A": 20, "U": 21, "B": 22, "h": 23, "Y": 24,
	"D": 25, "n": 26, "N": 27, "7": 28, "I": 29, "v": 30, "i": 31, "p": 32,
	"T": 33, "L": 34, "a": 35, "0": 36, "6": 37, "M": 38, "m": 39, "q": 40,
	"y": 41, "J": 42, "j": 43, "s": 44, "r": 45, "4": 46, "o": 47, "G": 48,
	"Q": 49, "1": 50, "K": 51, "9": 52, "S": 53, "l": 54, "e": 55, "g": 56,
	"d": 57, "w": 58, "C": 59, "R": 60, "W": 61}

var couchbase_server = map[string]map[int]map[string]string{
	"imageinfo": {
		0: {"ip": "10.10.0.216", "port": "8091", "bucket": "imageinfo", "prefix": "i"},
		1: {"ip": "10.10.0.217", "port": "8091", "bucket": "imageinfo", "prefix": "i"},
		2: {"ip": "10.10.0.218", "port": "8091", "bucket": "imageinfo", "prefix": "i"},
		3: {"ip": "10.10.0.219", "port": "8091", "bucket": "imageinfo", "prefix": "i"},
		4: {"ip": "10.10.0.220", "port": "8091", "bucket": "imageinfo", "prefix": "i"},
		5: {"ip": "10.10.0.221", "port": "8091", "bucket": "imageinfo", "prefix": "i"}}}

/* 初始化配置 */
func init_config() {
	if DEBUG_LEVEl {
		couchbase_server = map[string]map[int]map[string]string{
			"imageinfo": {
				0: {"ip": "172.16.0.120", "port": "8091", "bucket": "imageinfo", "prefix": "i"},
				1: {"ip": "172.16.0.121", "port": "8091", "bucket": "imageinfo", "prefix": "i"},
				2: {"ip": "172.16.0.122", "port": "8091", "bucket": "imageinfo", "prefix": "i"}}}

	}
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

	if size == "100" {
		size = "200"
	}

	if hash != "" && uid != "" && sign != "" {
		real_uid := uid

		if encode != "" {
			real_uid = strconv.Itoa(get_dec_s(uid))
		}

		iencode, _ := strconv.Atoi(encode)
		suid := real_uid

		if iencode > 1 {
			suid = uid
		}

		encode_level := "0"

		if iencode >= 5 {
			encode_level = encode
		}

		check_sign := make_imgload_sign(hash, suid, encode_level)

		if check_sign != "" && check_sign == sign {
			client := agent2code(r)
			fmt.Println("client:", client)
			mclients := map[int]bool{0: true, 3: true, 8: true, 9: true, 10: true}
			is_mobile, err := mclients[client]
			fmt.Println("is_mobile:", is_mobile)

			if err == false {
				is_mobile = false
			}

			is_mobile = !is_mobile
			fmt.Println("size:", size)
			if size == "" || size == "0" {
				if !is_mobile {
					size = "800"
				} else {
					size = "480"
					ttime := time.Date(time.Now().Year(), time.Now().Month()+1, 31, 0, 0, 0, 0, time.UTC)
					ttl, _ = strconv.Atoi(strconv.FormatInt(ttime.Unix(), 10))
				}
			}
			fmt.Println("size1:", size)
			sizes := map[string]bool{"100": true, "200": true, "320": true, "480": true, "800": true, "1440": true}
			is_size, err := sizes[size]

			if err == false {
				is_size = false
			}

			if !is_size {
				size = "800"
			}
			fmt.Println("size2:", size)
			static_gif := static != "" || is_mobile

			if static_gif {
				if ttype == "" {
					ttype = get_ttype(hash)
				}

				if ttype == "" || ttype != "gif" {
					static_gif = false
				}
			}

			url = get_pic_url(hash, size, ttl, static_gif)
		}
	}

	if url == "" {
		url = "http://q.115.com/static/images/404.gif"
	}

	//fmt.Fprintln(w, url)
	res.Message = url
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

/* 获取参数 */
func get_param(r *http.Request, column string) string {
	if len(r.Form[column]) > 0 {
		return r.Form[column][0]
	}

	return ""
}

/* 生成sign */
func make_imgload_sign(hash string, uid string, encode_level string) string {
	hash = strings.Trim(hash, " ")
	uid = strings.Trim(uid, " ")
	is_match, _ := regexp.MatchString("[A-F0-9]{40}", hash)

	if hash == "" || !is_match || uid == "" {
		return ""
	}

	rstr := hash + string(uid[0]) + "@q.115" + uid + string(hash[0]) + "@u.img"
	iencode_level, _ := strconv.Atoi(encode_level)

	if iencode_level > 0 {
		rstr += encode_level
	}

	res := base64_encode(gomd5(rstr, true))
	res = strings.Replace(res, "/", "-", -1)
	res = strings.Replace(res, "+", ".", -1)
	res = strings.Replace(res, "=", "", -1)
	return res
}

/* 获取类型 */
func get_ttype(hash string) string {
	//bucket := connect_couchbase("imageinfo")
	return ""
}

/* 获取网盘缩略图地址 */
func get_pic_url(hash string, stype string, cache_time int, static_gif bool) string {
	key := "654321"
	per := "/thumb/"
	store_path := ""
	fmt.Println("stype:", stype)

	switch stype {
	case "100":
		fallthrough
	case "200":
		fallthrough
	case "480":
		fallthrough
	case "800":
		width := stype

		if static_gif && (stype == "100" || stype == "200" || stype == "480") {
			width += "s"
		}
		store_path = Substr(hash, 0, 1) + "/" + Substr(hash, 1, 2) + "/" + Substr(hash, 3, 2) + "/" + hash + "_" + width + "_" + stype
	case "320":
		store_path = Substr(hash, 0, 1) + "/" + Substr(hash, 1, 2) + "/" + Substr(hash, 3, 2) + "/" + hash + "_" + stype + "_960"
	case "1440":
		store_path = Substr(hash, 0, 1) + "/" + Substr(hash, 1, 2) + "/" + Substr(hash, 3, 2) + "/" + hash + "_" + stype + "_900"
	default:
		return ""
	}

	if cache_time == 0 {
		cache_time, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	}

	fmt.Println("cache_time:", cache_time)
	//fmt.Println("key:", key, ", per:", per, ", store_path:", store_path, ",cache_time:", cache_time)
	//fmt.Println("md5:", gomd5(key+per+store_path+strconv.Itoa(cache_time), true))
	sign := base64_encode(gomd5(key+per+store_path+strconv.Itoa(cache_time), true))
	fmt.Println("base64:", sign)
	sign = strings.Replace(sign, "+", "-", -1)
	sign = strings.Replace(sign, "/", "_", -1)
	sign = strings.Replace(sign, "=", "", -1)
	fmt.Println("sign:", sign)

	host := "http://thumb.115.com"
	fmt.Println("store_path:", store_path)
	return host + per + store_path + "?s=" + sign + "&t=" + strconv.Itoa(cache_time) + "&sync=1"
}

/* 超进位数转换为十进位数（加密模式--超数第一位为进位标识） */
func get_dec_s(number string) int {
	dd := string(number[0])
	_digit := _super2dec_base(dd)

	if _digit == -1 {
		return 0
	}

	_digit = _digit + 1
	re, _ := regexp.Compile(dd)
	number = re.ReplaceAllString(number, "")
	return get_dec(number, _digit)
}

/* 超进位数转换为十进位数 */
func get_dec(number string, _digit int) int {
	renum := 0
	nlength := len(number) - 1
	j := 0
	digit := 0
	dd := ""

	for i := nlength; i > -1; i-- {
		dd = (string)(number[i])
		digit = _super2dec_base(dd)

		if digit == -1 {
			return 0
		}

		renum = renum + digit*int(math.Pow(float64(_digit), float64(j)))
		j++
	}

	return renum
}

/* 超进位数和十进位数的个位对照数 */
func _super2dec_base(number string) int {
	index, err := _super2dec_arr[number]

	if err == false {
		return -1
	}

	return index
}

/* md5函数封装 */
func gomd5(str string, raw bool) string {
	m := md5.New()
	m.Write([]byte(str))

	if raw {
		return string(m.Sum(nil))
	}

	return hex.EncodeToString(m.Sum(nil))
}

/* 获取客户端代号 */
func agent2code(r *http.Request) int {
	code := 0
	agent := strings.ToLower(r.UserAgent())

	if agent != "" {
		if strings.Contains(agent, "iphone") || strings.Contains(agent, "udown") {
			code = 2 // from iPhone
		} else if strings.Contains(agent, "ipad") {
			code = 3 // from iPad
		} else if strings.Contains(agent, "android") || strings.Contains(agent, "115disk") {
			code = 4 // from Android device
		} else if strings.Contains(agent, "windows phone") {
			code = 5 // from Windows Phone
		} else if strings.Contains(agent, "mac os") {
			code = 8 // from MAC OS X
		} else if strings.Contains(agent, "115chrome") {
			code = 9 // from 115 chrome browse
		} else if strings.Contains(agent, "startos") {
			code = 10 // from Start OS
		} else {
			if strings.Contains(agent, "mobi") {
				code = 6 //from other mobile device
			} else {
				fmt.Println(r.Header)
			}
		}
	}

	return code
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

/* couchbase连接 */
func connect_couchbase(key string) *couchbase.Bucket {
	conf, err := couchbase_server[key]

	if err == false {
		return nil
	}

	for _, server := range conf {
		u := "http://" + server["ip"] + ":" + server["port"] + "/"
		fmt.Println("c:", u)
		client, err := couchbase.Connect(u)

		if err != nil {
			fmt.Printf("Connect failed %v", err)
			continue
		}

		cbpool, err := client.GetPool("default")
		if err != nil {
			fmt.Printf("Failed to connect to default pool %v", err)
			continue
		}

		cbbucket, err := cbpool.GetBucket(server["bucket"])

		if err != nil {
			fmt.Printf("Failed to connect to bucket %s %v", server["bucket"], err)
			continue
		}

		return cbbucket
	}

	return nil
}

/* base64编码*/
func base64_encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func base64_decode(str string) string {
	res, _ := base64.StdEncoding.DecodeString(str)
	return string(res)
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
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " server start:")
	route()
	init_config()
	bucket := connect_couchbase("imageinfo")
	//bucket.Set("i_test", 0, []string{"an", "example", "list"})
	ob := map[string]interface{}{}
	err1 := bucket.Get("i_ACC9D5333CC559EF05F4DFB39446C29EC9320735", &ob)
	fmt.Println(err1, "|||||||||", ob)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " server listen 8889:")
	err := http.ListenAndServe(":8889", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
