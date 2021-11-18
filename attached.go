package fmdu

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	TblAttached         = "TBL_ATTACHED"
	ViewFmduAttached    = "VIEW_FMDU_ATTACHED"
	ViewIntroAttached   = "VIEW_INTRO_ATTACHED"
	ViewDensyouAttached = "VIEW_DENSYOU_ATTACHED"
)

type ParameterAttached struct {
	*Parameter
	Kind      string
	AppKind   uint8
	ID        uint64
	FID       uint64
	TourokuNo uint64
	BuilNo    uint64
	BuilName  string
	Path      string
	Name      string
	Title     string
	Remarks   string
	Table     string
}

func (dbpara *ParameterAPI) InitParameterAttached(c echo.Context) (para *ParameterAttached, err error) {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	kind := c.QueryParam("kind")
	if len(kind) == 0 {
		kind = DirAttached
	}
	app, _ := strconv.Atoi(c.QueryParam("app"))
	fid, _ := strconv.Atoi(c.QueryParam("fid"))
	tourokuno, _ := strconv.Atoi(c.QueryParam("tourokuno"))
	builno, _ := strconv.Atoi(c.QueryParam("builno"))
	builname := c.QueryParam("builname")
	path := c.QueryParam("path")
	name := c.QueryParam("name")
	title := c.QueryParam("title")
	remarks := c.QueryParam("remarks")

	para = new(ParameterAttached)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.Kind = kind
	para.AppKind = uint8(app)
	para.ID = uint64(id)
	para.FID = uint64(fid)
	para.TourokuNo = uint64(tourokuno)
	para.BuilNo = uint64(builno)
	para.BuilName = builname
	para.Path = path
	para.Name = name
	para.Title = title
	para.Remarks = remarks

	tbl := ""
	switch para.AppKind {
	case AppIdxIntro:
		tbl = ViewIntroAttached
	case AppIdxDensyou:
		tbl = ViewDensyouAttached
	default:
		tbl = ViewFmduAttached
	}
	para.Table = tbl

	return
}

type Attached struct {
	ID        uint64
	FID       uint64
	AppKind   uint8
	Path      string
	Name      string
	Title     sql.NullString
	Remarks   sql.NullString
	IP        sql.NullString
	UserID    sql.NullString
	Ctim      time.Time
	TourokuNo uint64
	BuilNo    uint64
	BuilName  string
}
type Attacheds []*Attached

