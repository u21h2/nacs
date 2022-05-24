package services

import (
	"fmt"
	"nacs/common"
	nonweb_utils "nacs/nonweb/utils"
	"nacs/utils"
	"nacs/utils/logger"
	"strings"
	"time"
)

func MemcachedScan(info nonweb_utils.HostInfo) (err error) {
	realhost := fmt.Sprintf("%s:%v", info.Host, info.Port)
	client, err := utils.WrapperTcpWithTimeout("tcp", realhost, time.Duration(info.Timeout)*time.Second)
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	if err == nil {
		err = client.SetDeadline(time.Now().Add(time.Duration(info.Timeout) * time.Second))
		if err == nil {
			_, err = client.Write([]byte("stats\n")) //Set the key randomly to prevent the key on the server from being overwritten
			if err == nil {
				rev := make([]byte, 1024)
				n, err := client.Read(rev)
				if err == nil {
					if strings.Contains(string(rev[:n]), "STAT") {
						result := fmt.Sprintf("Memcached %s unauthorized", realhost)
						logger.Success(result)
					}
				} else {
					errlog := fmt.Sprintf(" Memcached %v:%v %v", info.Host, info.Port, err)
					if common.RunningInfo.BruteDebug {
						logger.Failed(errlog)
					}
				}
			}
		}
	}
	return err
}
