package fmdu

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	TblIntro  = "TBL_INTRO"
	ViewIntro = "VIEW_INTRO"
)

type ParameterIntro struct {
	*Parameter
	ID                 uint64
	BuilNo             uint64
	BuilAddr           string
	BuilName           string
	BuilNoStr          string
	TourokuNoStr       string // ADD 2021/10/28
	FtthArea           string // ADD 2021/11/06
	RoomCnt            int8
	Eq                 uint8
	IsGen              int8 // ADD 2021/10/28
	IsFin              int8 // ADD 2021/10/28
	DirectStatus       int8
	DMFlag             int8
	DMDay              string
	ExplanationFlag    int8
	ExplanationPerson  string // ADD 2021/11/16
	ExplanationRemarks string
	ExplanationDay     string
	GenchouCompany     int8
	GenchouDay         string
	GenchouRtnDay      string // ADD 2021/10/26
	Yesno              int8
	YesnoReason        string
	KoujiName          string
	KoujiCompany       int8
	OrderDay           string
	KoujiDay           string
	DropDelFlag        int8   // ADD 2021/10/26
	DropDelDay         string // ADD 2021/10/26
	Remarks            string
}

func (dbpara *ParameterAPI) InitParameterIntrto(c echo.Context) (para *ParameterIntro, err error) {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	builno, _ := strconv.Atoi(c.QueryParam("builno"))
	builnostr := c.QueryParam("builnostr")
	builaddr, _ := decodeURLString(c.QueryParam("builaddr"))
	builname, _ := decodeURLString(c.QueryParam("builname"))
	tourokunostr := c.QueryParam("tourokunostr")
	fttharea := c.QueryParam("fttharea")
	roomcnt := -1
	tmp := c.QueryParam("roomcnt")
	if 0 < len(tmp) {
		roomcnt, _ = strconv.Atoi(tmp)
	}
	isgen := -1
	tmp = c.QueryParam("isgen")
	if 0 < len(tmp) {
		isgen, _ = strconv.Atoi(tmp)
	}
	isfin := -1
	tmp = c.QueryParam("isfin")
	if 0 < len(tmp) {
		isfin, _ = strconv.Atoi(tmp)
	}
	eq, _ := strconv.Atoi(c.QueryParam("eq"))
	dirstat := -1
	tmp = c.QueryParam("status")
	if 0 < len(tmp) {
		dirstat, _ = strconv.Atoi(tmp)
	}
	dmflag := -1
	tmp = c.QueryParam("dmflag")
	if 0 < len(tmp) {
		dmflag, _ = strconv.Atoi(tmp)
	}
	dmday := c.QueryParam("dmday")
	expflag := -1
	tmp = c.QueryParam("expflag")
	if 0 < len(tmp) {
		expflag, _ = strconv.Atoi(tmp)
	}
	expperson, _ := decodeURLString(c.QueryParam("expperson"))
	expremark := c.QueryParam("expremark")
	expday := c.QueryParam("expday")

	gencomp := -1
	tmp = c.QueryParam("gencomp")
	if 0 < len(tmp) {
		gencomp, _ = strconv.Atoi(tmp)
	}
	genday := c.QueryParam("genday")
	genrtnday := c.QueryParam("genrtnday")

	yesno := -1
	tmp = c.QueryParam("yesno")
	if 0 < len(tmp) {
		yesno, _ = strconv.Atoi(tmp)
	}
	yesnoreason := c.QueryParam("yesnoreason")

	koujiname, _ := decodeURLString(c.QueryParam("koujiname"))
	orderday := c.QueryParam("orderday")
	koujiday := c.QueryParam("koujiday")
	koujicomp := -1
	tmp = c.QueryParam("koujicomp")
	if 0 < len(tmp) {
		koujicomp, _ = strconv.Atoi(tmp)
	}
	dropdelflag := -1
	tmp = c.QueryParam("dropdelflag")
	if 0 < len(tmp) {
		dropdelflag, _ = strconv.Atoi(tmp)
	}
	dropdelday := c.QueryParam("dropdelday")

	remarks1 := c.QueryParam("remarks")

	para = new(ParameterIntro)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.ID = uint64(id)
	para.BuilNo = uint64(builno)
	para.BuilAddr = builaddr
	para.BuilName = builname

	para.BuilNoStr = builnostr
	para.TourokuNoStr = tourokunostr
	para.FtthArea = fttharea

	para.RoomCnt = int8(roomcnt)
	if 0 < eq {
		para.Eq = uint8(eq)
	}
	para.IsGen = int8(isgen)
	para.IsFin = int8(isfin)

	para.DirectStatus = int8(dirstat)
	para.DMFlag = int8(dmflag)
	para.DMDay = dmday
	para.ExplanationFlag = int8(expflag)
	para.ExplanationPerson = expperson
	para.ExplanationRemarks = expremark
	para.ExplanationDay = expday

	para.GenchouCompany = int8(gencomp)
	para.GenchouDay = genday
	para.GenchouRtnDay = genrtnday

	para.Yesno = int8(yesno)
	para.YesnoReason = yesnoreason

	para.KoujiName = koujiname
	para.KoujiCompany = int8(koujicomp)
	para.OrderDay = orderday
	para.KoujiDay = koujiday

	para.DropDelFlag = int8(dropdelflag)
	para.DropDelDay = dropdelday

	para.Remarks = remarks1

	return
}