// GetAttached ... 添付ファイルデータ取得
func (para *ParameterAttached) GetAttached() (as Attacheds, err error) {

	if len(para.Table) == 0 {
		err = fmt.Errorf("AppKind is Null")
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,app_kind,fid,touroku_no,buil_no,buil_name,path,name,title,remarks,ip,user_id,ctim FROM %v", para.Table)
	tmp := ""
	addcnt := 0

	if 0 < para.ID {
		if addcnt == 0 {
			tmp = fmt.Sprintf(" (id = %v) ", para.ID)
		} else {
			tmp = fmt.Sprintf(" %v and (id = %v) ", tmp, para.ID)
		}
		addcnt++
	}
	if 0 < para.TourokuNo {
		if addcnt == 0 {
			tmp = fmt.Sprintf(" (touroku_no = %v) ", para.TourokuNo)
		} else {
			tmp = fmt.Sprintf(" %v and (touroku_no = %v) ", tmp, para.TourokuNo)
		}
		addcnt++
	}
	if 0 < para.BuilNo {
		if addcnt == 0 {
			tmp = fmt.Sprintf(" (buil_no = %v) ", para.BuilNo)
		} else {
			tmp = fmt.Sprintf(" %v and (buil_no = %v) ", tmp, para.BuilNo)
		}
		addcnt++
	}
	if 0 < len(para.BuilName) {
		if addcnt == 0 {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(buil_name) like %v) ", getNoCaseTextSQLStr(para.BuilName))
		} else {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(buil_name) like %v) ", tmp, getNoCaseTextSQLStr(para.BuilName))
		}
		addcnt++
	}
	if 0 < len(para.Name) {
		if addcnt == 0 {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(name) like %v) ", getNoCaseTextSQLStr(para.Name))
		} else {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(name) like %v) ", tmp, getNoCaseTextSQLStr(para.Name))
		}
		addcnt++
	}
	if 0 < len(para.Title) {
		if addcnt == 0 {
			tmp = fmt.Sprintf(" (GET_NOCASE_TEXT(title) like %v) ", getNoCaseTextSQLStr(para.Title))
		} else {
			tmp = fmt.Sprintf(" %v and (GET_NOCASE_TEXT(title) like %v) ", tmp, getNoCaseTextSQLStr(para.Title))
		}
		addcnt++
	}
	if 0 < len(para.Remarks) {
		if addcnt == 0 {
			tmp = fmt.Sprintf(" (remarks like '%%%v%%') ", para.Remarks)
		} else {
			tmp = fmt.Sprintf(" %v and (remarks like '%%%v%%') ", tmp, para.Remarks)
		}
		addcnt++
	}
	if addcnt == 0 {
		sql = fmt.Sprintf("%v order by id desc;", sql)
	} else {
		sql = fmt.Sprintf("%v where %v order by id desc;", sql, tmp)
	}

	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	as = make(Attacheds, 0)
	for rows.Next() {
		a := new(Attached)
		if err := rows.Scan(&a.ID, &a.AppKind, &a.FID, &a.TourokuNo, &a.BuilNo, &a.BuilName, &a.Path, &a.Name, &a.Title.String,
			&a.Remarks.String, &a.IP.String, &a.UserID.String, &a.Ctim); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		as = append(as, a)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

const (
	DataPath = "data"

	DirAttached = "attached"
	//DirDensyou = "densyou"
)

func checkAttachedKind(kind string) (ok bool) {
	ok = false
	switch kind {
	case DirAttached:
		ok = true
	default:
	}

	return ok
}

func (para *ParameterAttached) getDstPath() (dstpath string) {
	dstpath = filepath.Join(para.AppPath, DataPath, para.Kind)
	return
}

// InsertAttached ... 添付ファイルのレコード追加
func (para *ParameterAttached) InsertAttached() (as Attacheds, err error) {
	if len(para.AppPath) == 0 {
		err = fmt.Errorf("*** Error *** AppPathが不正")
		return
	}
	if para.AppKind == 0 {
		err = fmt.Errorf("*** Error *** AppKindが不正")
		return
	}
	if para.FID == 0 {
		err = fmt.Errorf("*** Error *** FIDが不正")
		return
	}
	if para.TourokuNo == 0 {
		err = fmt.Errorf("*** Error *** TourokuNoが不正")
		return
	}
	if len(para.Path) == 0 {
		err = fmt.Errorf("*** Error *** Pathが不正")
		return
	}
	if len(para.Name) == 0 {
		err = fmt.Errorf("*** Error *** FileNameが不正")
		return
	}
	ok := checkAttachedKind(para.Kind)
	if !ok {
		err = fmt.Errorf("*** Error *** Kind(データ種)が不正")
		return
	}

	// ファイルコピー
	dstpath := getDstFilePath(para.getDstPath(), filepath.Ext(para.Name), para.TourokuNo)
	err = copyFile(para.Path, dstpath)
	if err != nil {
		return
	}
	para.Path = dstpath

	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	var id uint
	err = db.QueryRow("INSERT INTO TBL_ATTACHED (app_kind,touroku_no,fid,path,name,title,remarks,ip,user_id,ctim ) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,now()) RETURNING id",
		para.AppKind, para.TourokuNo, para.FID, para.Path, para.Name, para.Title, para.Remarks, para.IP, para.UserID).Scan(&id)

	if err != nil {
		deleteFile(para.Path)
		return
	}
	fmt.Println(id)
	para.ID = uint64(id)

	as, err = para.GetAttached()

	return
}

// UpdateAttached ... 添付ファイルのレコード更新
func (para *ParameterAttached) UpdateAttached() (as Attacheds, err error) {
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

	sql := fmt.Sprintf("SELECT id FROM %v where (id = %v)", para.Table, para.ID)
	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	idx := 0
	var a *Attached
	for rows.Next() {
		a = new(Attached)
		if err := rows.Scan(&a.ID); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		idx++
		break
	}
	if err = rows.Err(); err != nil {
		return
	}
	if idx < 1 {
		err = fmt.Errorf("*** Error *** レコードが存在しない[%v]", para.ID)
		return
	}
	now := time.Now()
	_, err = db.Query("UPDATE TBL_ATTACHED SET remarks = $1, ip = $2, user_id = $3, ctim = $4 "+
		" WHERE ID = $5",
		para.Remarks, para.IP, para.UserID, now, para.ID)
	if err != nil {
		return
	}

	as, err = para.GetAttached()

	return
}

// DeleteAttached ... 添付ファイルの削除
func (para *ParameterAttached) DeleteAttached() (as Attacheds, err error) {
	if para.ID == 0 {
		err = fmt.Errorf("*** Error *** IDが不正")
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT path FROM %v where (id = %v)", para.Table, para.ID)
	rows, err := db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	path := ""
	var a *Attached
	for rows.Next() {
		a = new(Attached)
		if err := rows.Scan(&a.Path); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		path = a.Path
		break
	}
	if err = rows.Err(); err != nil {
		return
	}
	if len(path) < 1 {
		err = fmt.Errorf("*** Error *** パスが取得できない[%v]", para.ID)
		return
	}

	// ファイル削除
	err = deleteFile(path)
	if err != nil {
		fmt.Println(err)
	}

	// レコード削除
	_, err = db.Query("DELETE FROM TBL_ATTACHED WHERE ID = $1", para.ID)
	if err != nil {
		return
	}

	return
}
