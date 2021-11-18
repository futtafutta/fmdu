package fmdu

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	TblCpl  = "tbl_cpl"
	ViewCpl = "view_cpl"
)

type ParameterCpl struct {
	*Parameter
	ID              uint64
	FID             uint64
	FtypeID         uint8
	DropNo          uint8
	CplName         string
	BrCnt           uint8
	CplNo           string
	InstallationDay string
	CancelDay       string
	Remarks         string
}

func (dbpara *ParameterAPI) InitParameterCpl(c echo.Context) (para *ParameterCpl, err error) {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	fid, _ := strconv.Atoi(c.QueryParam("fid"))
	ftypeid, _ := strconv.Atoi(c.QueryParam("ftypeid"))
	dropno, _ := strconv.Atoi(c.QueryParam("dropno"))
	cplname := c.QueryParam("cplname")
	brcnt, _ := strconv.Atoi(c.QueryParam("brcnt"))
	installday := c.QueryParam("installday")
	cancelday := c.QueryParam("cancelday")
	remarks := c.QueryParam("remarks")

	cplno := c.QueryParam("cplno")

	para = new(ParameterCpl)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.ID = uint64(id)
	para.FID = uint64(fid)
	para.FtypeID = uint8(ftypeid)
	para.DropNo = uint8(dropno)
	para.CplName = cplname
	para.BrCnt = uint8(brcnt)
	para.InstallationDay = installday
	para.CancelDay = cancelday
	para.Remarks = remarks
	para.CplNo = cplno

	return
}

type Cpl struct {
	ID              uint64
	FID             uint64
	FtypeID         uint8
	CplNo           string
	DropNo          uint8
	CplName         string
	BrCnt           uint8
	UsePort         uint8
	EmptyPort       uint8
	installday      *time.Time
	cancelday       *time.Time
	InstallationDay string
	CancelDay       string
	Remarks         sql.NullString
	IP              sql.NullString
	UserID          sql.NullString
	Ctim            time.Time
	TourokuNo       uint64
	BuilNo          uint64
	BuilName        string
}
type Cpls []*Cpl

