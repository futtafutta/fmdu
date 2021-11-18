package fmdu

import "fmt"

const (
	TblFtype = "TBL_FTYPE"
)

type Ftype struct {
	ID   uint
	Name string
}
type Ftypes []*Ftype

func (para *Parameter) GetFtype() (fs Ftypes, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,name FROM %v", TblFtype)

	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	fs = make(Ftypes, 0)
	for rows.Next() {
		f := new(Ftype)
		if err := rows.Scan(&f.ID, &f.Name); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		fs = append(fs, f)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}
