package main

import (
	//	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	//	"time"

	"github.com/gorilla/mux"
)

const (
	INSERT_SQL = "INSERT INTO TB_APP(APP_NAME,APP_COUNT,APP_PASSWORD,UNAME,LASTUPDATE)VALUES(?,?,?,?,?)"
	DELETE_SQL = "DELETE FROM TB_APP WHERE UNAME=?"
	QUERY_SQL  = "SELECT APP_NAME,APP_COUNT,APP_PASSWORD,LASTUPDATE FROM TB_APP WHERE UNAME=?"
)

type APPINFO struct {
	APPNAME     string `json:"appname"`
	APPCOUNT    string `json:"appcount"`
	APPPASSWORD string `json:"apppassword"`
	LASTUPDATE  string `json:"updatetime"`
}

type APPINFOS struct {
	Uname string    `json:"username"`
	INFOS []APPINFO `json:"infos"`
}

func TransferHandler(response http.ResponseWriter, request *http.Request) {

	Log.Info("TransferHandler...")
	defer request.Body.Close()
	rbody, err := ioutil.ReadAll(request.Body)
	var infos APPINFOS
	if err != nil {
		processError(500, "获取body失败", response)
		return
	} else {
		Log.Info(string(rbody))
		err = json.Unmarshal(rbody, &infos)
		if err != nil {
			processError(500, "解析body失败", response)
			return
		} else if len(infos.INFOS) > 0 {
			//删除与该用户相关的旧的记录
			Custom_MysqlManager.db.Exec(DELETE_SQL, infos.Uname)
			//插入
			for _, appinfo := range infos.INFOS {

				_, err = Custom_MysqlManager.db.Exec(INSERT_SQL,
					appinfo.APPNAME, appinfo.APPCOUNT, appinfo.APPPASSWORD, infos.Uname, appinfo.LASTUPDATE)
				if err != nil {
					Log.Error("Insert failed at app: ", appinfo.APPNAME)
					Log.Error(err)
				}
			}
		} else {
			processError(500, "length is 0! ", response)
			return
		}

	}

	response.Write(convertParam2byte(200, "transfer success"))
}

func PullHandler(response http.ResponseWriter, request *http.Request) {
	Log.Info("PullHandler...")
	vars := mux.Vars(request)
	uname := vars["uname"]
	appinfos := new(APPINFOS)
	var name, count, pwd, updatetime string
	flag := 0
	if uname != "" {
		appinfos.Uname = uname
		var infoarr []APPINFO
		rows, err := Custom_MysqlManager.db.Query(QUERY_SQL, uname)
		if err != nil {
			processError(500, err.Error(), response)
			return
		} else {
			for rows.Next() {
				flag++
				err = rows.Scan(&name, &count, &pwd, &updatetime)
				if err != nil {
					Log.Error("读取数据失败。。")
				} else {
					appinfo := APPINFO{name, count, pwd, updatetime}
					infoarr = append(infoarr, appinfo)
				}
			}
			if flag == 0 {
				processError(400, "No rows", response)
				return
			}
			appinfos.INFOS = infoarr
			results, err := json.Marshal(appinfos)
			Log.Info(string(results))
			if err != nil {
				processError(500, err.Error(), response)
				return
			}
			response.Write(convertParam2byte(200, string(results)))

		}
	} else {
		processError(500, "参数为空！", response)
	}

}

func processError(code int, msg string, res http.ResponseWriter) {
	Log.Error(msg)
	res.Write(convertParam2byte(code, msg))
}

func convertParam2byte(code int, message string) []byte {
	data, _ := json.Marshal(map[string]interface{}{"code": code, "msg": message})
	return data
}
