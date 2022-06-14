package services

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"nacs/common"
	nonweb_utils "nacs/nonweb/utils"
	"nacs/utils"
	"nacs/utils/logger"
	"net"
	"strings"
	"time"
)

func SshScan(info nonweb_utils.HostInfo) (tmperr error) {
	if common.RunningInfo.NoBrute {
		return
	}
	starttime := time.Now().Unix()
	for _, user := range common.Userdict["ssh"] {
		for _, pass := range common.Passwords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := SshConn(info.Host, info.Port, user, pass)
			if flag == true && err == nil {
				return err
			} else {
				errlog := fmt.Sprintf("ssh %v:%v %v %v %v", info.Host, info.Port, user, pass, err)
				if common.RunningInfo.BruteDebug {
					logger.Failed(errlog)
				}
				tmperr = err
				if utils.CheckErrs(err) {
					return err
				}
				if time.Now().Unix()-starttime > (int64(len(common.Userdict["ssh"])*len(common.Passwords)) * int64(common.RunningInfo.BruteTimeout)) {
					return err
				}
			}
			if common.RunningInfo.SSHKey != "" {
				return err
			}
		}
	}
	return tmperr
}

func SshConn(Host, Port, user, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := Host, Port, user, pass
	Auth := []ssh.AuthMethod{}
	if common.RunningInfo.SSHKey != "" {
		pemBytes, err := ioutil.ReadFile(common.RunningInfo.SSHKey)
		if err != nil {
			return false, errors.New("read key failed" + err.Error())
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return false, errors.New("parse key failed" + err.Error())
		}
		Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		Auth = []ssh.AuthMethod{ssh.Password(Password)}
	}

	config := &ssh.ClientConfig{
		User:    Username,
		Auth:    Auth,
		Timeout: time.Duration(common.RunningInfo.BruteTimeout) * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", Host, Port), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		if err == nil {
			defer session.Close()
			flag = true
			var result string
			if common.RunningInfo.Command != "" {
				combo, _ := session.CombinedOutput(common.RunningInfo.Command)
				result = fmt.Sprintf("SSH:%v:%v %v %v `%v`: %v", Host, Port, Username, Password, common.RunningInfo.Command, strings.Replace(string(combo), "\n", "", 1))
				if common.RunningInfo.SSHKey != "" {
					result = fmt.Sprintf("SSH:%v:%v sshkey correct \n %v", Host, Port, string(combo))
				}
				logger.Success(result)
			} else {
				result = fmt.Sprintf("SSH:%v:%v:%v %v", Host, Port, Username, Password)
				if common.RunningInfo.SSHKey != "" {
					result = fmt.Sprintf("SSH:%v:%v sshkey correct", Host, Port)
				}
				logger.Success(result)
			}
		}
	}
	return flag, err

}
