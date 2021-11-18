package fmdu

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

// Login ... ログイン処理サービス
func (dbpara *ParameterAPI) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var us Users
		para, err := dbpara.InitParameterUser(c)
		if err != nil {
			return para.retunClientHandler(c, us, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		us, err = para.Login()
		if err != nil {
			return para.retunClientHandler(c, us, 0, err, time.Since(start))
		}
		if para.Debug {
			spew.Dump(us)
		}

		length := len(us)
		return para.retunClientHandler(c, us, length, err, time.Since(start))
	}
}

// Logout ... ログアウト処理サービス
func (dbpara *ParameterAPI) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()
		var us Users
		para, err := dbpara.InitParameterUser(c)
		if err != nil {
			return para.retunClientHandler(c, us, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		return para.retunClientHandler(c, us, 0, err, time.Since(start))
	}
}

// GetUser ... ユーザー取得処理サービス
func (dbpara *ParameterAPI) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

		start := time.Now()

		var us Users
		para, err := dbpara.InitParameterUser(c)
		if err != nil {
			return para.retunClientHandler(c, us, 0, err, time.Since(start))
		}
		//dbpara.DB = para.DB
		if para.Debug {
			spew.Dump(para)
		}
		us, err = para.GetUser()
		if err != nil {
			return para.retunClientHandler(c, us, 0, err, time.Since(start))
		}
		if para.Debug {
			spew.Dump(us)
		}

		length := len(us)
		return para.retunClientHandler(c, us, length, err, time.Since(start))
	}
}
