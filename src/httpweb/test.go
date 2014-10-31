package main

// http://www.ttbiji.com/

import (
	_ "crypto/md5"
	_ "database/sql"
	_ "encoding/hex"
	"encoding/json"
	_ "fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/bitly/go-simplejson" // for json get
	_ "html/template"
	_ "io"
	_ "log"
	_ "net"
	_ "net/http"
	"os"
	_ "strconv"
	_ "strings"
	_ "time"
)

type Server struct {
	// ID 不会导出到JSON中
	ID string `json:"ID"`

	// ServerName 的值会进行二次JSON编码
	ServerName  string `json:"serverName"`
	ServerName2 string `json:"serverName2,string"`

	// 如果 ServerIP 为空，则不输出到JSON串中
	ServerIP string `json:"serverIP,omitempty"`
}

func maint() {
	s := Server{
		ID:          "3",
		ServerName:  `Go "1.0" `,
		ServerName2: `Go "1.0" `,
		ServerIP:    `sdff`,
	}

	b, _ := json.Marshal(s)
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))
}
