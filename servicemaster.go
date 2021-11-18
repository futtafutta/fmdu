package fmdu

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

func (dbpara *ParameterAPI) GetFtype() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var fs Ftypes
		para, err := dbpara.InitParameter(c)
		if err != nil {
			return para.retunClientHandler(c, fs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		fs, err = para.GetFtype()
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

func (dbpara *ParameterAPI) GetCompany() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var cs Companys
		para, err := dbpara.InitParameter(c)
		if err != nil {
			return para.retunClientHandler(c, cs, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		cs, err = para.GetCompany()
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

func (dbpara *ParameterAPI) GetIntroStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var ss Statuss
		//para, err := dbpara.InitParameter(c)
		para, err := dbpara.InitParameterIntroStatus(c)
		if err != nil {
			return para.retunClientHandler(c, ss, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		ss, err = para.GetIntroStatus()
		if err != nil {
			return para.retunClientHandler(c, ss, 0, err, time.Since(start))
		}
		length := len(ss)
		if para.Debug {
			//spew.Dump(fs)
			fmt.Printf("*** Info *** RecordCount : %v \n", length)
		}
		return para.retunClientHandler(c, ss, length, err, time.Since(start))
	}
}

func (dbpara *ParameterAPI) GetArea() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var as Araes
		para, err := dbpara.InitParameter(c)
		if err != nil {
			return para.retunClientHandler(c, as, 0, err, time.Since(start))
		}
		if para.Debug {
			spew.Dump(para)
		}
		as, err = para.GetArea()
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