type Intro struct {
	ID                 uint64
	IsGenchou          string // ADD 2021/10/28
	IsFinish           string // ADD 2021/10/28
	isGenchou          uint8  // ADD 2021/10/28
	isFinish           uint8  // ADD 2021/10/28
	TourokuNo          uint64 // ADD 2021/10/26
	BuilNo             uint64
	BuilAddr           string
	BuilName           string
	FtthArea           string // ADD 2021/11/06
	BuilCnt            sql.NullInt32
	BuilStatus         sql.NullString
	BuilTV             sql.NullString
	BuilNet            sql.NullString
	BuilPhone          sql.NullString
	CustCnt            uint64
	FtthCnt            uint64
	HfcCnt             uint64
	HfcTvCnt           uint64
	StatusID           uint8
	Status             string
	DmFlag             uint8
	DmFlagStr          string
	DmDay              sql.NullTime
	ExplanationFlag    uint8
	ExplanationFlagStr string
	ExplanationPerson  sql.NullString
	ExplanationRemarks sql.NullString
	ExplanationDay     sql.NullTime
	GenchouCompanyID   sql.NullInt32
	GenchouCompany     sql.NullString
	GenchouDay         sql.NullTime
	GenchouRtnDay      sql.NullTime // ADD 2021/10/26
	Yesno              sql.NullInt32
	YesnoStr           string
	YesnoReason        sql.NullString
	KoujiName          sql.NullString
	KoujiCompanyID     sql.NullInt32
	KoujiCompany       sql.NullString
	OrderDay           sql.NullTime
	KoujiDay           sql.NullTime
	DropDelFlag        uint8        // ADD 2021/10/26
	DropDelFlagStr     string       // ADD 2021/10/26
	DropDelDay         sql.NullTime // ADD 2021/10/26
	Longitude          sql.NullFloat64
	Latitude           sql.NullFloat64
	Remarks1           sql.NullString
	IP                 sql.NullString
	UserID             sql.NullString
	Ctim               time.Time
}

type Intros []*Intro

