package fmdu

import (
	"database/sql"
	"fmt"
	"strconv"
)

const (
	ViewCadSb = "VIEW_CAD_SB"
)

type lnglat struct {
	longitude sql.NullFloat64
	latitude  sql.NullFloat64
}
type lnglatMap map[uint64]lnglat

func (para *Parameter) getLngLatMap() (llmap lnglatMap, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	stat := fmt.Sprintf(`SELECT fac_no,longitude,latitude FROM %v `, ViewCadSb)

	if para.Debug {
		fmt.Println(stat)
	}
	rows, err := db.Query(stat)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	llmap = make(lnglatMap)
	for rows.Next() {
		var facno string
		var lng sql.NullFloat64
		var lat sql.NullFloat64
		if err := rows.Scan(&facno, &lng, &lat); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
			continue
		}
		var builno uint64
		tmp, err := strconv.Atoi(facno)
		if err != nil {
			fmt.Println(err)
			continue
		}
		builno = uint64(tmp)
		ll := new(lnglat)
		ll.longitude = lng
		ll.latitude = lat
		llmap[builno] = *ll
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (llmap lnglatMap) getLngLat(builno uint64) (lng sql.NullFloat64, lat sql.NullFloat64) {
	ll, ok := llmap[builno]
	if ok {
		lng = ll.longitude
		lat = ll.latitude
		return
	}
	return
}
