package fmdu

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	// TblUser ... RDBテーブル名
	TblUser = "TBL_USER_MANAGE"
)

// User ... ユーザー情報
type User struct {
	ID         uint
	UserName   string
	Password   string
	Authority  uint
	AuthorityS string
	DateCreate time.Time
	DateUpdate time.Time
}

// Users ... ユーザー情報配列
type Users []*User

// ParameterUser ... ユーザー（ログイン）用パラメータ構造体
type ParameterUser struct {
	*Parameter
	ID          uint
	UserName    string
	Password    string
	OldPassword string
	Authority   uint
}

// InitParameterUser ... ユーザー情報パラメータ初期化
func (dbpara *ParameterAPI) InitParameterUser(c echo.Context) (para *ParameterUser, err error) {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	uname := c.QueryParam("uname")
	passwd := c.QueryParam("passwd")
	oldpasswd := c.QueryParam("oldpasswd")
	auth, _ := strconv.Atoi(c.QueryParam("auth"))

	para = new(ParameterUser)
	para.Parameter, err = dbpara.InitParameter(c)
	if err != nil {
		return
	}
	para.ID = uint(id)
	para.UserName = uname
	para.Password = passwd
	para.OldPassword = oldpasswd
	para.Authority = uint(auth)

	return
}

// Login ... ログイン（ユーザー名+パスワードチェック）
func (para *ParameterUser) Login() (us Users, err error) {
	if (len(para.UserName) == 0) || (len(para.Password) == 0) {
		err = fmt.Errorf("ユーザー名またはパスワードが未入力です。")
		return
	}
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,user_id,password,authority FROM %v where user_id='%v' order by id;", TblUser, para.UserName)

	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	us = make(Users, 0)
	for rows.Next() {
		u := new(User)
		if err := rows.Scan(&u.ID, &u.UserName, &u.Password, &u.Authority); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
		}
		u.getAuthName()
		us = append(us, u)
	}
	if err = rows.Err(); err != nil {
		return
	}
	if len(us) == 0 {
		err = fmt.Errorf("ユーザー名が登録されていません。")
		return
	}
	u := us[0]
	hash := getHashMD5(para.Password)
	if hash == u.Password {
		if para.Debug {
			fmt.Println(hash)
		}
	} else {
		err = fmt.Errorf("パスワードが一致しません。")
		us = make(Users, 0)
	}

	return
}

// MD5暗号化文字列を取得する
func getHashMD5(str string) (hash string) {
	b := []byte(str)
	md5 := md5.Sum(b)
	hash = hex.EncodeToString(md5[:])
	return
}

// ユーザー権限整数値
const (
	// UserAuthAdmin ... 管理者ユーザー
	UserAuthAdmin = 99
	// UserAuthPower ... パワーユーザー
	UserAuthPower = 50
	// UserAuthNormal ... 一般ユーザー
	UserAuthNormal = 1
)

func (u *User) getAuthName() {
	var auth = ""
	switch u.Authority {
	case UserAuthAdmin:
		auth = "管理者"
	case UserAuthPower:
		auth = "パワーユーザー"
	default:
		auth = "一般ユーザー"
	}
	u.AuthorityS = auth
}

