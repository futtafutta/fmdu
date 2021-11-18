package fmdu

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	defaultPort      = "8888"
	defaultDBHost    = "localhost"
	defaultDBPort    = "5432"
	defaultlog       = "user.log"
	defaultCrawllog  = "crawl.log"
	defaultBackuplog = "backup.log"
)

// ServerConfig ... API Server設定情報構造体
type ServerConfig struct {
	Config *Config
}

// Config ... TOML記述ファイルの設定項目格納構造体
type Config struct {
	Port          string `toml:"port"`
	Kei           uint8  `toml:"kei"`
	Log           string `toml:"log"`
	DBHost        string `toml:"dbhost"`
	DBPort        string `toml:"dbport"`
	DBName        string `toml:"dbname"`
	DBUser        string `toml:"dbuser"`
	DBPasswd      string `toml:"dbpasswd"`
	AppPath       string `toml:"apppath"`
	SmsPath       string `toml:"sms_path"`
	CsvPath       string `toml:"csv_path"`
	CsvMdu        string `toml:"csv_mdu"`
	CsvCust       string `toml:"csv_cust"`
	CrawlLogPath  string `toml:"crawllog"`
	BackupPath    string `toml:"backup"`
	BackupLogPath string `toml:"backuplog"`
}

func StartServer(path string) (err error) {
	confs, err := getConfig(path)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("*** Error *** Configの解析に失敗[%v]\n", path)
		return
	}
	conf := confs.Config
	spew.Dump(conf)
	if len(conf.DBName) == 0 {
		fmt.Println("*** Error *** ConfigにDBHostの指定が無い")
		return
	}
	if len(conf.Port) == 0 {
		conf.Port = defaultPort
	}
	port, err := strconv.ParseUint(conf.Port, 10, 64)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("*** Error *** ConfigのPort番号が不正[%v]\n", conf.Port)
		return
	}
	if conf.Kei <= 0 {
		fmt.Printf("*** Error *** ConfigのKEIが不正[%v]\n", conf.Kei)
		return
	}
	if len(conf.AppPath) == 0 {
		fmt.Println("*** Error *** ConfigにAppPathの指定が無い")
		return
	}
	if len((conf.SmsPath)) == 0 {
		fmt.Println("*** Error *** Configにsms_pathの指定が無い")
		return
	}
	if len((conf.CsvPath)) == 0 {
		fmt.Println("*** Error *** Configにcsv_pathの指定が無い")
		return
	}
	if len((conf.CsvMdu)) == 0 {
		fmt.Println("*** Error *** Configにcsv_mduの指定が無い")
		return
	}
	if len((conf.CsvCust)) == 0 {
		fmt.Println("*** Error *** Configにcsv_custの指定が無い")
		return
	}

	para := new(ParameterAPI)
	dbpara := InitParameterRDB(conf.DBHost, conf.DBPort, conf.DBName, conf.DBUser, conf.DBPasswd)
	para.ParameterRDB = dbpara
	//para.JSONPath = conf.JSONPath
	para.Kei = conf.Kei
	para.AppPath = conf.AppPath
	//para.CrawlExe = conf.CrawlExe
	//para.CrawlConf = conf.CrawlConf
	log := defaultlog
	if 0 < len(conf.Log) {
		log = conf.Log
	}
	para.LogPath = log
	para.SmsPath = conf.SmsPath
	para.CsvPath = conf.CsvPath
	para.CsvMdu = conf.CsvMdu
	para.CsvCust = conf.CsvCust

	log = defaultCrawllog
	if 0 < len(conf.CrawlLogPath) {
		log = conf.CrawlLogPath
	}
	para.CrawlLogPath = log

	backuppath := filepath.Join(para.AppPath, defaultBackupDir)
	if 0 < len(conf.BackupPath) {
		backuppath = conf.BackupPath
	}
	para.BackupPath = backuppath

	log = defaultBackuplog
	if 0 < len(conf.BackupLogPath) {
		log = conf.BackupLogPath
	}
	para.BackupLogPath = log

	// Echoのインスタンス作る
	e := echo.New()

	// 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// ルーティング
	//e.GET("/lonlat/:ll", handler.MainPage())

	apppath := AppPath

	//e.Static("/", "f-mdu")

	e.GET(apppath+APIMethodLogin, para.Login())
	e.GET(apppath+APIMethodLogout, para.Logout())

	e.GET(apppath+APIMethodFtype, para.GetFtype())
	e.GET(apppath+APIMethodCompany, para.GetCompany())
	e.GET(apppath+APIMethodIntroStatus, para.GetIntroStatus())
	e.GET(apppath+APIMethodArea, para.GetArea())

	e.GET(apppath+APIMethodFmdu, para.GetFmdu())
	e.POST(apppath+APIMethodFmdu, para.PostFmdu())

	e.GET(apppath+APIMethodCpl, para.GetCpl())
	e.POST(apppath+APIMethodCpl, para.PostCpl())
	e.PUT(apppath+APIMethodCpl, para.PutCpl())
	e.DELETE(apppath+APIMethodCpl, para.DeleteCpl())

	e.GET(apppath+APIMethodAttached, para.GetAttached())
	e.POST(apppath+APIMethodAttached, para.PostAttached())
	e.PUT(apppath+APIMethodAttached, para.PutAttached())
	e.DELETE(apppath+APIMethodAttached, para.DeleteAttached())

	e.GET(apppath+APIMethodIntro, para.GetIntro())
	e.PUT(apppath+APIMethodIntro, para.PutIntro())

	e.GET(apppath+APIMethodDensyou, para.GetDensyou())
	e.PUT(apppath+APIMethodDensyou, para.PutDensyou())

	e.GET(apppath+APIMethodCustomer, para.GetCustomer())
	e.PUT(apppath+APIMethodCustomer, para.PutCustomer())

	e.GET(apppath+APIMethodImport, para.Import())
	e.GET(apppath+APIMethodImportLog, para.GetImportLog())

	e.GET(apppath+APIMethodBackup, para.Backup())

	e.GET(apppath+APIMethodLog, para.GetLog())

	// サーバー起動
	e.Start(fmt.Sprintf(":%v", port)) //ポート番号指定してね
	//e.Logger.Fatal(e.StartAutoTLS(":1323"))

	return
}

func getConfig(tml string) (conf *ServerConfig, err error) {
	fmt.Println("config")
	_, err = toml.DecodeFile(tml, &conf)
	if err != nil {
		fmt.Println(err)
	}

	return
}
