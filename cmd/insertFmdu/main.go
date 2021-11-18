package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	//_ "github.com/lib/pq"

	"bitbucket.org/inadenoita/fmdu"
	"bitbucket.org/inadenoita/futalib"
)

const (
	myVersion = "1.0.0"
	buildDay  = "yyyy-mm-dd"
)

/*
const (
	TblMDU  = "TBL_MDU"
	TblFmdu = "TBL_FMDU"
	ViewMDU = "VIEW_MDU"
)

type mdu struct {
	BuilNo int
}
type mdus []*mdu

func connDB() (*sql.DB, error) {
	conn, err := sql.Open("postgres", getDBConStr())
	return conn, err
}

func getDBConStr() (connstr string) {
	connstr = fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable", "localhost", 5432, "fmdu", "inamap", "@Hanazono")
	return
}

func getMdus() (ms mdus, err error) {
	db, err := connDB()
	if err != nil {
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT COL30 FROM %v where ftth=1;", ViewMDU)

	fmt.Println(sql)

	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	ms = make(mdus, 0)
	for rows.Next() {
		m := new(mdu)
		//if err := rows.Scan(&b.ID, &b.Cell, &b.Akey, &b.Lng, &b.Lat, &b.BlockCell); err != nil {
		if err := rows.Scan(&m.BuilNo); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
		}
		ms = append(ms, m)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (ms mdus) insertFmdu() (idx uint) {
	idx = 0
	db, err := connDB()
	if err != nil {
		return
	}
	defer db.Close()

	//_, err = db.Query("TRUNCATE TBL_FMDU CASCADE;")

	for _, m := range ms {
		sql := fmt.Sprintf("SELECT buil_no FROM %v where buil_no=%v;", TblFmdu, m.BuilNo)
		//fmt.Println(sql)
		rows, err := db.Query(sql)
		if err != nil {
			//log.Fatal(err)
			return
		}
		defer rows.Close()
		tmpidx := 0
		for rows.Next() {
			m := new(mdu)
			if err := rows.Scan(&m.BuilNo); err != nil {
				fmt.Printf("値の取得に失敗しました。: %v\n", err)
			}
			tmpidx++
		}
		if 0 < tmpidx {
			continue
		}

		var fid uint64
		err = db.QueryRow("INSERT INTO TBL_FMDU(buil_no,ctim) VALUES($1, now()) RETURNING id", m.BuilNo).Scan(&fid)
		if err != nil {
			fmt.Println(err)
		}
		idx++
	}
	return
}
*/

func main() {
	start := time.Now()

	pg := os.Args[0]

	futalib.ShowPgInfo(pg, myVersion, buildDay)

	/*
		ms, err := getMdus()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		//spew.Dump(ms)
	*/

	//idx := ms.insertFmdu()
	para := new(fmdu.ParameterAPI)
	para.ParameterRDB = fmdu.InitParameterRDB("localhost", "5432", "fmdu", "inamap", "@Hanazono")

	idx, err := para.InsertFmdu()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("%v 件のレコードを登録しました。\n", idx)

	idx, err = para.InsertIntro()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Printf("%v 件のレコードを登録しました。\n", idx)

	end := time.Now()
	fmt.Printf("%f秒\n", (end.Sub(start)).Seconds())

	os.Exit(0)
}

// 引数取得
func getArgs() []string {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
%s [OPTIONS] inpath
Options
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var ()
	flag.Parse()

	return flag.Args()
}
