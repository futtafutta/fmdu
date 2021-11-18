package fmdu

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	TblCustomer       = "TBL_CUSTOMER"
	ViewCustomer      = "VIEW_FMDU_CUST"
	ViewIntroCustomer = "VIEW_INTRO_CUST"

	TblServiceSegment = "TBL_SERVICE_SEGMENT"
)

const (
	AppIdxFmdu    uint8 = 1
	AppIdxIntro   uint8 = 2
	AppIdxDensyou uint8 = 3
)

type ParameterCustomer struct {
	*Parameter
	FID uint64
	//BuilNo uint64
	ClientNo      uint64
	CplNo         string
	AppKind       uint8
	ProcessStatus int8
	KoujiStatus   int8
	Person        string
	Remarks       string
}

func (dbpara *ParameterAPI) InitParameterCustomer(c echo.Context) (para *ParameterCustomer, err error) {
	fid, _ := strconv.Atoi(c.QueryParam("fid"))
	//builno, _ := strconv.Atoi(c.QueryParam("builno"))
	cplno := c.QueryParam("cplno")
	clino, _ := strconv.Atoi(c.QueryParam("clino"))
	appidx, _ := strconv.Atoi(c.QueryParam("app"))
	procstat := -1
	tmp := c.QueryParam("procstat")
	if 0 < len(tmp) {
		procstat, _ = strconv.Atoi(tmp)
	}
	koujistat := -1
	tmp = c.QueryParam("koujistat")
	if 0 < len(tmp) {
		koujistat, _ = strconv.Atoi(tmp)
	}
	person := c.QueryParam("person")
	remarks := c.QueryParam("remarks")

	para = new(ParameterCustomer)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.FID = uint64(fid)
	//para.BuilNo = uint64(builno)
	para.ClientNo = uint64(clino)
	para.CplNo = cplno
	para.AppKind = uint8(appidx)

	para.ProcessStatus = int8(procstat)
	para.KoujiStatus = int8(koujistat)
	para.Person = person
	para.Remarks = remarks

	return
}

type Customer struct {
	FID              uint64
	COL01            uint64
	COL02            uint16
	COL03            string
	COL04            uint16
	COL05            string
	COL06            sql.NullString
	COL07            sql.NullString
	COL08            sql.NullString
	COL09            sql.NullString
	COL10            sql.NullString
	COL11            sql.NullString
	COL12            sql.NullString
	COL13            sql.NullString
	COL14            sql.NullString
	COL15            sql.NullString
	COL16            sql.NullString
	COL17            sql.NullString
	COL18            sql.NullString
	COL19            sql.NullString
	COL20            sql.NullString
	COL21            sql.NullString
	COL22            sql.NullInt32
	COL23            sql.NullInt64
	COL24            sql.NullString
	COL25            sql.NullString
	COL26            sql.NullString
	COL27            sql.NullString
	COL28            sql.NullString
	COL29            sql.NullString
	COL30            sql.NullString
	COL31            sql.NullString
	COL32            sql.NullString
	COL33            sql.NullString
	COL34            sql.NullString
	COL35            sql.NullString
	COL36            sql.NullString
	COL37            string
	COL38            string
	COL39            string
	dbCOL37          sql.NullTime
	dbCOL38          sql.NullTime
	dbCOL39          sql.NullTime
	IsFTTH           uint8
	IsPaid           uint8
	IsDensyouOnly    uint8
	DensyouOnlyStr   string
	Longitude        sql.NullFloat64
	Latitude         sql.NullFloat64
	ProcessStatus    uint8
	KoujiStatus      uint8
	ProcessStatusStr string
	KoujiStatusStr   string
	Person           sql.NullString
	Remarks          sql.NullString
}
type Customers []*Customer

