package fmdu

import (
	"encoding/base64"
	"encoding/json"
	"strconv"

	"github.com/labstack/echo/v4"
)

/*
// ルート種類定数
const (
	RouteDown    = 0
	RouteUp      = 1
	RouteOpt     = 2
	RouteCurrent = 99
)
*/

const (
	defalutMaxLen = 10000
)

// Parameter ... iNA-Leaf WEB API 基本パラメータ構造体
type Parameter struct {
	*ParameterAPI
	KyokuCode string
	/*
		Longitude  float64
		Latitude   float64
		MinLng     float64
		MaxLng     float64
		MinLat     float64
		MaxLat     float64
	*/
	Debug      bool
	MaxRecSize int
	IP         string
	UserID     string
	Binfo      struct {
		Platform  string `json:"platform"`
		UserAgent string `json:"userAgent"`
	} `json:"binfo"`
}

// InitParameter ... パラメータ初期化
func (dbpara *ParameterAPI) InitParameter(c echo.Context) (*Parameter, error) {
	/*
		lat, _ := strconv.ParseFloat(c.QueryParam("lat"), 64)
		lng, _ := strconv.ParseFloat(c.QueryParam("lng"), 64)

		minx, _ := strconv.ParseFloat(c.QueryParam("minx"), 64)
		miny, _ := strconv.ParseFloat(c.QueryParam("miny"), 64)
		maxx, _ := strconv.ParseFloat(c.QueryParam("maxx"), 64)
		maxy, _ := strconv.ParseFloat(c.QueryParam("maxy"), 64)
	*/

	debug, _ := strconv.ParseBool(c.QueryParam("debug"))

	maxlen, _ := strconv.Atoi(c.QueryParam("maxlen"))
	if maxlen == 0 {
		maxlen = defalutMaxLen
	}
	ip := c.QueryParam("ip")
	userid := c.QueryParam("userid")
	binfostr := c.QueryParam("binfo")
	binfobyte, _ := decodeString(binfostr)

	para := &Parameter{}
	para.ParameterAPI = dbpara

	/*
		para.Latitude = lat
		para.Longitude = lng

		para.MinLng = minx
		para.MinLat = miny
		para.MaxLng = maxx
		para.MaxLat = maxy
	*/
	para.Debug = debug

	para.MaxRecSize = maxlen

	para.IP = ip
	para.UserID = userid

	json.Unmarshal(binfobyte, &para.Binfo)

	return para, nil
}

// Base64エンコード文字列パラメータのデコード(js<-->golang)
func decodeString(encstr string) (dec []byte, err error) {
	dec, err = base64.URLEncoding.DecodeString(encstr)
	return
}

// Base64エンコード文字列パラメータのデコード(js<-->golang)
func decodeURLString(encstr string) (decstr string, err error) {
	dec, err := base64.URLEncoding.DecodeString(encstr)
	if err != nil {
		return
	}
	decstr = string(dec)
	return
}
