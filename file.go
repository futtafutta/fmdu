package fmdu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/inadenoita/futalib"
)

func getDstFilePath(dstpath, ext string, builno uint64) (path string) {
	t := time.Now()
	tstr := t.Format("20060102_150405.000000000")
	tstr = strings.ReplaceAll(tstr, ".", "_")

	builstr := fmt.Sprintf("%v", builno)
	path = filepath.Join(dstpath, builstr)

	filename := fmt.Sprintf("%v%v", tstr, ext)
	path = filepath.Join(path, filename)

	return
}

// phpのアップロードtempパスからデータ格納パスへとコピーする
func copyFile(srcpath, dstpath string) (err error) {
	d, _ := filepath.Split(dstpath)
	_, err = os.Stat((d))
	if err != nil {
		ok := futalib.MkdirNewdir(d)
		if !ok {
			err = fmt.Errorf("*** Error *** Directory作成失敗[%v]", d)
			return
		}
	}
	err = futalib.CopyFile(srcpath, dstpath)

	return
}

func deleteFile(path string) (err error) {
	_, err = os.Stat(path)
	if err != nil {
		return
	}
	err = futalib.DelFile(path)

	return
}
