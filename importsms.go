package fmdu

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/inadenoita/futalib"
	"github.com/davecgh/go-spew/spew"
)

const (
	cancelKeyword = "集合解約"
)

const (
	DelimOld   = "\",\""
	DelimNew   = "|"
	ColCntMdu  = 31
	ColCntCust = 39

	SystemCsvMduFileName  = "mdu.csv"
	SystemCsvCustFileName = "cust.csv"
)

func (para *ParameterImport) CopyCsvFromSMS() (err error) {
	_, err = os.Stat(para.SmsPath)
	if err != nil {
		return
	}
	csvs := futalib.GetExtFiles(&[]string{}, para.SmsPath, para.SmsPath, ".csv")
	csvcnt := len(*csvs)
	if csvcnt == 0 {
		err = fmt.Errorf("csvファイルがない [%v]", para.SmsPath)
		return
	}

	for i, csv := range *csvs {
		_, f := filepath.Split(csv)
		dstpath := filepath.Join(para.CsvPath, f)
		fmt.Printf("[%v / %v] copy  %v  to  %v \n", i+1, csvcnt, csv, dstpath)
		err = futalib.CopyFile(csv, dstpath)
		if err != nil {
			fmt.Println(err)
		}
	}

	return
}

// 集合CSVのレコード行以外を削除
func (para *ParameterImport) modfCsvMdu() (err error) {
	newlist := make([]string, 0)
	outpath := para.CsvMduSystem
	_, err = os.Stat(outpath)
	if err == nil {
		err = futalib.DelFile(outpath)
		if err != nil {
			return
		}
	}
	err = futalib.NkfExe("nkf", para.CsvMdu, futalib.ModeSjis2Utf)

	list, _ := futalib.Csv2Slice(para.CsvMdu)
	for _, l := range *list {
		ls := strings.Split(l, ",")
		if len(ls) < ColCntMdu {
			continue
		}
		l = strings.ReplaceAll(l, DelimOld, DelimNew)
		l = strings.ReplaceAll(l, "\"", "")
		newlist = append(newlist, l)
	}

	if !futalib.OutputCsvFromArrayUtf8(newlist, outpath) {
		err = fmt.Errorf("CSVファイル出力失敗[%v]", outpath)
		return
	}

	return
}

// 顧客CSVのレコード行以外を削除
func (para *ParameterImport) modfCsvCust() (err error) {
	newlist := make([]string, 0)
	outpath := para.CsvCustSystem
	_, err = os.Stat(outpath)
	if err == nil {
		err = futalib.DelFile(outpath)
		if err != nil {
			return
		}
	}
	err = futalib.NkfExe("nkf", para.CsvCust, futalib.ModeSjis2Utf)

	list, _ := futalib.Csv2Slice(para.CsvCust)
	for _, l := range *list {
		ls := strings.Split(l, ",")
		if len(ls) < ColCntCust {
			continue
		}
		l = strings.ReplaceAll(l, DelimOld, DelimNew)
		l = strings.ReplaceAll(l, "\"", "")
		newlist = append(newlist, l)
	}

	if !futalib.OutputCsvFromArrayUtf8(newlist, outpath) {
		err = fmt.Errorf("CSVファイル出力失敗[%v]", outpath)
		return
	}

	return
}

//  *** MDU CSV Import ***
func (para *ParameterImport) importMdu() (idx int, err error) {
	idx = 0
	imppara := initParameterImport()
	imppara.tbl = TblMDU
	imppara.delimiter = DelimNew
	imppara.isheader = true
	imppara.column = "(COL01,COL02,COL03,COL04,COL05,COL06,COL07,COL08,COL09,COL10,COL11,COL12,COL13,COL14,COL15,COL16,COL17,COL18,COL19,COL20,COL21,COL22,COL23,COL24,COL25,COL26,COL27,COL28,COL29,COL30,COL31)"
	imppara.csv = para.CsvMduSystem
	para.ParameterRDB.parameterImport = imppara
	spew.Dump(para.ParameterRDB)
	cnt, err := para.ParameterRDB.execCopySQL()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
		return
	}
	idx += cnt

	// 解約削除
	cmd := fmt.Sprintf("DELETE FROM %v WHERE COL06 = '%v'", para.tbl, cancelKeyword)
	cnt, err = para.execSQL(cmd)
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
		return
	}
	fmt.Printf("Record for %v : %v count\n", para.tbl, cnt)

	idx -= cnt

	fmt.Printf("Import Record for %v : %v count\n", imppara.tbl, idx)

	return
}

type mduinfo struct {
	COL03 sql.NullInt64
	COL07 sql.NullString
	COL08 sql.NullInt64
	COL09 sql.NullString
	COL10 sql.NullInt64
	COL11 sql.NullString
	COL12 sql.NullInt64
}

