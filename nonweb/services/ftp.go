package services

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"nacs/common"
	nonweb_utils "nacs/nonweb/utils"
	"nacs/utils"
	"nacs/utils/logger"
	"strings"
	"time"
)

func FtpScan(info nonweb_utils.HostInfo) (tmperr error) {
	if common.RunningInfo.NoBrute {
		return
	}
	starttime := time.Now().Unix()
	flag, err := FtpConn(info, "anonymous", "")
	if flag == true && err == nil {
		return err
	} else {
		errlog := fmt.Sprintf("ftp://%v:%v %v %v", info.Host, info.Port, "anonymous", err)
		if common.RunningInfo.BruteDebug {
			logger.Failed(errlog)
		}
		tmperr = err
		if utils.CheckErrs(err) {
			return err
		}
	}

	for _, user := range common.Userdict["ftp"] {
		for _, pass := range common.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := FtpConn(info, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("[-] ftp://%v:%v %v %v %v", info.Host, info.Port, user, pass, err)
				if common.RunningInfo.BruteDebug {
					logger.Failed(errlog)
				}
				tmperr = err
				if utils.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common.Userdict["ftp"])*len(common.Passwords)) * info.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func FtpConn(info nonweb_utils.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Port, user, pass
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", Host, Port), time.Duration(info.Timeout)*time.Second)
	if err == nil {
		err = conn.Login(Username, Password)
		if err == nil {
			flag = true
			result := fmt.Sprintf("ftp://%v:%v:%v %v", Host, Port, Username, Password)
			dirs, err := conn.List("")
			//defer conn.Logout()
			if err == nil {
				if len(dirs) > 0 {
					for i := 0; i < len(dirs); i++ {
						if len(dirs[i].Name) > 50 {
							result += "\n   [->]" + dirs[i].Name[:50]
						} else {
							result += "\n   [->]" + dirs[i].Name
						}
						if i == 5 {
							break
						}
					}
				}
			}
			logger.Success(result)
		}
	}
	return flag, err
}
