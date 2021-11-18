package fmdu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// AppPath ... アプリケーションパス
const AppPath = "/api/v1/"

// API Method名の定義
const (
	// ログイン/ログアウト関連
	APIMethodLogin  = "login"
	APIMethodLogout = "logout"

	// マスタ関連
	APIMethodFtype       = "ftype"
	APIMethodCompany     = "company"
	APIMethodIntroStatus = "introstatus"
	APIMethodArea        = "area"

	APIMethodCustomer = "cust"

	APIMethodFmdu    = "fmdu"
	APIMethodIntro   = "intro"
	APIMethodDensyou = "densyou"

	APIMethodCpl      = "cpl"
	APIMethodAttached = "attached"

	APIMethodImport    = "import"
	APIMethodImportLog = "importlog"

	APIMethodBackup = "backup"

	APIMethodLog = "log"
)

func getNoLogMethod() (m map[string]interface{}) {
	m = make(map[string]interface{})
	m[APIMethodFtype] = nil
	m[APIMethodCompany] = nil
	m[APIMethodIntroStatus] = nil
	m[APIMethodArea] = nil
	return
}

// ParameterAPI ... WEB API用パラメータ
type ParameterAPI struct {
	*ParameterRDB
	Kei           uint8
	LogPath       string
	AppPath       string
	SmsPath       string
	CsvPath       string
	CsvMdu        string
	CsvCust       string
	CrawlLogPath  string
	BackupPath    string
	BackupLogPath string
}

// クライアントに処理結果を返す
func (p *Parameter) retunClientHandler(c echo.Context, rec interface{}, length int, err error, dur time.Duration) error {
	var status uint8 = 1
	if err != nil {
		status = 0
	}
	j, _ := GetAPIJSONString(rec, length, err, dur)
	p.writeLog(c, status)
	return c.String(http.StatusOK, j)
}

// ログ書き込み
func (p *Parameter) writeLog(c echo.Context, status uint8) {
	file := p.LogPath

	r := c.Request()
	met := r.Method
	sv := getServiceName(c.Path())
	/*
		if (met == "GET") && (sv == APIMethodFigure) {
			return
		}
	*/
	mmap := getNoLogMethod()
	_, ok := mmap[sv]
	if ok {
		return
	}

	logfile, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//const layout = "2006-01-02 15:04:05"

	msg := fmt.Sprintf("|%v|%v|%v|%v|%v|%v|%v", status, met, sv, p.UserID, p.IP, p.Binfo.Platform, p.Binfo.UserAgent)
	if err != nil {
		fmt.Printf("cannnot open %s: %s", file, msg)
	}
	defer logfile.Close()

	// io.MultiWriteで、
	// 標準出力とファイルの両方を束ねて、
	// logの出力先に設定する
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	//log.SetFlags(log.Ldate | log.Ltime)
	log.Println(msg)
}

func getServiceName(path string) (s string) {
	s = path
	s = strings.Replace(s, AppPath, "", -1)
	return
}

// APIJSON ... API返却値のJSON型
type APIJSON struct {
	Length   int
	Error    string
	Record   interface{}
	Duration string
}

// NewAPIJSON ... API返却値のJSON型生成
func NewAPIJSON(rec interface{}, len int, er error, dur time.Duration) *APIJSON {
	errstr := ""
	if er != nil {
		errstr = er.Error()
	}
	aj := new(APIJSON)
	aj.Length = len
	aj.Error = errstr
	aj.Record = rec
	aj.Duration = dur.String()

	return aj
}

// GetAPIJSONString ... JSON文字列を返す
func GetAPIJSONString(rec interface{}, length int, er error, dur time.Duration) (jsonstr string, err error) {
	jsonBytes, err := json.Marshal(NewAPIJSON(rec, length, er, dur))

	if err != nil {
		//log.Fatal("JSON Marshal error:", err)
		return
	}
	jsonstr = string(jsonBytes)
	return
}

// GetAPIJSONStringIndent ... JSON整形文字列を返す
func GetAPIJSONStringIndent(rec interface{}, length int, er error, dur time.Duration) (jsonstr string, err error) {
	jsonBytes, err := json.Marshal(NewAPIJSON(rec, length, er, dur))
	if err != nil {
		//log.Fatal("JSON Marshal error:", err)
		return
	}
	buf := new(bytes.Buffer)
	err = json.Indent(buf, jsonBytes, "", "	")
	if err != nil {
		return
	}
	jsonstr = buf.String()
	return
}

// GetJSONStringIndent ... JSON整形文字列を返す
func GetJSONStringIndent(rec interface{}) (jsonstr string, err error) {
	jsonBytes, err := json.Marshal(rec)
	if err != nil {
		//log.Fatal("JSON Marshal error:", err)
		return
	}
	buf := new(bytes.Buffer)
	err = json.Indent(buf, jsonBytes, "", "	")
	if err != nil {
		return
	}
	jsonstr = buf.String()
	return
}