// 建物番号違いの集合情報をセット -- ADD 2021/10/25
func (para *Parameter) setBuilInfoFromTourokuNo() (err error) {
	ms, err := para.getMdus()
	if err != nil {
		return
	}
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(MaxOpenConns)

	for _, m := range ms {
		if m.BuilCnt.Valid {
			continue
		}
		var mi *mduinfo
		cmd := fmt.Sprintf("select distinct on (col03,col07,col08,col09,col10,col11,col12) col03,col07,col08,col09,col10,col11,col12 from %v where col24 = %v and (0<=col03);", TblMDU, m.TourokuNo)
		if para.Debug {
			fmt.Println(cmd)
		}
		rows, tmperr := db.Query(cmd)
		if tmperr != nil {
			fmt.Println(tmperr)
			continue
		}
		defer rows.Close()

		idx := 0
		for rows.Next() {
			mi = new(mduinfo)
			if err := rows.Scan(&mi.COL03, &mi.COL07, &mi.COL08, &mi.COL09, &mi.COL10, &mi.COL11, &mi.COL12); err != nil {
				fmt.Printf("値の取得に失敗しました。: %v\n", err)
				continue
			}
			idx++
			if 1 <= idx {
				break
			}
		}

		if idx == 0 {
			continue
		}
		if !mi.COL03.Valid {
			continue
		}
		cmd = fmt.Sprintf(`update %v set COL03 = %v, COL07= '%v', COL08 = %v, COL09 = '%v', COL10 = %v, COL11 = '%v', COL12 = %v WHERE (COL30 = %v)`,
			TblMDU, mi.COL03.Int64, mi.COL07.String, mi.COL08.Int64, mi.COL09.String, mi.COL10.Int64, mi.COL11.String, mi.COL12.Int64, m.BuilNo)
		if para.Debug {
			fmt.Println(cmd)
		}
		_, err = db.Query(cmd)

		if err != nil {
			return
		}

	}

	return
}

//  *** CUST CSV Import ***
func (para *ParameterImport) importCust() (idx int, err error) {
	idx = 0
	imppara := initParameterImport()
	imppara.tbl = TblCustomer
	imppara.delimiter = DelimNew
	imppara.isheader = true
	imppara.column = `(COL01,COL02,COL03,COL04,COL05,COL06,COL07,COL08,COL09,COL10,COL11,COL12,COL13,COL14,COL15,COL16,COL17,COL18,COL19,COL20,COL21,COL22,COL23,COL24,COL25,COL26,COL27,COL28,COL29,COL30,COL31,COL32,COL33,COL34,COL35,COL36,COL37,COL38,COL39)`
	imppara.csv = para.CsvCustSystem
	para.ParameterRDB.parameterImport = imppara
	spew.Dump(para.ParameterRDB)
	cnt, err := para.ParameterRDB.execCopySQL()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
		return
	}
	idx += cnt

	// 顧客番号 is Null（未加入者） 削除
	cmd := fmt.Sprintf("DELETE FROM %v WHERE COL23 is NULL", para.tbl)
	cnt, err = para.execSQL(cmd)
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
		return
	}
	fmt.Printf("Record for %v : %v count\n", para.tbl, cnt)

	idx -= cnt
	fmt.Printf("Import Record for %v : %v count\n", para.tbl, idx)

	// ISFTTHセット
	err = para.setIsFTTH()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
	}
	// ISFTTH TVセット
	err = para.setIsFTTHTv()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
	}
	// ISPaidセット
	err = para.setIsPaid()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
	}
	// 電障のみセット
	err = para.setIsDensyouOnly()
	if err != nil {
		para.delLockFile()
		fmt.Println(err)
	}

	return
}

func (para *Parameter) setIsFTTH() (err error) {
	cs, err := para.getCustomer()
	if err != nil {
		return
	}

	sm, err := para.getServiceSegment()
	if err != nil {
		return
	}

	fmt.Println(len(cs))

	cntftth := 0
	cnthfc := 0
	list := make([]string, 0)
	for _, c := range cs {
		if !c.isFTTH(sm) {
			cnthfc++
			continue
		}
		cntftth++

		custno := fmt.Sprintf("%v", c.COL23.Int64)

		list = append(list, custno)
	}

	if len(list) == 0 {
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(MaxOpenConns)

	liststr := strings.Join(list, ",")

	sql := fmt.Sprintf("UPDATE TBL_CUSTOMER set FTTH = 1 where COL23 IN ( %v );", liststr)
	fmt.Println(sql)

	_, err = db.Query(sql)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cntftth)
	fmt.Println(cnthfc)

	return
}

