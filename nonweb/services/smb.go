package services

import (
	"errors"
	"fmt"
	"github.com/stacktitan/smb/smb"
	"nacs/common"
	nonweb_utils "nacs/nonweb/utils"
	"nacs/utils"
	"nacs/utils/logger"
	"strings"
	"time"
)

func SmbScan(info nonweb_utils.HostInfo) (tmperr error) {
	if common.RunningInfo.NoBrute {
		return nil
	}
	starttime := time.Now().Unix()
	for _, user := range common.Userdict["smb"] {
		for _, pass := range common.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := doWithTimeOut(info, user, pass)
			if flag == true && err == nil {
				var result string
				if info.Domain != "" {
					result = fmt.Sprintf("[+] SMB:%v:%v:%v\\%v %v", info.Host, info.Port, info.Domain, user, pass)
				} else {
					result = fmt.Sprintf("[+] SMB:%v:%v:%v %v", info.Host, info.Port, user, pass)
				}
				logger.Success(result)
				return err
			} else {
				errlog := fmt.Sprintf("smb %v:%v %v %v %v", info.Host, 445, user, pass, err)
				errlog = strings.Replace(errlog, "\n", "", -1)
				if common.RunningInfo.BruteDebug {
					logger.Failed(errlog)
				}
				tmperr = err
				if utils.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common.Userdict["smb"])*len(common.Passwords)) * info.Timeout) {
					return err
				}
			}
		}
	}
	return tmperr
}

func SmblConn(info nonweb_utils.HostInfo, user string, pass string, signal chan struct{}) (flag bool, err error) {
	flag = false
	Host, Username, Password := info.Host, user, pass
	options := smb.Options{
		Host:        Host,
		Port:        445,
		User:        Username,
		Password:    Password,
		Domain:      info.Domain,
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			flag = true
		}
	}
	signal <- struct{}{}
	return flag, err
}

func doWithTimeOut(info nonweb_utils.HostInfo, user string, pass string) (flag bool, err error) {
	signal := make(chan struct{})
	go func() {
		flag, err = SmblConn(info, user, pass, signal)
	}()
	select {
	case <-signal:
		return flag, err
	case <-time.After(time.Duration(info.Timeout) * time.Second):
		return false, errors.New("time out")
	}
}
