package services

import (
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	"nacs/common"
	nonweb_utils "nacs/nonweb/utils"
	"nacs/utils"
	"nacs/utils/logger"
	"strings"
	"time"
)

func OracleScan(info nonweb_utils.HostInfo) (tmperr error) {
	if common.RunningInfo.NoBrute {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range common.Userdict["oracle"] {
		for _, pass := range common.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := OracleConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("oracle %v:%v %v %v %v", info.Host, info.Port, user, pass, err)
				if common.RunningInfo.BruteDebug {
					logger.Failed(errlog)
				}
				tmperr = err
				if utils.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common.Userdict["oracle"])*len(common.Passwords)) * info.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func OracleConn(info nonweb_utils.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Port, user, pass
	dataSourceName := fmt.Sprintf("oracle://%s:%s@%s:%s/orcl", Username, Password, Host, Port)
	db, err := sql.Open("oracle", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(info.Timeout) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(info.Timeout) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("oracle:%v:%v:%v %v", Host, Port, Username, Password)
			logger.Success(result)
			flag = true
		}
	}
	return flag, err
}
