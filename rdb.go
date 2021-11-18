package fmdu

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	// postgresql Driver

	"github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
)

// DB種別定義
const (
	DBTypePsql = "postgres"
	//DBTypeSQLitesql = "sqlite3"

	DBCmd     = "psql"
	DBDumpCmd = "pg_dump"

	extDump   = ".dump"
	dumpFname = "fmdu" + extDump
)

// RDBデフォルト値
const (
	DefaultDB     = "demo"
	DefaultHost   = "localhost"
	DefaultUser   = "inamap"
	DefaultPasswd = "@Hanazono"
	DefaultPort   = "5432"
)

const (
	MaxOpenConns = 10000
)

const (
	NULL = "NULL"
)

// ParameterRDB ... RDB接続用パラメータ
type ParameterRDB struct {
	DBHost   string
	DBPort   string
	DBName   string
	DBUser   string
	DBPasswd string
	*parameterImport
}

// InitParameterRDB ... RDB接続用パラメータ設定
func InitParameterRDB(host, port, db, user, passwd string) (para *ParameterRDB) {
	para = new(ParameterRDB)
	if len(host) == 0 {
		para.DBHost = DefaultHost
	} else {
		para.DBHost = host
	}
	if len(port) == 0 {
		para.DBPort = DefaultPort
	} else {
		para.DBPort = port
	}
	if len(db) == 0 {
		para.DBName = DefaultDB
	} else {
		para.DBName = db
	}
	if len(user) == 0 {
		para.DBUser = DefaultUser
	} else {
		para.DBUser = user
	}
	if len(passwd) == 0 {
		para.DBPasswd = DefaultPasswd
	} else {
		para.DBPasswd = passwd
	}

	return
}

// ConnDB ... DBへの接続
func (para *ParameterRDB) ConnDB() (*sql.DB, error) {
	conn, err := sql.Open(DBTypePsql, para.GetDBConStr())
	return conn, err
}

// GetDBConStr ... DB接続文字列を得る
func (para *ParameterRDB) GetDBConStr() (connstr string) {
	connstr = fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable", para.DBHost, para.DBPort, para.DBName, para.DBUser, para.DBPasswd)
	return
}

// GetDBConStrDefault ... DB接続文字列を得る※DB名以外はデフォルト値を使用
func GetDBConStrDefault(db string) (connstr string) {
	connstr = fmt.Sprintf("host=%v port=%v dbname=%v user=%v password=%v sslmode=disable", DefaultHost, DefaultPort, DefaultDB, DefaultUser, DefaultPasswd)
	return
}

const (
	defaultDelim    = ","
	defaultQuote    = ""
	defaultEncoding = "UTF8"
	optionHeader    = "header"
)

type parameterImport struct {
	csv       string
	tbl       string
	delimiter string
	//quote     string
	column   string
	isheader bool
	isdelete bool
}

func initParameterImport() (para *parameterImport) {
	para = new(parameterImport)
	para.isdelete = true
	para.isheader = false
	para.delimiter = defaultDelim
	return
}

func (para *ParameterRDB) execCopySQL() (cnt int, err error) {
	if len(para.csv) == 0 {
		err = fmt.Errorf("CSV Pathが未指定")
		return
	}
	if _, err = os.Stat(para.csv); err != nil {
		return
	}
	if len(para.tbl) == 0 {
		err = fmt.Errorf("Table名が未指定")
		return
	}
	if len(para.delimiter) == 0 {
		para.delimiter = defaultDelim
	}
	//if len(para.quote) == 0 {
	//	para.quote = defaultQuote
	//}
	header := ""
	if para.isheader {
		header = optionHeader
	}

	cmd := DBCmd
	var cmdargs []string

	// 既存レコード削除
	if para.isdelete {
		cmdargs = []string{"-U", fmt.Sprintf("%v", para.DBUser), "-d", fmt.Sprintf("%v", para.DBName), "-c", fmt.Sprintf(`TRUNCATE %v RESTART IDENTITY;`, para.tbl)}
		fmt.Println((spew.Sdump(cmd, cmdargs)))
		_, err = exec.Command(cmd, cmdargs...).CombinedOutput()
		if err != nil {
			return
		}
	}

	// CSVインポート
	cmdargs = []string{"-U", fmt.Sprintf("%v", para.DBUser), "-d", fmt.Sprintf("%v", para.DBName),
		"-c", fmt.Sprintf(`\copy %v %v from '%v' delimiter '%v' csv %v ;`, para.tbl, para.column, para.csv, para.delimiter, header)}
	fmt.Println((spew.Sdump(cmd, cmdargs)))
	_, err = exec.Command(cmd, cmdargs...).CombinedOutput()
	if err != nil {
		return
	}

	cnt, err = para.countTable()

	return
}

func (para *ParameterRDB) execSQL(cmd string) (cnt int, err error) {
	if len(para.tbl) == 0 {
		err = fmt.Errorf("Table名が未指定")
		return
	}
	if len(cmd) == 0 {
		err = fmt.Errorf("SQL文が未指定")
		return
	}

	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(MaxOpenConns)

	fmt.Println(cmd)
	_, err = db.Query(cmd)
	if err != nil {
		return
	}
	cnt, err = para.countTable()
	return
}

func (para *ParameterRDB) countTable() (cnt int, err error) {
	if len(para.tbl) == 0 {
		err = fmt.Errorf("Table名が未指定")
		return
	}
	db, err := para.ConnDB()
	if err != nil {
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT count(*) FROM %v ;", para.tbl)

	fmt.Println(sql)

	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	cnt = 0
	for rows.Next() {
		//m := new(Mdu)
		//if err := rows.Scan(&b.ID, &b.Cell, &b.Akey, &b.Lng, &b.Lat, &b.BlockCell); err != nil {
		if err := rows.Scan(&cnt); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
		}
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

// 全角・半角混在検索関数使用用の文字列を返す
func getNoCaseTextSQLStr(arg string) (str string) {
	str = fmt.Sprintf(" '%%' || GET_NOCASE_TEXT('%v') || '%%' ", arg)
	return
}

// RDBバックアップコマンド実行
func (para *ParameterBackup) execExportRDB() (err error) {
	if len(para.backupPathToday) == 0 {
		err = fmt.Errorf("Backup Pathが未指定")
		return
	}
	if len(para.DBName) == 0 {
		err = fmt.Errorf("DB名が未指定")
		return
	}
	if len(para.DBUser) == 0 {
		err = fmt.Errorf("DBUser名が未指定")
		return
	}

	if _, err = os.Stat(para.backupPathToday); err != nil {
		return
	}

	dumppath := filepath.Join(para.backupPathToday, dumpFname)
	//cmd := fmt.Sprintf("%v %v > %v", DBDumpCmd, para.DBName, dumppath)
	cmd := DBDumpCmd

	// DUMP生成
	cmdargs := []string{"-U", fmt.Sprintf("%v", para.DBUser), "-d", fmt.Sprintf("%v", para.DBName), "-f", fmt.Sprintf("%v", dumppath)}
	fmt.Println((spew.Sdump(cmd, cmdargs)))
	_, err = exec.Command(cmd, cmdargs...).CombinedOutput()
	if err != nil {
		return
	}

	return
}