func (para *Parameter) setIsFTTHTv() (err error) {
	cs, err := para.getCustomer()
	if err != nil {
		return
	}

	sm, err := para.getServiceSegment()
	if err != nil {
		return
	}

	fmt.Println(len(cs))

	cnttvftth := 0
	list := make([]string, 0)
	for _, c := range cs {
		if !c.isTvFtth(sm) {
			continue
		}
		cnttvftth++

		custno := fmt.Sprintf("%v", c.COL23.Int64)

		list = append(list, custno)
	}

	if len(list) == 0 {
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(MaxOpenConns)

	liststr := strings.Join(list, ",")

	sql := fmt.Sprintf("UPDATE TBL_CUSTOMER set FTTH_TV = 1 where COL23 IN ( %v );", liststr)
	fmt.Println(sql)

	_, err = db.Query(sql)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cnttvftth)

	return
}

func (para *Parameter) setIsPaid() (err error) {
	cs, err := para.getCustomer()
	if err != nil {
		return
	}

	sm, err := para.getServiceSegment()
	if err != nil {
		return
	}

	fmt.Println(len(cs))

	cntpaid := 0
	cntfree := 0

	list := make([]string, 0)
	for _, c := range cs {
		if !c.isPaid(sm) {
			cntfree++
			continue
		}
		cntpaid++

		custno := fmt.Sprintf("%v", c.COL23.Int64)

		list = append(list, custno)
	}

	if len(list) == 0 {
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(MaxOpenConns)

	liststr := strings.Join(list, ",")

	sql := fmt.Sprintf("UPDATE TBL_CUSTOMER set PAID = 1 where COL23 IN ( %v );", liststr)
	fmt.Println(sql)

	_, err = db.Query(sql)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cntpaid)
	fmt.Println(cntfree)

	return
}

func (para *Parameter) setIsDensyouOnly() (err error) {
	cs, err := para.getCustomer(true)
	if err != nil {
		return
	}
	fmt.Println(len(cs))
	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(MaxOpenConns)

	cntdenonly := 0
	cntpaid := 0

	list := make([]string, 0)
	for _, c := range cs {
		sql := fmt.Sprintf("SELECT count(*) as cnt FROM %v where COL01 = %v;", TblCustomer, c.COL01)
		tmprows, err := db.Query(sql)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer tmprows.Close()
		var cnt int
		for tmprows.Next() {
			if err := tmprows.Scan(&cnt); err != nil {
				fmt.Printf("値の取得に失敗しました。: %v\n", err)
				cnt = 0
			}
		}
		if 1 < cnt {
			cntpaid++
			continue
		}
		cntdenonly++
		custno := fmt.Sprintf("%v", c.COL23.Int64)

		list = append(list, custno)

	}

	if len(list) == 0 {
		return
	}

	liststr := strings.Join(list, ",")

	sql := fmt.Sprintf("UPDATE TBL_CUSTOMER set DENSYOU_ONLY = 1 where COL23 IN ( %v );", liststr)
	fmt.Println(sql)

	_, err = db.Query(sql)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cntdenonly)
	fmt.Println(cntpaid)

	return
}

func (para *Parameter) getCustomer(isdensyou ...bool) (cs Customers, err error) {
	densyouflag := false
	if 0 < len(isdensyou) {
		densyouflag = isdensyou[0]
	}
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	tbl := TblCustomer
	if densyouflag {
		tbl = ViewCustDensyou
	}
	sql := fmt.Sprintf(`SELECT COL01,COL02,COL03,COL04,COL05,COL06,COL07,COL10,COL11,COL12,COL13,COL14,COL15,COL16,COL17,COL18,COL19,COL20,COL21,COL22,COL23,COL24,COL25,COL26,COL27,COL28,COL29,COL30,COL31,COL32,COL33,COL34,COL35,COL36,COL37,COL38,COL39 FROM %v`,
		tbl)

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
		if err := rows.Scan(&c.COL01, &c.COL02, &c.COL03, &c.COL04, &c.COL05, &c.COL06, &c.COL07, &c.COL10,
			&c.COL11, &c.COL12, &c.COL13, &c.COL14, &c.COL15, &c.COL16, &c.COL17, &c.COL18, &c.COL19, &c.COL20,
			&c.COL21, &c.COL22, &c.COL23, &c.COL24, &c.COL25, &c.COL26, &c.COL27, &c.COL28, &c.COL29, &c.COL30,
			&c.COL31, &c.COL32, &c.COL33, &c.COL34, &c.COL35, &c.COL36, &c.dbCOL37, &c.dbCOL38, &c.dbCOL39); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
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
