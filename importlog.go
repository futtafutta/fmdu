package fmdu

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/inadenoita/futalib"
)

type ImportLog struct {
	idx            uint
	StartTime      string
	EndTime        string
	durationStr    string
	Duration       float64
	IP             string
	UserID         string
	status         uint8
	StatusStr      string
	UpdateCategory string
}

type ImportLogs []*ImportLog

const (
	layoutTimestamp = "2006-01-02 15:04:05"
	LayoutDay       = "2006-01-02"
)

func (para *ParameterImport) writeLog(succ int, start, end time.Time, dur time.Duration) (err error) {

	line := fmt.Sprintf("%v|%v|%v|%v|%v|%v", start.Format(layoutTimestamp), end.Format(layoutTimestamp), dur, para.IP, para.UserID, succ)
	fmt.Println(line)

	file, err := os.OpenFile(para.CrawlLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	fmt.Fprintln(file, line) //書き込み

	return
}

const (
	statusStrSuccess = "成功"
	statusStrFailure = "失敗"
)

func newImportLog(idx int, line string) (i *ImportLog, err error) {
	i = new(ImportLog)
	i.idx = uint(idx)
	ls := strings.Split(line, "|")
	if len(ls) < 6 {
		err = fmt.Errorf("*** Error *** invalid Columns Number[%v]", idx+1)
		return
	}
	i.StartTime = ls[0]
	i.EndTime = ls[1]
	i.durationStr = ls[2]
	tmps := i.durationStr
	tmps = strings.ReplaceAll(tmps, "s", "")
	tmps = strings.ReplaceAll(tmps, "S", "")
	i.Duration, _ = strconv.ParseFloat(tmps, 64)
	i.Duration = futalib.Round(i.Duration)
	i.IP = ls[3]
	i.UserID = ls[4]
	tmp, _ := strconv.Atoi(ls[5])
	i.status = uint8(tmp)
	if i.status == 0 {
		i.StatusStr = statusStrFailure
	} else {
		i.StatusStr = statusStrSuccess
	}
	if i.UserID == SystemUser {
		i.UpdateCategory = "日時バッチ"
	} else {
		i.UpdateCategory = "手動更新"
	}

	return
}

func (para *ParameterImport) GetImportLog() (is ImportLogs, err error) {
	_, err = os.Stat(para.CrawlLogPath)
	if err != nil {
		return
	}
	list, err := futalib.Csv2Slice(para.CrawlLogPath)
	if err != nil {
		return
	}
	is = make(ImportLogs, 0)
	for i, l := range *list {
		imp, err := newImportLog(i, l)
		if err != nil {
			continue
		}
		is = append(is, imp)
	}
	sort.SliceStable(is, func(i, j int) bool { return is[i].idx > is[j].idx })

	return
}
