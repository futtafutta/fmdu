package fmdu

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	TblFmdu     = "TBL_FMDU"
	ViewFmdu    = "VIEW_FMDU"
	ViewCplCust = "VIEW_FMDU_CPL_CUST"
)

type ParameterFmdu struct {
	*Parameter
	ID           uint64
	Ftype        uint8
	TourokuNo    uint64 // ADD 2021/10/28
	BuilNo       uint64
	BuilAddr     string
	BuilName     string
	Remarks      string
	BuilNoStr    string
	TourokuNoStr string // ADD 2021/10/28
	FtthArea     string // ADD 2021/11/06
	IsCpl        int8
	RoomCnt      int8
	Eq           uint8
}

func (dbpara *ParameterAPI) InitParameterFmdu(c echo.Context) (para *ParameterFmdu, err error) {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	ftype, _ := strconv.Atoi(c.QueryParam("ftype"))
	builno, _ := strconv.Atoi(c.QueryParam("builno"))
	tourokuno, _ := strconv.Atoi(c.QueryParam("tourokuno"))
	builnostr := c.QueryParam("builnostr")
	builaddr, _ := decodeURLString(c.QueryParam("builaddr"))
	builname, _ := decodeURLString(c.QueryParam("builname"))
	remarks1 := c.QueryParam("remarks")
	tourokunostr := c.QueryParam("tourokunostr")
	fttharea := c.QueryParam("fttharea")
	iscpl := -1
	tmp := c.QueryParam("iscpl")
	if 0 < len(tmp) {
		iscpl, _ = strconv.Atoi(tmp)
	}
	roomcnt := -1
	tmp = c.QueryParam("roomcnt")
	if 0 < len(tmp) {
		roomcnt, _ = strconv.Atoi(tmp)
	}

	eq, _ := strconv.Atoi(c.QueryParam("eq"))

	para = new(ParameterFmdu)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.ID = uint64(id)
	para.Ftype = uint8(ftype)
	para.TourokuNo = uint64(tourokuno)
	para.BuilNo = uint64(builno)
	para.BuilAddr = builaddr
	para.BuilName = builname
	para.Remarks = remarks1
	para.BuilNoStr = builnostr
	para.TourokuNoStr = tourokunostr
	para.FtthArea = fttharea
	para.IsCpl = int8(iscpl)
	para.RoomCnt = int8(roomcnt)
	para.Eq = uint8(eq)

	return
}

type Fmdu struct {
	ID       uint64
	FtypeID  uint8
	Remarks1 sql.NullString
	Remarks2 sql.NullString
	Remarks3 sql.NullString
	IP       sql.NullString
	UserID   sql.NullString
	Ctim     time.Time
	CplCnt   uint
	//CplPortCntDB uint64
	CplPortCnt  int64
	CplUseCnt   int64
	CplEmptyCnt int64
	TourokuNo   uint64
	BuilNo      uint64
	BuilAddr    string
	BuilName    string
	FtthArea    string
	BuilCnt     sql.NullInt32
	BuilStatus  sql.NullString
	BuilTV      sql.NullString
	BuilNet     sql.NullString
	Longitude   sql.NullFloat64
	Latitude    sql.NullFloat64
}

type Fmdus []*Fmdu