// GetCpl ... カプラデータ取得
func (para *ParameterCpl) GetCpl() (cs Cpls, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,fid,touroku_no,buil_no,buil_name,ftype_id,cpl_no,drop_no,cpl_name,br_cnt,installation_day,cancel_day,remarks,ip,user_id,ctim FROM %v", ViewCpl)
	addcnt := 0
	tmp := ""
	if 0 < para.ID {
		tmp = fmt.Sprintf(" (id = %v) ", para.ID)
		addcnt++
	}
	if 0 < para.FID {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (fid = %v) ", tmp, para.FID)
		} else {
			tmp = fmt.Sprintf(" (fid = %v) ", para.FID)
		}
		addcnt++
	}
	if 0 < para.FtypeID {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (ftype_id = %v) ", tmp, para.FtypeID)
		} else {
			tmp = fmt.Sprintf(" (ftype_id = %v) ", para.FtypeID)
		}
		addcnt++
	}
	if 0 < len(para.CplNo) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(cpl_no) like %v) ", tmp, getNoCaseTextSQLStr(para.CplNo))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(cpl_no) like %v) ", getNoCaseTextSQLStr(para.CplNo))
		}
		addcnt++
	}
	if 0 < len(para.Remarks) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (remarks like '%%%v%%') ", tmp, para.Remarks)
		} else {
			tmp = fmt.Sprintf(" (remarks like '%%%v%%') ", para.Remarks)
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

	cs = make(Cpls, 0)
	for rows.Next() {
		c := new(Cpl)
		if err := rows.Scan(&c.ID, &c.FID, &c.TourokuNo, &c.BuilNo, &c.BuilName, &c.FtypeID, &c.CplNo, &c.DropNo, &c.CplName, &c.BrCnt,
			&c.installday, &c.cancelday, &c.Remarks, &c.IP, &c.UserID, &c.Ctim); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		if c.installday != nil {
			c.InstallationDay = c.installday.Format("2006-01-02")
		}
		if c.cancelday != nil {
			c.CancelDay = c.cancelday.Format("2006-01-02")
		}

		// 顧客情報からカプラ使用端子数を取得
		sql = fmt.Sprintf("SELECT COUNT(*) AS cnt FROM %v WHERE (COL01 = %v) AND  (GET_NOCASE_TEXT(COL19) like %v);", ViewCustomer, c.BuilNo, getNoCaseTextSQLStr(c.CplNo))
		stmt, err := db.Prepare(sql)
		if err != nil {
			return cs, err
		}
		defer stmt.Close()

		cnt := 0
		crows, _ := stmt.Query()
		for crows.Next() {
			crows.Scan(&cnt)
		}
		c.UsePort = uint8(cnt)
		c.EmptyPort = c.BrCnt - c.UsePort

		cs = append(cs, c)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (para *ParameterCpl) InsertCpl() (cs Cpls, err error) {
	if para.FID == 0 {
		err = fmt.Errorf("*** Error *** FIDが不正")
		return
	}
	if para.FtypeID == 0 {
		err = fmt.Errorf("*** Error *** FtypeIDが不正")
		return
	}
	if para.DropNo == 0 {
		err = fmt.Errorf("*** Error *** DropNoが不正")
		return
	}
	if len(para.CplName) == 0 {
		err = fmt.Errorf("*** Error *** CplNameが不正")
		return
	}
	if para.BrCnt == 0 {
		err = fmt.Errorf("*** Error *** BrCnt不正")
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	cmd := fmt.Sprintf("SELECT id,buil_no FROM %v where (fid = %v) and (drop_no = %v) and (cpl_name = '%v')", ViewCpl, para.FID, para.DropNo, para.CplName)
	rows, err := db.Query(cmd)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	idx := 0
	var c *Cpl
	for rows.Next() {
		c = new(Cpl)
		if err := rows.Scan(&c.ID, &c.BuilNo); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		idx++
		break
	}
	if err = rows.Err(); err != nil {
		return
	}
	if 0 < idx {
		err = fmt.Errorf("*** Error *** カプラ番号重複 建物番号：%v カプラ番号 : %v%v", c.BuilNo, para.DropNo, para.CplName)
		return
	}
	var fid uint
	if (0 < len(para.InstallationDay)) && (0 < len(para.CancelDay)) {
		err = db.QueryRow("INSERT INTO TBL_CPL (fid,ftype_id,drop_no,cpl_name,br_cnt,installation_day,cancel_day,remarks,ip,user_id,ctim ) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,now()) RETURNING id",
			para.FID, para.FtypeID, para.DropNo, para.CplName, para.BrCnt, para.InstallationDay, para.CancelDay, para.Remarks, para.IP, para.UserID).Scan(&fid)
	} else if 0 < len(para.InstallationDay) {
		err = db.QueryRow("INSERT INTO TBL_CPL (fid,ftype_id,drop_no,cpl_name,br_cnt,installation_day,remarks,ip,user_id,ctim ) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,now()) RETURNING id",
			para.FID, para.FtypeID, para.DropNo, para.CplName, para.BrCnt, para.InstallationDay, para.Remarks, para.IP, para.UserID).Scan(&fid)
	} else if 0 < len(para.CancelDay) {
		err = db.QueryRow("INSERT INTO TBL_CPL (fid,ftype_id,drop_no,cpl_name,br_cnt,cancel_day,remarks,ip,user_id,ctim ) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,now()) RETURNING id",
			para.FID, para.FtypeID, para.DropNo, para.CplName, para.BrCnt, para.CancelDay, para.Remarks, para.IP, para.UserID).Scan(&fid)
	} else {
		err = db.QueryRow("INSERT INTO TBL_CPL (fid,ftype_id,drop_no,cpl_name,br_cnt,remarks,ip,user_id,ctim ) VALUES($1,$2,$3,$4,$5,$6,$7,$8,now()) RETURNING id",
			para.FID, para.FtypeID, para.DropNo, para.CplName, para.BrCnt, para.Remarks, para.IP, para.UserID).Scan(&fid)
	}

	if err != nil {
		return
	}
	fmt.Println(fid)
	para.ID = uint64(fid)

	cs, err = para.GetCpl()

	return
}

func (para *ParameterCpl) UpdateCpl() (cs Cpls, err error) {
	if para.ID == 0 {
		err = fmt.Errorf("*** Error *** IDが不正")
		return
	}
	if para.FID == 0 {
		err = fmt.Errorf("*** Error *** FIDが不正")
		return
	}
	if len(para.CplNo) == 0 {
		err = fmt.Errorf("*** Error *** CplNoが不正")
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	sql := fmt.Sprintf("SELECT id,buil_no FROM %v where (fid = %v) and (cpl_no = '%v') and (id != %v)", ViewCpl, para.FID, para.CplNo, para.ID)
	rows, err := db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	idx := 0
	var c *Cpl
	for rows.Next() {
		c = new(Cpl)
		if err := rows.Scan(&c.ID, &c.BuilNo); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		idx++
		break
	}
	if err = rows.Err(); err != nil {
		return
	}
	if 0 < idx {
		err = fmt.Errorf("*** Error *** カプラ番号重複 建物番号：%v カプラ番号 : %v%v", c.BuilNo, para.DropNo, para.CplName)
		return
	}

	sql = fmt.Sprintf("update %v ", TblCpl)
	tmpsql := ""
	//addcnt := 0
	tmpsql = " set CTIM = now() "
	if 0 < para.FID {
		tmpsql = fmt.Sprintf("%v, FID = %v ", tmpsql, para.FID)
	}
	if 0 < para.FtypeID {
		tmpsql = fmt.Sprintf("%v, FTYPE_ID = %v ", tmpsql, para.FtypeID)
	}
	if 0 < para.DropNo {
		tmpsql = fmt.Sprintf("%v, DROP_NO = %v ", tmpsql, para.DropNo)
	}
	if 0 < len(para.CplName) {
		tmpsql = fmt.Sprintf("%v, CPL_NAME = '%v' ", tmpsql, para.CplName)
	}
	if 0 < para.BrCnt {
		tmpsql = fmt.Sprintf("%v, BR_CNT = %v ", tmpsql, para.BrCnt)
	}
	if 0 < len(para.InstallationDay) {
		if para.InstallationDay == NULL {
			tmpsql = fmt.Sprintf("%v, INSTALLATION_DAY = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, INSTALLATION_DAY = '%v' ", tmpsql, para.InstallationDay)
		}
	}
	if 0 < len(para.CancelDay) {
		if para.CancelDay == NULL {
			tmpsql = fmt.Sprintf("%v, CANCEL_DAY = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, CANCEL_DAY = '%v' ", tmpsql, para.CancelDay)
		}
	}
	if 0 < len(para.Remarks) {
		if para.Remarks == NULL {
			tmpsql = fmt.Sprintf("%v, REMARKS = %v ", tmpsql, NULL)
		} else {
			tmpsql = fmt.Sprintf("%v, REMARKS = '%v' ", tmpsql, para.Remarks)
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

	cs, err = para.GetCpl()

	return
}

// DeleteCpl ... カプラの削除
func (para *ParameterCpl) DeleteCpl() (cs Cpls, err error) {
	if para.ID == 0 {
		err = fmt.Errorf("*** Error *** IDが不正")
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()

	// 図形座標削除
	_, err = db.Query("DELETE FROM TBL_CPL WHERE ID = $1", para.ID)
	if err != nil {
		return
	}

	return
}
