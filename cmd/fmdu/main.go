// ************************************************************************************************************************
// ProgramName:	fmdu
// Project:		FMDU
// Description:	芯線情報管理サービス For SCN
// Usage:		fmdu  (configDirectoryPath)
//				$ sudo service fmdu start
// Env:			linux/amd64
//				Linuxでビルドすること
//				ビルドした実行ファイルをDaemon化して登録すること
// Histrory:
//				2021/07/28	Ver:1.0.0 Go1.16.5   初リリース版
//				2021/08/04	Ver:1.0.1 Go1.16.5   "log"サービス追加実装
//				2021/10/28	Ver:1.1.0 Go1.17.2   10/22のレビュー結果を反映
//				2021/11/05	Ver:1.1.1 Go1.17.2   importCust()内imppara.column内の改行を除去
//				2021/11/05	Ver:1.1.2 Go1.17.2   Backupサービス追加（RDB&添付データ）
//				2021/11/06	Ver:1.1.3 Go1.17.2   [intro]のHFC数（うち放送数）=>（うち放送のみ数）となるように処理修正
//				2021/11/06	Ver:1.1.4 Go1.17.2   [densyou]の電障のみ？の判定結果が逆になっている不具合の修正
//				2021/11/06	Ver:1.1.5 Go1.17.2   各リスト検索項目に「FTTHエリア」を追加+リスト表示
//												 RDBに「TBL_AREA」追加し、各VIEWにFTTHエリア名を追加
//				2021/11/11	Ver:1.1.6 Go1.17.2   マスタ取得関連のAPI処理をlogに出力しないように設定
//				2021/11/12	Ver:1.1.7 Go1.17.2   パラメータ文字列中に半角スペースを含むケースの対策としてencode<=>decodeするように修正
//				2021/11/12	Ver:1.1.8 Go1.17.2   [intro]のフィルタ状況に「進捗状況」DirectStatusを追加
//				2021/11/12	Ver:1.1.9  Go1.17.2  [introstatus]API用にParameterIntroStatus構造体を追加
//				2021/11/16	Ver:1.1.10 Go1.17.2  [TBL_INTRO]に文字列型 列「EXPLANATION_PERSON」（出席担当）追加による拡張
//
// ************************************************************************************************************************

package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/inadenoita/fmdu"

	"bitbucket.org/inadenoita/futalib"

	"github.com/davecgh/go-spew/spew"
)

const (
	version  = "1.1.9"
	buildday = "2021-11-12"
)

func main() {
	futalib.ShowPgInfo(os.Args[0], version, buildday)

	args := getArgs()

	if len(args) < 1 {
		fmt.Println("*** Error *** 引数にConfig Directoryパスが無い")
		flag.Usage()
		os.Exit(-1)
	}
	path := args[0]
	tmls := futalib.GetExtFiles(&[]string{}, path, path, ".tml")
	spew.Dump(tmls)

	wg, semaphore := futalib.MakeSemaphore(len(*tmls))
	for _, tml := range *tmls {
		wg.Add(1)
		go func(tml string) {
			defer wg.Done()
			semaphore <- 1
			fmdu.StartServer(tml)
			<-semaphore
		}(tml)
	}
	wg.Wait()

}

// 引数取得
func getArgs() (args []string) {
	// -hオプション用文言
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `
Usage of %s:
%s [OPTIONS] (Config Directory Path)
Options
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	var ()
	flag.Parse()

	args = flag.Args()

	return
}
