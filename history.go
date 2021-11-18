package fmdu

import (
	"fmt"
)

const (
	//ViewCplCust = "view_fmdu_cpl_cust"
	TblCustHistory = "TBL_CUSTOMER_HISTORY"
)

type CplCust struct {
	CplNo  string
	BuilNo uint64
	CustNo uint64
}

//type CplCusts []*CplCust
type CplCustMap map[string][]interface{}

func (para *Parameter) getCplHistory() (cmap CplCustMap, err error) {
	cmap = make(CplCustMap)
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf(`SELECT COL19,COL01,COL23 FROM %v`, TblCustHistory)
	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		c := new(CplCust)
		if err := rows.Scan(&c.CplNo, &c.BuilNo, &c.CustNo); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		key := fmt.Sprintf("%v-%v-%v", c.BuilNo, c.CplNo, c.CustNo)
		cmap[key] = nil
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (para *Parameter) updateCplHistory(cmap CplCustMap) (err error) {

	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf(`SELECT COL19,COL01,COL23 FROM %v`, ViewCplCust)

	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		c := new(Customer)
		if err := rows.Scan(&c.COL19, &c.COL01, &c.COL23); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		key := fmt.Sprintf("%v-%v-%v", c.COL01, c.COL19.String, c.COL23.Int64)
		_, ok := cmap[key]
		if ok {
			continue
		}
		// 履歴テーブルに存在しない顧客番号であれば顧客レコードをInsertする
		err = para.insertHistory(uint64(c.COL23.Int64))
		if err != nil {
			fmt.Println(err)
		}

	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (para *Parameter) insertHistory(custno uint64) (err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	var id uint
	err = db.QueryRow(`INSERT INTO TBL_CUSTOMER_HISTORY (COL01,COL02,COL03,COL04,COL05,COL06,COL07,COL10,
		COL11,COL12,COL13,COL14,COL15,COL16,COL17,COL18,COL19,COL20,COL21,COL22,COL23,COL24,COL25,
		COL26,COL27,COL28,COL29,COL30,COL31,COL32,COL33,COL34,COL35,COL36,COL37,COL38,COL39
		) select 
			COL01,COL02,COL03,COL04,COL05,COL06,COL07,COL10,
			COL11,COL12,COL13,COL14,COL15,COL16,COL17,COL18,COL19,COL20,COL21,COL22,COL23,COL24,COL25,
			COL26,COL27,COL28,COL29,COL30,COL31,COL32,COL33,COL34,COL35,COL36,COL37,COL38,COL39
		from VIEW_FMDU_CPL_CUST where COL23 = $1 RETURNING id`, custno).Scan(&id)

	return
}

func (para *Parameter) UpdateCplHistory() (err error) {
	cmap, err := para.getCplHistory()
	if err != nil {
		return
	}

	//spew.Dump(cmap)

	err = para.updateCplHistory(cmap)

	return
}