func (para *ParameterIntro) GetIntro() (is Intros, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf(`select id,is_genchou,is_finish,ftth_area,touroku_no,buil_no,buil_addr,buil_name,buil_cnt,
	buil_status,buil_tv,buil_net,buil_phone,cust_cnt,ftth_cnt,hfc_cnt,hfc_tv_cnt,direct_status,status,
	dm_flag,dm_day,explanation_flag,explanation_person,explanation_remarks,explanation_day,
	genchou_company_id,genchou_company,genchou_day,genchou_return_day,yesno,yesno_reason,
	kouji_name,kouji_company_id,kouji_company,order_day,kouji_day,drop_del_flag,drop_del_day,
	remarks1,user_id,ctim from %v `, ViewIntro)

	addcnt := 0
	tmp := ""
	if 0 < para.ID {
		tmp = fmt.Sprintf(" (id = %v) ", para.ID)
		addcnt++
	}
	if 0 < para.BuilNo {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (buil_no = %v) ", tmp, para.BuilNo)
		} else {
			tmp = fmt.Sprintf(" (buil_no = %v) ", para.BuilNo)
		}
		addcnt++
	}
	if 0 < len(para.BuilNoStr) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(buil_no_str) like %v ) ", tmp, getNoCaseTextSQLStr(para.BuilNoStr))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(buil_no_str) like %v ) ", getNoCaseTextSQLStr(para.BuilNoStr))
		}
		addcnt++
	}
	if 0 < len(para.BuilAddr) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(buil_addr) like %v ) ", tmp, getNoCaseTextSQLStr(para.BuilAddr))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(buil_addr) like %v ) ", getNoCaseTextSQLStr(para.BuilAddr))
		}
		addcnt++
	}
	if 0 < len(para.BuilName) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(buil_name) like %v ) ", tmp, getNoCaseTextSQLStr(para.BuilName))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(buil_name) like %v ) ", getNoCaseTextSQLStr(para.BuilName))
		}
		addcnt++
	}
	if 0 < len(para.TourokuNoStr) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(touroku_no_str) like %v ) ", tmp, getNoCaseTextSQLStr(para.TourokuNoStr))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(touroku_no_str) like %v ) ", getNoCaseTextSQLStr(para.TourokuNoStr))
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

	if -1 < para.RoomCnt {
		if 0 < addcnt {
			if para.Eq == 1 {
				tmp = fmt.Sprintf(" %v and (buil_cnt <= %v) ", tmp, para.RoomCnt)
			} else {
				tmp = fmt.Sprintf(" %v and (buil_cnt >= %v) ", tmp, para.RoomCnt)
			}
		} else {
			if para.Eq == 1 {
				tmp = fmt.Sprintf(" (buil_cnt <= %v) ", para.RoomCnt)
			} else {
				tmp = fmt.Sprintf(" (buil_cnt >= %v) ", para.RoomCnt)
			}
		}
		addcnt++
	}

	if -1 < para.IsGen {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (is_genchou = %v) ", tmp, para.IsGen)
		} else {
			tmp = fmt.Sprintf("  (is_genchou = %v) ", para.IsGen)
		}
		addcnt++
	}

	if -1 < para.IsFin {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (is_finish = %v) ", tmp, para.IsFin)
		} else {
			tmp = fmt.Sprintf(" (is_finish = %v) ", para.IsFin)
		}
		addcnt++
	}
	// ※Add 2021/11/12
	if -1 < para.DirectStatus {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (direct_status = %v) ", tmp, para.DirectStatus)
		} else {
			tmp = fmt.Sprintf(" (direct_status = %v) ", para.DirectStatus)
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

	llmap, err := para.getLngLatMap()
	if err != nil {
		fmt.Println(err)
	}

	is = make(Intros, 0)
	for rows.Next() {
		i := new(Intro)
		if err := rows.Scan(&i.ID, &i.isGenchou, &i.isFinish, &i.FtthArea, &i.TourokuNo,
			&i.BuilNo, &i.BuilAddr, &i.BuilName, &i.BuilCnt, &i.BuilStatus, &i.BuilTV, &i.BuilNet, &i.BuilPhone,
			&i.CustCnt, &i.FtthCnt, &i.HfcCnt, &i.HfcTvCnt,
			&i.StatusID, &i.Status, &i.DmFlag, &i.DmDay, &i.ExplanationFlag, &i.ExplanationPerson, &i.ExplanationRemarks, &i.ExplanationDay,
			&i.GenchouCompanyID, &i.GenchouCompany, &i.GenchouDay, &i.GenchouRtnDay, &i.Yesno, &i.YesnoReason,
			&i.KoujiName, &i.KoujiCompanyID, &i.KoujiCompany, &i.OrderDay, &i.KoujiDay, &i.DropDelFlag, &i.DropDelDay,
			&i.Remarks1, &i.UserID, &i.Ctim); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		i.setFlag2Str()
		i.setYesnoStr()
		i.setIsIntroCompleted()
		i.Longitude, i.Latitude = llmap.getLngLat(i.BuilNo)

		is = append(is, i)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

// 光ダイレクト・調査結果 セット
func (i *Intro) setIsIntroCompleted() {
	if 0 < i.isGenchou {
		i.IsGenchou = "○"
	}
	if 0 < i.isFinish {
		i.IsFinish = "○"
	}
}

func (i *Intro) setFlag2Str() {
	ynstr := ""
	if 0 < i.DmFlag {
		ynstr = "○"
	}
	i.DmFlagStr = ynstr

	ynstr = ""
	if i.ExplanationFlag == 1 {
		ynstr = "○"
	} else if i.ExplanationFlag == 2 {
		ynstr = "(予)"
	}
	i.ExplanationFlagStr = ynstr

	ynstr = ""
	if 0 < i.DropDelFlag {
		ynstr = "○"
	}
	i.DropDelFlagStr = ynstr
}

func (i *Intro) setYesnoStr() {
	ynstr := ""
	switch i.Yesno.Int32 {
	case 1:
		ynstr = "○"
	case 2:
		ynstr = "×"
	case 3:
		ynstr = "△"
	}
	i.YesnoStr = ynstr
}

func (para *ParameterIntro) UpdateIntro() (is Intros, err error) {
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

	sql := fmt.Sprintf("SELECT COUNT(*) AS cnt FROM %v WHERE id = %v;", ViewIntro, para.ID)
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

	sql = fmt.Sprintf("update %v ", TblIntro)
	tmpsql := ""
	//addcnt := 0
	tmpsql = " set CTIM = now() "
	if 0 < len(para.IP) {
		tmpsql = fmt.Sprintf("%v, IP = '%v' ", tmpsql, para.IP)
	}
	if 0 < len(para.UserID) {
		tmpsql = fmt.Sprintf("%v, USER_ID = '%v' ", tmpsql, para.UserID)
	}
	if -1 < para.DirectStatus {
		tmpsql = fmt.Sprintf("%v, DIRECT_STATUS = %v ", tmpsql, para.DirectStatus)
	}

	if -1 < para.DMFlag {
		tmpsql = fmt.Sprintf("%v, DM_FLAG = %v ", tmpsql, para.DMFlag)
	}
	if 0 < len(para.DMDay) {
		if para.DMDay == NULL {
			tmpsql = fmt.Sprintf("%v, DM_DAY = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, DM_DAY = '%v' ", tmpsql, para.DMDay)
		}
	}

	if -1 < para.ExplanationFlag {
		tmpsql = fmt.Sprintf("%v, EXPLANATION_FLAG = %v ", tmpsql, para.ExplanationFlag)
	}

	if 0 < len(para.ExplanationPerson) {
		if para.ExplanationPerson == NULL {
			tmpsql = fmt.Sprintf("%v, EXPLANATION_PERSON = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, EXPLANATION_PERSON = '%v' ", tmpsql, para.ExplanationPerson)
		}
	}

	if 0 < len(para.ExplanationRemarks) {
		if para.ExplanationRemarks == NULL {
			tmpsql = fmt.Sprintf("%v, EXPLANATION_REMARKS = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, EXPLANATION_REMARKS = '%v' ", tmpsql, para.ExplanationRemarks)
		}
	}

	if 0 < len(para.ExplanationDay) {
		if para.ExplanationDay == NULL {
			tmpsql = fmt.Sprintf("%v, EXPLANATION_DAY = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, EXPLANATION_DAY = '%v' ", tmpsql, para.ExplanationDay)
		}
	}
	if -1 < para.GenchouCompany {
		tmpsql = fmt.Sprintf("%v, GENCHOU_COMPANY_ID = %v ", tmpsql, para.GenchouCompany)
	}
	if 0 < len(para.GenchouDay) {
		if para.GenchouDay == NULL {
			tmpsql = fmt.Sprintf("%v, GENCHOU_DAY = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, GENCHOU_DAY = '%v' ", tmpsql, para.GenchouDay)
		}
	}
	if 0 < len(para.GenchouRtnDay) {
		if para.GenchouRtnDay == NULL {
			tmpsql = fmt.Sprintf("%v, GENCHOU_RETURN_DAY = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, GENCHOU_RETURN_DAY = '%v' ", tmpsql, para.GenchouRtnDay)
		}
	}
	if -1 < para.Yesno {
		tmpsql = fmt.Sprintf("%v, YESNO = %v ", tmpsql, para.Yesno)
	}
	if 0 < len(para.YesnoReason) {
		if para.YesnoReason == NULL {
			tmpsql = fmt.Sprintf("%v, YESNO_REASON = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, YESNO_REASON = '%v' ", tmpsql, para.YesnoReason)
		}
	}
	if 0 < len(para.KoujiName) {
		if para.KoujiName == NULL {
			tmpsql = fmt.Sprintf("%v, KOUJI_NAME = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, KOUJI_NAME = '%v' ", tmpsql, para.KoujiName)
		}
	}
	if -1 < para.KoujiCompany {
		tmpsql = fmt.Sprintf("%v, KOUJI_COMPANY_ID = %v ", tmpsql, para.KoujiCompany)
	}
	if 0 < len(para.OrderDay) {
		if para.OrderDay == NULL {
			tmpsql = fmt.Sprintf("%v, ORDER_DAY = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, ORDER_DAY = '%v' ", tmpsql, para.OrderDay)
		}
	}
	if 0 < len(para.KoujiDay) {
		if para.KoujiDay == NULL {
			tmpsql = fmt.Sprintf("%v, KOUJI_DAY = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, KOUJI_DAY = '%v' ", tmpsql, para.KoujiDay)
		}
	}

	if -1 < para.DropDelFlag {
		tmpsql = fmt.Sprintf("%v, DROP_DEL_FLAG = %v ", tmpsql, para.DropDelFlag)
	}

	if 0 < len(para.DropDelDay) {
		if para.DropDelDay == NULL {
			tmpsql = fmt.Sprintf("%v, DROP_DEL_DAY = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, DROP_DEL_DAY = '%v' ", tmpsql, para.DropDelDay)
		}
	}

	if 0 < len(para.Remarks) {
		if para.Remarks == NULL {
			tmpsql = fmt.Sprintf("%v, REMARKS1 = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, REMARKS1 = '%v' ", tmpsql, para.Remarks)
		}

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

	is, err = para.GetIntro()

	return
}
