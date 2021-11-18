package fmdu

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/labstack/echo/v4"
)

func (dbpara *ParameterAPI) GetFmdu() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var fs Fmdus
		para, err := dbpara.InitParameterFmdu(c)
		if err != nil {
			return para.retunClientHandler(c, fs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		fs, err = para.GetFmdu()
		if err != nil {
			return para.retunClientHandler(c, fs, 0, err, time.Since(start))
		}
		length := len(fs)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, fs, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) PostFmdu() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var fs Fmdus
		para, err := dbpara.InitParameterFmdu(c)
		if err != nil {
			return para.retunClientHandler(c, fs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		fs, err = para.PostFmdu()
		if err != nil {
			return para.retunClientHandler(c, fs, 0, err, time.Since(start))
		}
		length := len(fs)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, fs, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) GetCustomer() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var cs Customers
		para, err := dbpara.InitParameterCustomer(c)
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		cs, err = para.GetCustomer()
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		length := len(cs)
		if para.Debug {
			//spew.Dump(cs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, cs, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) PutCustomer() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var cs Customers
		para, err := dbpara.InitParameterCustomer(c)
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		cs, err = para.UpdateCustomer()
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		length := len(cs)
		if para.Debug {
			//spew.Dump(cs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, cs, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) GetCpl() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var cs Cpls
		para, err := dbpara.InitParameterCpl(c)
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		cs, err = para.GetCpl()
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		length := len(cs)
		if para.Debug {
			//spew.Dump(cs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, cs, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) PostCpl() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var cs Cpls
		para, err := dbpara.InitParameterCpl(c)
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		cs, err = para.InsertCpl()
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		length := len(cs)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, cs, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) PutCpl() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var cs Cpls
		para, err := dbpara.InitParameterCpl(c)
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		cs, err = para.UpdateCpl()
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		length := len(cs)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, cs, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) DeleteCpl() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var cs Cpls
		para, err := dbpara.InitParameterCpl(c)
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		cs, err = para.DeleteCpl()
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		length := len(cs)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, cs, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) GetAttached() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var as Attacheds
		para, err := dbpara.InitParameterAttached(c)
		if err != nil {
			return para.retunClientHandler(c, as, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		as, err = para.GetAttached()
		if err != nil {
			return para.retunClientHandler(c, as, 0, err, time.Since(start))
		}
		length := len(as)
		if para.Debug {
			//spew.Dump(cs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, as, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) PostAttached() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var as Attacheds
		para, err := dbpara.InitParameterAttached(c)
		if err != nil {
			return para.retunClientHandler(c, as, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		as, err = para.InsertAttached()
		if err != nil {
			return para.retunClientHandler(c, as, 0, err, time.Since(start))
		}
		length := len(as)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, as, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) PutAttached() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var as Attacheds
		para, err := dbpara.InitParameterAttached(c)
		if err != nil {
			return para.retunClientHandler(c, as, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		as, err = para.UpdateAttached()
		if err != nil {
			return para.retunClientHandler(c, as, 0, err, time.Since(start))
		}
		length := len(as)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, as, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) DeleteAttached() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var as Attacheds
		para, err := dbpara.InitParameterAttached(c)
		if err != nil {
			return para.retunClientHandler(c, as, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		as, err = para.DeleteAttached()
		if err != nil {
			return para.retunClientHandler(c, as, 0, err, time.Since(start))
		}
		length := len(as)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, as, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) GetIntro() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var is Intros
		para, err := dbpara.InitParameterIntrto(c)
		if err != nil {
			return para.retunClientHandler(c, is, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		is, err = para.GetIntro()
		if err != nil {
			return para.retunClientHandler(c, is, 0, err, time.Since(start))
		}
		length := len(is)
		if para.Debug {
			//spew.Dump(cs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, is, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) PutIntro() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var is Intros
		para, err := dbpara.InitParameterIntrto(c)
		if err != nil {
			return para.retunClientHandler(c, is, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		is, err = para.UpdateIntro()
		if err != nil {
			return para.retunClientHandler(c, is, 0, err, time.Since(start))
		}
		length := len(is)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, is, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) GetDensyou() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var ds Densyous
		para, err := dbpara.InitParameterDensyou(c)
		if err != nil {
			return para.retunClientHandler(c, ds, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		ds, err = para.GetDensyou()
		if err != nil {
			return para.retunClientHandler(c, ds, 0, err, time.Since(start))
		}
		length := len(ds)
		if para.Debug {
			//spew.Dump(cs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, ds, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) PutDensyou() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
		start := time.Now()

		var ds Densyous
		para, err := dbpara.InitParameterDensyou(c)
		if err != nil {
			return para.retunClientHandler(c, ds, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		ds, err = para.UpdateDensyou()
		if err != nil {
			return para.retunClientHandler(c, ds, 0, err, time.Since(start))
		}
		length := len(ds)
		if para.Debug {
			//spew.Dump(cs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, ds, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) Import() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
		succ := 0

		start := time.Now()

		//var cs Companys
		para, err := dbpara.InitParameterImport(c)
		if err != nil {

			return para.retunClientHandler(c, nil, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		length, err := para.ImportSmsData()
		end := time.Now()
		dur := end.Sub(start)
		if err != nil {
			para.writeLog(succ, start, end, dur)
			return para.retunClientHandler(c, nil, 0, err, dur)
		}
		succ = 1
		//length := 0
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		para.writeLog(succ, start, end, dur)
		return para.retunClientHandler(c, 0, length, err, dur)
	}
}

func (dbpara *ParameterAPI) GetImportLog() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var is ImportLogs
		para, err := dbpara.InitParameterImport(c)
		if err != nil {

			return para.retunClientHandler(c, is, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		is, err = para.GetImportLog()
		if err != nil {
			return para.retunClientHandler(c, is, 0, err, time.Since(start))
		}
		length := len(is)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, is, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) GetLog() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var logs Logs
		para, err := dbpara.InitParameterLog(c)
		if err != nil {

			return para.retunClientHandler(c, logs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		logs, err = para.GetLog()
		if err != nil {
			return para.retunClientHandler(c, logs, 0, err, time.Since(start))
		}
		length := len(logs)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, logs, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) Backup() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
		succ := 0

		start := time.Now()

		//var cs Companys
		para, err := dbpara.InitParameterBackup(c)
		if err != nil {

			return para.retunClientHandler(c, nil, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		err = para.Backup()
		end := time.Now()

		dur := end.Sub(start)
		if err != nil {
			para.writeLog(succ, start, end, dur)
			return para.retunClientHandler(c, nil, 0, err, dur)
		}
		succ = 1
		if para.Debug {
			fmt.Printf("*** Info *** BackupStatus : %v \n", succ)
		}
		para.writeLog(succ, start, end, dur)
		return para.retunClientHandler(c, 0, 0, err, dur)
	}
}
