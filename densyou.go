package fmdu

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	TblDensyou      = "TBL_DENSYOU"
	ViewDensyou     = "VIEW_DENSYOU"
	TblDensyouCust  = "TBL_DENSYOU_CUSTOMER"
	ViewCustDensyou = "VIEW_CUSTOMER_DENSYOU"
	ViewDensyouCust = "VIEW_DENSYOU_CUSTOMER"
)

type ParameterDensyou struct {
	*Parameter
	ID             uint64
	DensyouCodeStr string
	FtthArea       string // ADD 2021/11/06
	Name           string
	Person         string
	Remarks        string
	ClientNo       uint64
	CustCnt        int8
	Eq             uint8
}

func (dbpara *ParameterAPI) InitParameterDensyou(c echo.Context) (para *ParameterDensyou, err error) {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	denstr := c.QueryParam("dencodestr")
	name, _ := decodeURLString(c.QueryParam("name"))
	fttharea := c.QueryParam("fttharea")
	remarks1 := c.QueryParam("remarks")
	person, _ := decodeURLString(c.QueryParam("person"))
	custcnt := -1
	tmp := c.QueryParam("custcnt")
	if 0 < len(tmp) {
		custcnt, _ = strconv.Atoi(tmp)
	}
	eq, _ := strconv.Atoi(c.QueryParam("eq"))

	para = new(ParameterDensyou)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.ID = uint64(id)
	para.DensyouCodeStr = denstr
	para.FtthArea = fttharea
	para.Name = name
	para.Person = person
	para.Remarks = remarks1
	para.CustCnt = int8(custcnt)
	para.Eq = uint8(eq)

	return
}

type Densyou struct {
	ID          uint64
	Name        string
	FtthArea    string // ADD 2021/11/06
	CustCnt     uint64
	FtthCnt     uint64
	HfcCnt      uint64
	FreeCnt     uint64
	FreeFtthCnt uint64
	FreeHfcCnt  uint64
	//ProcessStatus uint8
	//KoujiStatus   uint8
	//ProcessStatusStr string
	//KoujiStatusStr   string
	Person  sql.NullString
	Remarks sql.NullString
	IP      sql.NullString
	UserID  sql.NullString
	Ctim    time.Time
}
type Densyous []*Densyou

