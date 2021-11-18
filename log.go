package fmdu

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"bitbucket.org/inadenoita/futalib"
	"github.com/labstack/echo/v4"
)

type ParameterLog struct {
	*Parameter
	username string
}

func (dbpara *ParameterAPI) InitParameterLog(c echo.Context) (para *ParameterLog, err error) {
	para = new(ParameterLog)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.username = c.QueryParam("uname")

	return
}

type Log struct {
	idx     uint
	Time    string
	Status  uint8
	Method  string
	Service string
	UserID  string
	IP      string
	OS      string
	Browser string
}
type Logs []*Log

func (para *ParameterLog) GetLog() (logs Logs, err error) {

	_, err = os.Stat(para.LogPath)
	if err != nil {
		return
	}
	list, err := futalib.Csv2Slice(para.LogPath)
	if err != nil {
		return
	}
	logs = make(Logs, 0)
	for i, l := range *list {
		log, err := newLog(i, l)
		if err != nil {
			continue
		}
		if (0 < len(para.username)) && (para.username != log.UserID) {
			continue
		}
		logs = append(logs, log)
	}
	sort.SliceStable(logs, func(i, j int) bool { return logs[i].idx > logs[j].idx })
	return
}

func newLog(idx int, line string) (log *Log, err error) {
	ls := strings.Split(line, "|")
	if len(ls) < 8 {
		err = fmt.Errorf("*** Error *** invalid Columns Number[%v]", idx+1)
		return
	}
	if len(ls[0]) == 0 {
		err = fmt.Errorf("*** Error *** No Log Value Number[%v]", idx+1)
		return
	}
	log = new(Log)
	log.idx = uint(idx)
	log.Time = ls[0]
	tmp, _ := strconv.Atoi(ls[1])
	log.Status = uint8(tmp)
	log.Method = ls[2]
	log.Service = ls[3]
	log.UserID = ls[4]
	log.IP = ls[5]
	log.OS = ls[6]
	log.Browser = strings.ReplaceAll(ls[7], ",", " ")

	return
}
