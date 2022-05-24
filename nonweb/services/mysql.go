package services

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"nacs/common"
	nonweb_utils "nacs/nonweb/utils"
	"nacs/utils"
	"nacs/utils/logger"
	"strings"
	"time"
)

func MysqlScan(info nonweb_utils.HostInfo) (tmperr error) {
	if common.RunningInfo.NoBrute {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range common.Userdict["mysql"] {
		for _, pass := range common.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := MysqlConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("mysql %v:%v %v %v %v", info.Host, info.Port, user, pass, err)
				if common.RunningInfo.BruteDebug {
					logger.Failed(errlog)
				}

				tmperr = err
				if utils.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common.Userdict["mysql"])*len(common.Passwords)) * info.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func MysqlConn(info nonweb_utils.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Port, user, pass
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/mysql?charset=utf8&timeout=%v", Username, Password, Host, Port, time.Duration(info.Timeout)*time.Second)
	db, err := sql.Open("mysql", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(info.Timeout) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(info.Timeout) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("mysql:%v:%v:%v %v", Host, Port, Username, Password)
			logger.Success(result)
			flag = true
		}
	}
	return flag, err
}
