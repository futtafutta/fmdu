package fmdu

import (
	"database/sql"
	"fmt"
)

const (
	TblMDU  = "TBL_MDU"
	ViewMDU = "VIEW_MDU"
)

// -- Modf 2021/10/25 登録NO,世帯数もセットするように拡張
type Mdu struct {
	BuilNo    uint64
	TourokuNo uint64
	BuilCnt   sql.NullInt64
}
type Mdus []*Mdu

// Import用にMDUデータ取得
func (para ParameterAPI) getMdus(ftthflags ...int8) (ms Mdus, err error) {
	var ftthflag int8 = -1
	if 0 < len(ftthflags) {
		ftthflag = ftthflags[0]
	}
	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()

	sql := ""

	// -- Modf 2021/10/25 登録NO,世帯数もセットするように拡張
	if ftthflag == -1 {
		sql = fmt.Sprintf("SELECT COL30,COL24,COl03 FROM %v ;", ViewMDU)
	} else {
		sql = fmt.Sprintf("SELECT COL30,COL24,COL03 FROM %v where ftth=%v;", ViewMDU, ftthflag)
	}

	fmt.Println(sql)

	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	ms = make(Mdus, 0)
	for rows.Next() {
		m := new(Mdu)
		//if err := rows.Scan(&b.ID, &b.Cell, &b.Akey, &b.Lng, &b.Lat, &b.BlockCell); err != nil {
		if err := rows.Scan(&m.BuilNo, &m.TourokuNo, &m.BuilCnt); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
		}
		ms = append(ms, m)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}