func (para *ParameterFmdu) GetFmdu() (fs Fmdus, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,remarks1,user_id,ctim,cpl_cnt,cpl_port_cnt,touroku_no,buil_no,buil_addr,buil_name,buil_cnt,buil_status,buil_tv,buil_net,ftth_area FROM %v ", ViewFmdu)
	addcnt := 0
	tmp := ""
	if 0 < para.ID {
		tmp = fmt.Sprintf(" (id = %v) ", para.ID)
		addcnt++
	}
	if 0 < para.TourokuNo {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (touroku_no = %v) ", tmp, para.TourokuNo)
		} else {
			tmp = fmt.Sprintf(" (touroku_no = %v) ", para.TourokuNo)
		}
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
			tmp = fmt.Sprintf(" %v and (buil_no_str like '%%%v%%') ", tmp, para.BuilNoStr)
		} else {
			tmp = fmt.Sprintf(" (buil_no_str like '%%%v%%') ", para.BuilNoStr)
		}
		addcnt++
	}
	if 0 < len(para.BuilAddr) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and ( GET_NOCASE_TEXT(buil_addr) like %v ) ", tmp, getNoCaseTextSQLStr(para.BuilAddr))
		} else {
			tmp = fmt.Sprintf(" ( GET_NOCASE_TEXT(buil_addr) like %v ) ", getNoCaseTextSQLStr(para.BuilAddr))
		}
		addcnt++
	}
	if 0 < len(para.BuilName) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and ( GET_NOCASE_TEXT(buil_name) like %v ) ", tmp, getNoCaseTextSQLStr(para.BuilName))
		} else {
			tmp = fmt.Sprintf(" ( GET_NOCASE_TEXT(buil_name) like %v ) ", getNoCaseTextSQLStr(para.BuilName))
		}
		addcnt++
	}
	if 0 < len(para.TourokuNoStr) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and (touroku_no_str like '%%%v%%') ", tmp, para.TourokuNoStr)
		} else {
			tmp = fmt.Sprintf(" (touroku_no_str like '%%%v%%') ", para.TourokuNoStr)
		}
		addcnt++
	}
	if 0 < len(para.FtthArea) {
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and ( GET_NOCASE_TEXT(ftth_area) like %v ) ", tmp, getNoCaseTextSQLStr(para.FtthArea))
		} else {
			tmp = fmt.Sprintf(" ( GET_NOCASE_TEXT(ftth_area) like %v ) ", getNoCaseTextSQLStr(para.FtthArea))
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
	if -1 < para.IsCpl {
		ctmp := ""
		if para.IsCpl == 0 {
			ctmp = " (cpl_cnt < 1) "
		} else {
			ctmp = " (0 < cpl_cnt) "
		}
		if 0 < addcnt {
			tmp = fmt.Sprintf(" %v and %v ", tmp, ctmp)
		} else {
			tmp = ctmp
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

	fs = make(Fmdus, 0)
	for rows.Next() {
		f := new(Fmdu)
		if err := rows.Scan(&f.ID, &f.Remarks1, &f.UserID, &f.Ctim, &f.CplCnt, &f.CplPortCnt,
			&f.TourokuNo, &f.BuilNo, &f.BuilAddr, &f.BuilName, &f.BuilCnt, &f.BuilStatus, &f.BuilTV, &f.BuilNet, &f.FtthArea,
		); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		//sql = fmt.Sprintf("SELECT COUNT(id) FROM %v where id = %v", ViewCplCust, f.ID)
		/*
			if para.Debug {
				fmt.Println(sql)
			}
		*/
		var cnt int64 = 0
		err = db.QueryRow("SELECT COUNT(*) AS cnt FROM VIEW_FMDU_CPL_CUST WHERE id = $1;", f.ID).Scan(&cnt)
		if err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
		}
		//fmt.Println(cnt)
		//f.CplPortCnt = uint(f.CplPortCntDB.Int64)
		//f.CplEmptyCnt = f.CplPortCnt - uint(cnt)
		f.CplUseCnt = cnt
		f.CplEmptyCnt = f.CplPortCnt - cnt
		if f.CplEmptyCnt < 0 {
			f.CplEmptyCnt = 0
		}
		f.Longitude, f.Latitude = llmap.getLngLat(f.BuilNo)

		fs = append(fs, f)
	}
	if err = rows.Err(); err != nil {
		return
	}

	//portused := 0

	return
}

func (para *ParameterFmdu) PostFmdu() (fs Fmdus, err error) {
	if para.BuilNo == 0 {
		err = fmt.Errorf("*** Error *** 建物番号が不正")
		return
	}
	if para.TourokuNo == 0 {
		err = fmt.Errorf("*** Error *** 登録番号が不正")
		return
	}
	if para.Ftype < 1 {
		err = fmt.Errorf("*** Error *** 集合種別が不正")
		return
	}
	/*
		if para.BrCnt < 1 {
			err = fmt.Errorf("*** Error *** カプラ分岐数が不正")
			return
		}
	*/
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	var fid uint
	//var now = time.Now()
	err = db.QueryRow("INSERT INTO TBL_FMDU (buil_no,touroku_no,remarks,ip,user_id,ctim ) VALUES($1,$2,$3,$4,$5,now()) RETURNING id",
		para.BuilNo, para.TourokuNo, para.Remarks, para.IP, para.UserID).Scan(&fid)
	if err != nil {
		return
	}
	fmt.Println(fid)
	para.ID = uint64(fid)

	fs, err = para.GetFmdu()

	return
}
