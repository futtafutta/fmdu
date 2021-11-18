package fmdu

import "fmt"

const (
	TblConstCampany = "TBL_CONST_COMPANY"
)

type Company struct {
	ID   uint
	Name string
}
type Companys []*Company

func (para *Parameter) GetCompany() (cs Companys, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,name FROM %v", TblConstCampany)

	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	cs = make(Companys, 0)
	for rows.Next() {
		c := new(Company)
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		cs = append(cs, c)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}