// GetUser ... ユーザーリストを返す
func (para *ParameterUser) GetUser() (us Users, err error) {
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,user_id,password,authority,date_create,date_update FROM %v order by id;", TblUser)
	if 0 < para.ID {
		sql = fmt.Sprintf("SELECT id,user_id,password,authority,date_create,date_update FROM %v where id=%v order by id;", TblUser, para.ID)
	} else if 0 < len(para.UserName) {
		sql = fmt.Sprintf("SELECT id,user_id,password,authority,date_create,date_update FROM %v where user_id='%v' order by id;", TblUser, para.UserName)
	} else if 0 < para.Authority {
		sql = fmt.Sprintf("SELECT id,user_id,password,authority,date_create,date_update FROM %v where authority=%v order by id;", TblUser, para.Authority)
	}

	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	us = make(Users, 0)
	for rows.Next() {
		u := new(User)
		if err := rows.Scan(&u.ID, &u.UserName, &u.Password, &u.Authority, &u.DateCreate, &u.DateUpdate); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
		}
		u.getAuthName()
		us = append(us, u)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

// PostUser ... ユーザーを登録する
func (para *ParameterUser) PostUser() (us Users, err error) {
	if len(para.UserName) < 1 {
		err = fmt.Errorf("*** Error *** ユーザー名がない")
		return
	}
	if len(para.Password) < 1 {
		err = fmt.Errorf("*** Error *** パスワードがない")
		return
	}
	if para.Authority == 0 {
		para.Authority = 1
	}

	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,user_id,password,authority,date_create,date_update FROM %v where user_id='%v' order by id;", TblUser, para.UserName)
	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	isexist := false
	for rows.Next() {
		u := new(User)
		if err := rows.Scan(&u.ID); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
		}
		isexist = true
	}
	if err = rows.Err(); err != nil {
		return
	}
	if isexist {
		err = fmt.Errorf("*** Error *** [%v] は、既に登録済みユーザーです。", para.UserName)
		return
	}

	var uid uint
	//var now = time.Now()
	err = db.QueryRow("INSERT INTO TBL_USER_MANAGE (user_id,password,authority,date_create,date_update) VALUES($1,$2,$3,now(),now()) RETURNING id",
		para.UserName, getHashMD5(para.Password), para.Authority).Scan(&uid)
	if err != nil {
		return
	}
	fmt.Println(uid)
	para.ID = uid
	us, err = para.GetUser()

	return
}

// update前に更新データをチェックする
func (para *ParameterUser) checkUpdateValue(dbpasswd string, dbauth uint) (err error) {
	authchkOk := false
	passchkOk := false
	if para.ID < 1 {
		err = fmt.Errorf("*** Error *** IDがない")
		return
	}
	if 0 < para.Authority {
		if para.Authority != dbauth {
			authchkOk = true
		}
	}
	if 0 < len(para.Password) {
		if len(para.OldPassword) < 1 {
			err = fmt.Errorf("*** Error *** 現在のパスワードがない")
			return
		}
		if len(para.Password) < 1 {
			err = fmt.Errorf("*** Error *** 新パスワードがない")
			return
		}
		oldpwd := getHashMD5(para.OldPassword)
		newpwd := getHashMD5(para.Password)
		if oldpwd != dbpasswd {
			err = fmt.Errorf("*** Error *** 現在のパスワードが一致しません")
			return
		}
		if (oldpwd == newpwd) || (newpwd == dbpasswd) {
			err = fmt.Errorf("*** Error *** 新旧パスワードが同じです")
		} else {
			passchkOk = true
		}
	}

	if passchkOk || authchkOk {
		err = nil
	}
	return
}

// PutUser ... ユーザーを更新する
func (para *ParameterUser) PutUser() (us Users, err error) {
	if para.Authority == 0 {
		para.Authority = 1
	}

	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT id,user_id,password,authority,date_create,date_update FROM %v where id='%v' order by id;", TblUser, para.ID)
	if para.Debug {
		fmt.Println(sql)
	}
	rows, err := db.Query(sql)
	if err != nil {
		//log.Fatal(err)
		return
	}
	defer rows.Close()

	var u User
	flag := false
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.UserName, &u.Password, &u.Authority, &u.DateCreate, &u.DateUpdate); err != nil {
			fmt.Printf("値の取得に失敗しました。: %v\n", err)
		}
		flag = true
	}
	if !flag {
		err = fmt.Errorf("*** Error *** レコードが存在しません。")
		return
	}

	if err = para.checkUpdateValue(u.Password, u.Authority); err != nil {
		return
	}

	if err = rows.Err(); err != nil {
		return
	}
	sql = fmt.Sprintf("UPDATE TBL_USER_MANAGE SET password = '%v', authority = %v , date_update = now()  WHERE id = %v", getHashMD5(para.Password), para.Authority, para.ID)
	//var now = time.Now()
	_, err = db.Query(sql)
	if err != nil {
		return
	}
	us, err = para.GetUser()

	return
}

// DeleteUser ... ユーザーを削除する
func (para *ParameterUser) DeleteUser() (us Users, err error) {
	if para.ID < 1 {
		err = fmt.Errorf("*** Error *** IDがない")
		return
	}
	db, err := para.ConnDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sql := fmt.Sprintf("DELETE FROM TBL_USER_MANAGE WHERE id = %v", para.ID)
	//var now = time.Now()
	_, err = db.Query(sql)
	if err != nil {
		return
	}

	return
}
