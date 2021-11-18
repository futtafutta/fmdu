package fmdu

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"bitbucket.org/inadenoita/futalib"
	"github.com/labstack/echo/v4"
)

const (
	defaultBackupDir = "backup"
	extArchive       = ".tar.gz"
)

type ParameterBackup struct {
	*Parameter
	IsBackupRDB      bool
	IsBackupAttached bool
	//BackupPath       string
	backupPathToday string
	today           string
}

func (dbpara *ParameterAPI) InitParameterBackup(c echo.Context) (para *ParameterBackup, err error) {
	para = new(ParameterBackup)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}

	if len(para.AppPath) == 0 {
		err = fmt.Errorf("no AppPath")
		return
	}

	//para.BackupPath = filepath.Join(para.AppPath, defaultBackupDir)
	para.IsBackupRDB = true
	para.IsBackupAttached = true

	if len(para.IP) == 0 {
		para.Parameter.IP = SystemIP
	}
	if len(para.UserID) == 0 {
		para.Parameter.UserID = SystemUser
	}

	return
}

// ImportSmsData ... SMS-CSV import処理
func (para *ParameterBackup) Backup() (err error) {
	if !para.IsBackupRDB && !para.IsBackupAttached {
		return
	}
	err = para.initBackup()
	if err != nil {
		return
	}

	if para.IsBackupRDB {
		err = para.backupRDB()
		if err != nil {
			return
		}
	}

	if para.IsBackupAttached {
		err = para.backupAttached()
		if err != nil {
			return
		}
	}

	return
}

// RDBのBackup
func (para *ParameterBackup) backupRDB() (err error) {
	err = para.execExportRDB()
	return
}

// 添付データディレクトリのBackup
func (para *ParameterBackup) backupAttached() (err error) {
	srcpath := filepath.Join(para.AppPath, DataPath)
	dstpath := filepath.Join(para.backupPathToday, DataPath)
	cmd := fmt.Sprintf("cp -r %v %v", srcpath, dstpath)
	_, err = futalib.RunCmdStrShell(cmd)
	if err != nil {
		return
	}

	// 現在のディレクトリを保存してからcd
	prevDir, err := filepath.Abs(".")
	if err != nil {
		return
	}

	if err = os.Chdir(para.BackupPath); err != nil {
		return
	}

	//compresspath := fmt.Sprintf("%v%v", para.today, extArchive)
	// 圧縮
	cmd = fmt.Sprintf("tar zcvf %v%v %v", para.today, extArchive, para.today)
	_, err = futalib.RunCmdStrShell(cmd)
	if err != nil {
		return
	}
	// 元に戻す
	if err = os.Chdir(prevDir); err != nil {
		return
	}

	ret := futalib.DelFilesAll(para.backupPathToday)
	if !ret {
		err = fmt.Errorf("*** Error *** failuer remove dir [%v]", para.backupPathToday)
		return
	}

	return
}

// Backup初期化
func (para *ParameterBackup) initBackup() (err error) {
	now := time.Now()
	para.today = now.Format(LayoutDay)
	para.backupPathToday = filepath.Join(para.BackupPath, para.today)
	_, err = os.Stat(para.backupPathToday)
	if err != nil {
		ret := futalib.MkdirNewdir(para.backupPathToday)
		if !ret {
			err = fmt.Errorf("*** Error *** failuer newdir [%v]", para.backupPathToday)
			return
		}
	}
	err = nil
	return
}

// Backupログの保存
func (para *ParameterBackup) writeLog(succ int, start, end time.Time, dur time.Duration) (err error) {

	line := fmt.Sprintf("%v|%v|%v|%v|%v|%v", start.Format(layoutTimestamp), end.Format(layoutTimestamp), dur, para.IP, para.UserID, succ)
	fmt.Println(line)

	file, err := os.OpenFile(para.BackupLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	fmt.Fprintln(file, line) //書き込み

	return
}
