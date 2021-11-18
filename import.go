package fmdu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"

	"bitbucket.org/inadenoita/futalib"
)

const (
	CrawlLockFile = "clawl.lock"
	CrawLog       = "crawl.log"

	SystemIP   = "127.0.0.1"
	SystemUser = "root"
)

type ParameterImport struct {
	*Parameter
	//Path          string
	CsvMduSystem  string
	CsvCustSystem string
}

func (para *ParameterImport) isExistLockFile() (exist bool) {
	exist = false
	lfile := filepath.Join(para.CsvPath, CrawlLockFile)
	_, err := os.Stat(lfile)
	if err == nil {
		exist = true
	}
	return
}
func (para *ParameterImport) makeLockFile() (err error) {
	lfile := filepath.Join(para.CsvPath, CrawlLockFile)
	_, err = os.Create(lfile)
	return
}
func (para *ParameterImport) delLockFile() (err error) {
	lfile := filepath.Join(para.CsvPath, CrawlLockFile)
	err = futalib.DelFile(lfile)
	return
}

func (dbpara *ParameterAPI) InitParameterImport(c echo.Context) (para *ParameterImport, err error) {
	para = new(ParameterImport)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.CsvMduSystem = filepath.Join(para.CsvPath, SystemCsvMduFileName)
	para.CsvCustSystem = filepath.Join(para.CsvPath, SystemCsvCustFileName)

	if len(para.IP) == 0 {
		para.Parameter.IP = SystemIP
	}
	if len(para.UserID) == 0 {
		para.Parameter.UserID = SystemUser
	}

	return
}

// ImportSmsData ... SMS-CSV import処理
func (para *ParameterImport) ImportSmsData() (idx int, err error) {
	idx = 0
	cnt := 0
	if para.isExistLockFile() {
		err = fmt.Errorf("Crawl実行中")
		return
	}
	err = para.makeLockFile()
	if err != nil {
		return
	}

	err = para.CopyCsvFromSMS()
	if err != nil {
		para.delLockFile()
		return
	}

	_, err = os.Stat(para.CsvMdu)
	if err != nil {
		para.delLockFile()
		return
	}
	_, err = os.Stat(para.CsvCust)
	if err != nil {
		para.delLockFile()
		return
	}

	err = para.modfCsvMdu()
	if err != nil {
		para.delLockFile()
		return
	}
	err = para.modfCsvCust()
	if err != nil {
		para.delLockFile()
		return
	}

	//  *** MDU CSV Import ***
	cnt, err = para.importMdu()
	if err != nil {
		para.delLockFile()
		return
	}
	idx += cnt
	//  ADD 2021/10/25
	err = para.setBuilInfoFromTourokuNo()
	if err != nil {
		para.delLockFile()
		return
	}

	//  *** CUST CSV Import ***
	cnt, err = para.importCust()
	if err != nil {
		para.delLockFile()
		return
	}
	idx += cnt

	// カプラ履歴テーブル更新
	err = para.UpdateCplHistory()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
	}

	// FMDU Table Insert
	cnt, err = para.InsertFmdu()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
		return
	}
	fmt.Printf("Insert Record for %v : %v count\n", TblFmdu, cnt)
	idx += cnt

	// Intro Table Insert
	cnt, err = para.InsertIntro()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
		return
	}
	fmt.Printf("Insert Record for %v : %v count\n", TblIntro, cnt)
	idx += cnt

	// Densyou Table Insert
	cnt, err = para.InsertDensyou()
	if err != nil {
		fmt.Println(err)
		para.delLockFile()
		return
	}
	fmt.Printf("Insert Record for %v : %v count\n", TblDensyou, cnt)
	idx += cnt

	// DensyouCustomer Table Insert
	cnt, err = para.InsertDensyouCust()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
		return
	}
	fmt.Printf("Insert Record for %v : %v count\n", TblDensyouCust, cnt)
	idx += cnt

	para.delLockFile()

	return
}

func (para ParameterAPI) InsertFmdu() (idx int, err error) {
	ms, err := para.getMdus(1)
	if err != nil {
		return
	}
	idx, err = para.insertFmdu(ms)

	return
}