func (para *ParameterCustomer) GetCustomer() (cs Customers, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	tbl := ViewCustomer
	if para.AppKind == AppIdxIntro {
		tbl = ViewIntroCustomer
	} else if para.AppKind == AppIdxDensyou {
		tbl = ViewDensyouCust
	}

	sql := ""
	if para.AppKind == AppIdxDensyou {
		sql = fmt.Sprintf(`SELECT FID,COL01,COL02,COL03,COL04,COL05,COL06,COL07,COL08,COL09,COL10,COL11,COL12,COL13,COL14,COL15,COL16,COL17,COL18,COL19,COL20,
		COL21,COL22,COL23,COL24,COL25,COL26,COL27,COL28,COL29,COL30,COL31,COL32,COL33,COL34,COL35,COL36,COL37,COL38,COL39,
		FTTH,PAID,DENSYOU_ONLY,LONGITUDE,LATITUDE,PROCESS_STATUS,KOUJI_STATUS,PERSON,REMARKS1 FROM %v`, tbl)
	} else {
		sql = fmt.Sprintf(`SELECT FID,COL01,COL02,COL03,COL04,COL05,COL06,COL07,COL08,COL09,COL10,COL11,COL12,COL13,COL14,COL15,COL16,COL17,COL18,COL19,COL20,
		COL21,COL22,COL23,COL24,COL25,COL26,COL27,COL28,COL29,COL30,COL31,COL32,COL33,COL34,COL35,COL36,COL37,COL38,COL39,FTTH FROM %v`, tbl)
	}

	addcnt := 0
	tmp := ""
	if 0 < para.FID {
		if para.AppKind == AppIdxDensyou {
			tmp = fmt.Sprintf(" (FID = %v) ", para.FID)
			addcnt++
		} else {
			bnos := para.getFid2BuilNo(para.FID)
			if 0 < len(bnos) {
				tmp = fmt.Sprintf(" (COL01 in (%v) ) ", strings.Join(bnos, ","))
				addcnt++
			}
		}
	}
	if 0 < para.ClientNo {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (COL23 = %v) ", tmp, para.ClientNo)
		} else {
			tmp = fmt.Sprintf(" (COL23 = %v) ", para.ClientNo)
		}
		addcnt++
	}
	if 0 < len(para.CplNo) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(COL19) = %v ) ", tmp, getNoCaseTextSQLStr(para.CplNo))
		} else {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(COL19) = %v ) ", getNoCaseTextSQLStr(para.CplNo))
		}
		addcnt++
	}
	footer := " order by COL01, COL02;"
	if 0 < addcnt {
		sql = fmt.Sprintf("%v where %v %v", sql, tmp, footer)
	} else {
		sql = fmt.Sprintf("%v %v", sql, footer)
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

	cs = make(Customers, 0)
	for rows.Next() {
		c := new(Customer)
		if para.AppKind == AppIdxDensyou {
			if err := rows.Scan(&c.FID, &c.COL01, &c.COL02, &c.COL03, &c.COL04, &c.COL05, &c.COL06, &c.COL07, &c.COL08, &c.COL09, &c.COL10,
				&c.COL11, &c.COL12, &c.COL13, &c.COL14, &c.COL15, &c.COL16, &c.COL17, &c.COL18, &c.COL19, &c.COL20,
				&c.COL21, &c.COL22, &c.COL23, &c.COL24, &c.COL25, &c.COL26, &c.COL27, &c.COL28, &c.COL29, &c.COL30,
				&c.COL31, &c.COL32, &c.COL33, &c.COL34, &c.COL35, &c.COL36, &c.dbCOL37, &c.dbCOL38, &c.dbCOL39,
				&c.IsFTTH, &c.IsPaid, &c.IsDensyouOnly,
				&c.Longitude, &c.Latitude,
				&c.ProcessStatus, &c.KoujiStatus, &c.Person, &c.Remarks); err != nil {
				fmt.Printf("値の取得に失敗しました。: %v\n", err)
				continue
			}
			if c.IsFTTH == 0 {
				if c.ProcessStatus == 0 {
					c.ProcessStatusStr = ""
				} else {
					c.ProcessStatusStr = "○"
				}
				if c.KoujiStatus == 0 {
					c.KoujiStatusStr = ""
				} else {
					c.KoujiStatusStr = "○"
				}
			} else {
				c.ProcessStatusStr = "○"
				c.KoujiStatusStr = "○"
			}
			if c.IsDensyouOnly == 1 {
				c.DensyouOnlyStr = "電障のみ"
			} else {
				c.DensyouOnlyStr = "電障+有料サ"
			}
		} else {
			if err := rows.Scan(&c.FID, &c.COL01, &c.COL02, &c.COL03, &c.COL04, &c.COL05, &c.COL06, &c.COL07, &c.COL08, &c.COL09, &c.COL10,
				&c.COL11, &c.COL12, &c.COL13, &c.COL14, &c.COL15, &c.COL16, &c.COL17, &c.COL18, &c.COL19, &c.COL20,
				&c.COL21, &c.COL22, &c.COL23, &c.COL24, &c.COL25, &c.COL26, &c.COL27, &c.COL28, &c.COL29, &c.COL30,
				&c.COL31, &c.COL32, &c.COL33, &c.COL34, &c.COL35, &c.COL36, &c.dbCOL37, &c.dbCOL38, &c.dbCOL39, &c.IsFTTH); err != nil {
				fmt.Printf("値の取得に失敗しました。: %v\n", err)
				continue
			}
		}
		if c.dbCOL37.Valid {
			c.COL37 = c.dbCOL37.Time.String()
		}
		if c.dbCOL38.Valid {
			c.COL38 = c.dbCOL38.Time.String()
		}
		if c.dbCOL39.Valid {
			c.COL39 = c.dbCOL39.Time.String()
		}
		cs = append(cs, c)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (para *ParameterCustomer) getFid2BuilNo(fid uint64) (bnos []string) {
	bnos = make([]string, 0)
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	tbl := ViewFmdu
	if AppIdxFmdu < para.AppKind {
		tbl = ViewIntro
	}

	var tno uint64
	sql := fmt.Sprintf("select touroku_no from %v where id = %v limit 1 ;", tbl, fid)
	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		rows.Scan(&tno)
	}
	if tno == 0 {
		return
	}

	sql = fmt.Sprintf("select buil_no from %v where touroku_no = %v ;", tbl, tno)
	if para.Debug {
		fmt.Println(sql)
	}
	rows, err = db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	var b uint64
	for rows.Next() {
		rows.Scan(&b)
		bnos = append(bnos, fmt.Sprintf("%v", b))
	}

	return
}

func (para *ParameterCustomer) UpdateCustomer() (cs Customers, err error) {
	if para.ClientNo == 0 {
		err = fmt.Errorf("*** Error *** ClientNoが不正")
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT COUNT(*) AS cnt FROM %v WHERE col23 = %v;", ViewDensyouCust, para.ClientNo)
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
		err = fmt.Errorf("NoRecord[ClientNo=%v]", para.ClientNo)
		return
	}

	sql = fmt.Sprintf("update %v ", TblDensyouCust)
	tmpsql := ""
	addcnt := 0
	if -1 < para.ProcessStatus {
		tmpsql = fmt.Sprintf(" set PROCESS_STATUS = %v ", para.ProcessStatus)
		addcnt++
	}
	if -1 < para.KoujiStatus {
		if addcnt == 0 {
			tmpsql = fmt.Sprintf(" set KOUJI_STATUS = %v ", para.KoujiStatus)
		} else {
			tmpsql = fmt.Sprintf("%v, KOUJI_STATUS = %v ", tmpsql, para.KoujiStatus)
		}
		addcnt++
	}
	if 0 < len(para.Person) {
		if addcnt == 0 {
			tmpsql = fmt.Sprintf(" set PERSON = '%v' ", para.Person)
		} else {
			tmpsql = fmt.Sprintf("%v, PERSON = '%v' ", tmpsql, para.Person)
		}
		addcnt++
	}
	if 0 < len(para.Remarks) {
		if addcnt == 0 {
			tmpsql = fmt.Sprintf(" set REMARKS1 = '%v' ", para.Remarks)
		} else {
			tmpsql = fmt.Sprintf("%v, REMARKS1 = '%v' ", tmpsql, para.Remarks)
		}
		addcnt++
	}

	if 0 < len(tmpsql) {
		sql = fmt.Sprintf("%v %v", sql, tmpsql)
	}
	footer := fmt.Sprintf(" where CLIENT_NO = %v;", para.ClientNo)
	sql = fmt.Sprintf("%v %v", sql, footer)
	if para.Debug {
		fmt.Println(sql)
	}

	_, err = db.Query(sql)

	if err != nil {
		return
	}

	cs, err = para.getCustomer()

	return
}

const (
	//契約
	//利用中
	//休止
	//解約申出
	KaiyakuWord = "解約"
)

type ServiceSegment struct {
	Code     uint8
	Name     string
	Division string
	IsPaid   uint8
	IsFTTH   uint8
}

type ServiceSegmentMap map[uint8]*ServiceSegment

func (para *Parameter) getServiceSegment() (sm ServiceSegmentMap, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf(`SELECT CODE,NAME,DIVISION,PAID,FTTH FROM %v`, TblServiceSegment)
	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	sm = make(ServiceSegmentMap)
	for rows.Next() {
		s := new(ServiceSegment)
		if err := rows.Scan(&s.Code, &s.Name, &s.Division, &s.IsPaid, &s.IsFTTH); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		sm[s.Code] = s
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (c *Customer) isFTTHOK(codes string, stat string, sm ServiceSegmentMap) (ok bool) {
	ok = true
	if (len(codes) == 0) || (len(stat) == 0) {
		return
	}
	if strings.Contains(stat, KaiyakuWord) {
		return
	}
	code, err := strconv.Atoi(codes)
	if (err != nil) || (code == 0) {
		return
	}
	s := sm[uint8(code)]
	if s.IsFTTH == 0 {
		ok = false
	}

	return
}

func (c *Customer) isFTTH(sm ServiceSegmentMap) (isftth bool) {
	isftth = true

	isftth = c.isFTTHOK(c.COL25.String, c.COL27.String, sm)
	if !isftth {
		return
	}

	isftth = c.isFTTHOK(c.COL28.String, c.COL30.String, sm)
	if !isftth {
		return
	}

	isftth = c.isFTTHOK(c.COL31.String, c.COL33.String, sm)
	if !isftth {
		return
	}

	return
}

func (c *Customer) isTvFtth(sm ServiceSegmentMap) (isftth bool) {
	isftth = true

	isftth = c.isFTTHOK(c.COL25.String, c.COL27.String, sm)
	if !isftth {
		return
	}

	return
}

func (c *Customer) isPaidOK(codes string, stat string, sm ServiceSegmentMap) (ok bool) {
	ok = true
	if (len(codes) == 0) || (len(stat) == 0) {
		return
	}
	if strings.Contains(stat, KaiyakuWord) {
		return
	}
	code, err := strconv.Atoi(codes)
	if (err != nil) || (code == 0) {
		return
	}
	s := sm[uint8(code)]
	if s.IsPaid == 0 {
		ok = false
	}

	return
}

func (c *Customer) isPaid(sm ServiceSegmentMap) (ispaid bool) {
	ispaid = true

	ispaid = c.isPaidOK(c.COL25.String, c.COL27.String, sm)
	if !ispaid {
		return
	}

	ispaid = c.isPaidOK(c.COL28.String, c.COL30.String, sm)
	if !ispaid {
		return
	}

	ispaid = c.isPaidOK(c.COL31.String, c.COL33.String, sm)
	if !ispaid {
		return
	}

	return
}
