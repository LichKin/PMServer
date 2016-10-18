package main

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL env
const (
	DB_NAME = "PM_CLOUD"
	DB_USER = "root"
	DB_PW   = "qwe123"
	DB_HOST = "localhost"
	DB_PORT = "3306"
)

func CheckNil(key string, value string) string {
	if key == "" {
		return value
	}
	return key
}

var Custom_MysqlManager = new(MysqlStorageManager)

type MysqlStorageManager struct {
	db  *sql.DB
	dsn string // 用于连接的字符串
}

// GetMysqlDSN 方法负责生成Mysql连接字符串，传入Mysql的基本配置后会返回标准的连接字符串作为结果。
// 这个结果可以直接作为NewMysqlStorage方法的参数
func GetMysqlDSN(userName string, password string, protocal string, dbName string, params map[string]string) string {
	dsn := ""

	dsn += userName

	if password != "" {
		dsn += ":" + password
	}

	if protocal != "" {
		dsn += "@tcp(" + protocal + ")"
	}

	if dbName != "" {
		dsn += "/" + dbName
	}

	if len(params) > 0 {
		dsn += "?"
		paramList := []string{}
		for k, v := range params {
			paramList = append(paramList, k+"="+v+"&")
		}
		dsn += join(paramList, "&")
	}

	Log.Info("storage_mysql_config:", dsn)
	return dsn
}

func (this *MysqlStorageManager) SetDsn(dsn string) {
	this.dsn = dsn
}

func (this *MysqlStorageManager) Connect() (*sql.DB, error) {
	this.db, _ = sql.Open("mysql", GetMysqlDSN(
		CheckNil(os.Getenv("DB_USER"), DB_USER),
		CheckNil(os.Getenv("DB_PW"), DB_PW),
		CheckNil(os.Getenv("DB_HOST"), DB_HOST)+":"+CheckNil(os.Getenv("DB_PORT"), DB_PORT),
		CheckNil(os.Getenv("DB_NAME"), DB_NAME),
		map[string]string{}))

	this.db.SetMaxOpenConns(2000)
	this.db.SetMaxIdleConns(500)

	return this.db, nil
}

// Query 接受一个sql语句和其中包含的变量，返回查找到的结果
func (this *MysqlStorageManager) Query(sqlString string, args ...interface{}) (*sql.Rows, error) {
	//return this.db.Query(sqlString, args...)
	rows, err := this.db.Query(sqlString, args...)
	defer rows.Close()
	return rows, err
}

// Prepare 准备sql语句
func (this *MysqlStorageManager) Prepare(sqlString string) (*sql.Stmt, error) {
	return this.db.Prepare(sqlString)
}

// Exec 执行update，insert，delete操作，接收sql语句和保存的变量
func (this *MysqlStorageManager) Exec(sqlString string, args ...interface{}) (sql.Result, error) {
	return this.db.Exec(sqlString, args...)
}

func join(a []string, sep string) string {
	if len(a) == 0 {
		return ""
	}
	if len(a) == 1 {
		return a[0]
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i])
	}

	b := make([]byte, n)
	bp := copy(b, a[0])
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s)
	}
	return string(b)
}