func (para ParameterAPI) insertFmdu(ms Mdus) (idx int, err error) {
	idx = 0
	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(MaxOpenConns)

	// -- Modf 2021/10/25 登録NOもセットするように拡張
	for _, m := range ms {
		sql := fmt.Sprintf("SELECT buil_no,touroku_no FROM %v where buil_no=%v;", TblFmdu, m.BuilNo)
		//fmt.Println(sql)
		rows, err := db.Query(sql)
		if err != nil {
			//log.Fatal(err)
			return idx, err
		}
		defer rows.Close()
		tmpidx := 0
		for rows.Next() {
			m := new(Mdu)
			if err := rows.Scan(&m.BuilNo, &m.TourokuNo); err != nil {
				fmt.Printf("値の取得に失敗しました。: %v\n", err)
			}
			tmpidx++
		}
		if 0 < tmpidx {
			continue
		}

		var fid uint64
		err = db.QueryRow("INSERT INTO TBL_FMDU(buil_no,touroku_no,ctim) VALUES($1, $2, now()) RETURNING id", m.BuilNo, &m.TourokuNo).Scan(&fid)
		if err != nil {
			fmt.Println(err)
		}
		idx++
	}
	return
}

func (para ParameterAPI) InsertIntro() (idx int, err error) {
	// FTTHも含み全て取得
	ms, err := para.getMdus()
	if err != nil {
		return
	}
	idx, err = para.insertIntro(ms)

	return
}

func (para ParameterAPI) insertIntro(ms Mdus) (idx int, err error) {
	idx = 0
	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(MaxOpenConns)

	// -- Modf 2021/10/25 登録NOもセットするように拡張
	for _, m := range ms {
		sql := fmt.Sprintf("SELECT buil_no,touroku_no FROM %v where buil_no=%v;", TblIntro, m.BuilNo)
		//fmt.Println(sql)
		rows, err := db.Query(sql)
		if err != nil {
			//log.Fatal(err)
			return idx, err
		}
		defer rows.Close()
		tmpidx := 0
		for rows.Next() {
			m := new(Mdu)
			if err := rows.Scan(&m.BuilNo, &m.TourokuNo); err != nil {
				fmt.Printf("値の取得に失敗しました。: %v\n", err)
			}
			tmpidx++
		}
		if 0 < tmpidx {
			continue
		}

		var fid uint64
		err = db.QueryRow("INSERT INTO TBL_Intro(buil_no,touroku_no,ctim) VALUES($1, $2, now()) RETURNING id", m.BuilNo, m.TourokuNo).Scan(&fid)
		if err != nil {
			fmt.Println(err)
		}
		idx++
	}
	return
}

func (para Parameter) InsertDensyou() (idx int, err error) {
	ds, err := para.getDensyouListFromSMS()
	if err != nil {
		return
	}
	idx, err = para.insertDensyou(ds)

	return
}

func (para Parameter) insertDensyou(ds Densyous) (idx int, err error) {
	idx = 0
	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(MaxOpenConns)

	for _, d := range ds {
		sql := fmt.Sprintf("SELECT id FROM %v where id=%v;", TblDensyou, d.ID)
		//fmt.Println(sql)
		rows, err := db.Query(sql)
		if err != nil {
			//log.Fatal(err)
			return idx, err
		}
		defer rows.Close()
		tmpidx := 0
		for rows.Next() {
			d := new(Densyou)
			if err := rows.Scan(&d.ID); err != nil {
				fmt.Printf("値の取得に失敗しました。: %v\n", err)
			}
			tmpidx++
		}
		if 0 < tmpidx {
			continue
		}

		var fid uint64
		err = db.QueryRow("INSERT INTO TBL_DENSYOU(id,name,ctim) VALUES($1, $2, now()) RETURNING id", d.ID, d.Name).Scan(&fid)
		if err != nil {
			fmt.Println(err)
		}
		idx++
	}
	return
}

func (para Parameter) InsertDensyouCust() (idx int, err error) {
	cs, err := para.getDensyouCustFromSMS()
	if err != nil {
		return
	}
	idx, err = para.insertDensyouCust(cs)

	return
}

func (para Parameter) insertDensyouCust(cs densyouCusts) (idx int, err error) {
	idx = 0
	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(MaxOpenConns)

	for _, c := range cs {
		sql := fmt.Sprintf("SELECT id FROM %v where CLIENT_NO = %v;", TblDensyouCust, c.clientNo)
		//fmt.Println(sql)
		rows, err := db.Query(sql)
		if err != nil {
			//log.Fatal(err)
			return idx, err
		}
		defer rows.Close()
		tmpidx := 0
		for rows.Next() {
			//d := new(Densyou)
			var id uint64
			if err := rows.Scan(&id); err != nil {
				fmt.Printf("値の取得に失敗しました。: %v\n", err)
			}
			tmpidx++
		}
		if 0 < tmpidx {
			continue
		}

		var fid uint64
		err = db.QueryRow("INSERT INTO TBL_DENSYOU_CUSTOMER(densyou_code,client_no,ctim) VALUES($1, $2, now()) RETURNING id",
			c.densyouCode, c.clientNo).Scan(&fid)
		if err != nil {
			fmt.Println(err)
		}
		idx++
	}
	return
}
