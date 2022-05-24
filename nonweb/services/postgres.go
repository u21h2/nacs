package services

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"nacs/common"
	nonweb_utils "nacs/nonweb/utils"
	"nacs/utils"
	"nacs/utils/logger"
	"strings"
	"time"
)

func PostgresScan(info nonweb_utils.HostInfo) (tmperr error) {
	if common.RunningInfo.NoBrute {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range common.Userdict["postgresql"] {
		for _, pass := range common.Passwords {
			pass = strings.Replace(pass, "{user}", string(user), -1)
			flag, err := PostgresConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("psql %v:%v %v %v %v", info.Host, info.Port, user, pass, err)
				if common.RunningInfo.BruteDebug {
					logger.Failed(errlog)
				}
				tmperr = err
				if utils.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common.Userdict["postgresql"])*len(common.Passwords)) * info.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func PostgresConn(info nonweb_utils.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Port, user, pass
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", Username, Password, Host, Port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(info.Timeout) * time.Second)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("Postgres:%v:%v:%v %v", Host, Port, Username, Password)
			logger.Success(result)
			flag = true
		}
	}
	return flag, err
}