func (para ParameterDensyou) GetDensyou() (ds Densyous, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf(`select id,name,ftth_area,cust_cnt,ftth_cnt,hfc_cnt,free_cnt,free_ftth_cnt,free_hfc_cnt,person,remarks1,ip,user_id,ctim from %v `, ViewDensyou)

	addcnt := 0
	tmp := ""
	if 0 < para.ID {
		tmp = fmt.Sprintf(" (id = %v) ", para.ID)
		addcnt++
	}
	if 0 < len(para.DensyouCodeStr) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(densyou_code_str) like %v ) ", tmp, getNoCaseTextSQLStr(para.DensyouCodeStr))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(densyou_code_str) like %v ) ", getNoCaseTextSQLStr(para.DensyouCodeStr))
		}
		addcnt++
	}
	if 0 < len(para.FtthArea) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(ftth_area) like %v ) ", tmp, getNoCaseTextSQLStr(para.FtthArea))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(ftth_area) like %v ) ", getNoCaseTextSQLStr(para.FtthArea))
		}
		addcnt++
	}
	if 0 < len(para.Name) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(name) like %v ) ", tmp, getNoCaseTextSQLStr(para.Name))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(name) like %v ) ", getNoCaseTextSQLStr(para.Name))
		}
		addcnt++
	}
	if 0 < len(para.Person) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(person) like %v ) ", tmp, getNoCaseTextSQLStr(para.Person))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(person) like %v) ", getNoCaseTextSQLStr(para.Person))
		}
		addcnt++
	}
	if 0 < len(para.Remarks) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (remarks1 like '%%%v%%') ", tmp, para.Remarks)
		} else {
			tmp = fmt.Sprintf(" (remarks1 like '%%%v%%') ", para.Remarks)
		}
		addcnt++
	}
	if -1 < para.CustCnt {
		if 0 < addcnt {
			if para.Eq == 1 {
				tmp = fmt.Sprintf(" %v and (cust_cnt <= %v) ", tmp, para.CustCnt)
			} else {
				tmp = fmt.Sprintf(" %v and (cust_cnt >= %v) ", tmp, para.CustCnt)
			}
		} else {
			if para.Eq == 1 {
				tmp = fmt.Sprintf(" (cust_cnt <= %v) ", para.CustCnt)
			} else {
				tmp = fmt.Sprintf(" (cust_cnt >= %v) ", para.CustCnt)
			}
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
		return
	}
	defer rows.Close()

	ds = make(Densyous, 0)
	for rows.Next() {
		d := new(Densyou)
		if err := rows.Scan(&d.ID, &d.Name, &d.FtthArea, &d.CustCnt, &d.FtthCnt, &d.HfcCnt,
			&d.FreeCnt, &d.FreeFtthCnt, &d.FreeHfcCnt,
			&d.Person, &d.Remarks, &d.IP, &d.UserID, &d.Ctim); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		ds = append(ds, d)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (para *ParameterDensyou) UpdateDensyou() (ds Densyous, err error) {
	if para.ID == 0 {
		err = fmt.Errorf("*** Error *** IDが不正")
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT COUNT(*) AS cnt FROM %v WHERE id = %v;", ViewDensyou, para.ID)
	stmt, err := db.Prepare(sql)
	if err != nil {
		return
	}
	defer stmt.Close()

	cnt := 0
	crows, err := stmt.Query()
	for crows.Next() {
		crows.Scan(&cnt)
	}

	if cnt == 0 {
		err = fmt.Errorf("NoRecord[id=%v]", para.ID)
		return
	}

	sql = fmt.Sprintf("update %v ", TblDensyou)
	tmpsql := ""
	//addcnt := 0
	tmpsql = " set CTIM = now() "
	if 0 < len(para.Person) {
		if para.Person == NULL {
			tmpsql = fmt.Sprintf("%v, PERSON = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, PERSON = '%v' ", tmpsql, para.Person)
		}
	}
	if 0 < len(para.Remarks) {
		if para.Remarks == NULL {
			tmpsql = fmt.Sprintf("%v, REMARKS1 = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, REMARKS1 = '%v' ", tmpsql, para.Remarks)
		}
	}
	if 0 < len(para.IP) {
		tmpsql = fmt.Sprintf("%v, IP = '%v' ", tmpsql, para.IP)
	}
	if 0 < len(para.UserID) {
		tmpsql = fmt.Sprintf("%v, USER_ID = '%v' ", tmpsql, para.UserID)
	}
	if 0 < len(tmpsql) {
		sql = fmt.Sprintf("%v %v", sql, tmpsql)
	}
	footer := fmt.Sprintf(" where ID = %v;", para.ID)
	sql = fmt.Sprintf("%v %v", sql, footer)
	if para.Debug {
		fmt.Println(sql)
	}
	_, err = db.Query(sql)
	if err != nil {
		return
	}

	ds, err = para.GetDensyou()

	return
}

// SMSデータから電障コードと電障名を取得
func (para *Parameter) getDensyouListFromSMS() (ds Densyous, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf(`SELECT DISTINCT COL08,COL09 FROM %v order by COL08`, ViewCustDensyou)

	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	ds = make(Densyous, 0)
	for rows.Next() {
		d := new(Densyou)
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		ds = append(ds, d)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

type densyouCust struct {
	clientNo    uint64
	densyouCode uint
}
type densyouCusts []*densyouCust

// SMSデータから電障顧客を取得
func (para *Parameter) getDensyouCustFromSMS() (cs densyouCusts, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf(`SELECT COL23,COL08 FROM %v `, ViewCustDensyou)

	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	cs = make(densyouCusts, 0)
	for rows.Next() {
		c := new(densyouCust)
		if err := rows.Scan(&c.clientNo, &c.densyouCode); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}

		cs = append(cs, c)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
