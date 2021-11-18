package fmdu

import "fmt"

const (
	ViewAreaFtth = "VIEW_AREA_FTTH"
)

type Area struct {
	Ftth string
}
type Araes []*Area

func (para *Parameter) GetArea() (as Araes, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT ftth FROM %v ", ViewAreaFtth)
	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	as = make(Araes, 0)
	for rows.Next() {
		a := new(Area)
		if err := rows.Scan(&a.Ftth); err != nil {
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
