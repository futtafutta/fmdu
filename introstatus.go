package fmdu

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	TblIntroStatus = "TBL_INTRO_STATUS"
)

type PerameterIntroStatus struct {
	*Parameter
	ID   int64
	Name string
}

func (dbpara *ParameterAPI) InitParameterIntroStatus(c echo.Context) (para *PerameterIntroStatus, err error) {
	id := -1
	tmp := c.QueryParam("id")
	if 0 < len(tmp) {
		id, _ = strconv.Atoi(tmp)
	}
	name := c.QueryParam("name")

	para = new(PerameterIntroStatus)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.ID = int64(id)
	para.Name = name
	return
}

type Status struct {
	ID   uint
	Name string
}
type Statuss []*Status

func (para *PerameterIntroStatus) GetIntroStatus() (ss Statuss, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,name FROM %v", TblIntroStatus)
	addcnt := 0
	tmp := ""
	if -1 < para.ID {
		tmp = fmt.Sprintf(" (id = %v) ", para.ID)
		addcnt++
	}
	if 0 < len(para.Name) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (name = '%v') ", tmp, para.Name)
		} else {
			tmp = fmt.Sprintf(" (name = '%v') ", para.Name)
		}
		addcnt++
	}
	if 0 < addcnt {
		sql = fmt.Sprintf("%v where %v order by id;", sql, tmp)
	} else {
		sql = fmt.Sprintf("%v order by id;", sql)
	}

	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	ss = make(Statuss, 0)
	for rows.Next() {
		s := new(Status)
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		ss = append(ss, s)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}
