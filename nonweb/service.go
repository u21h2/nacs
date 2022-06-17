package nonweb

import (
	"nacs/common"
	"nacs/nonweb/services"
	"nacs/nonweb/utils"
	"nacs/utils/logger"
	"strconv"
)

func Service() {
	logger.Info("Start to process nonweb services")
	for _, discoverResult := range common.DiscoverResults {
		var info utils.HostInfo
		info.Host = discoverResult["host"].(string)
		info.Port = strconv.Itoa(discoverResult["port"].(int))
		info.Domain = ""
		info.Timeout = int64(common.RunningInfo.BruteTimeout)

		switch discoverResult["protocol"].(string) {
		case "ssh":
			logger.Info("[protocol] ssh " + info.Host)
			services.SshScan(info)
			continue
		case "redis":
			logger.Info("[protocol] redis " + info.Host + ":" + info.Port)
			services.RedisScan(info)
			continue
		case "smb":
			logger.Info("[protocol] smb-CVE-2020-0796 " + info.Host + ":" + info.Port)
			services.SmbGhostScan(info)
			logger.Info("[protocol] smb-17010 " + info.Host + ":" + info.Port)
			services.MS17010Scan(info)
			logger.Info("[protocol] smb-brute " + info.Host + ":" + info.Port)
			services.SmbScan(info)
			continue
		case "rdp":
			logger.Info("[protocol] rdp " + info.Host + ":" + info.Port)
			services.RdpScan(info)
			continue
		case "oracle":
			logger.Info("[protocol] oracle " + info.Host + ":" + info.Port)
			services.OracleScan(info)
			continue
		case "mysql":
			logger.Info("[protocol] mysql " + info.Host + ":" + info.Port)
			services.MysqlScan(info)
			continue
		case "mssql":
			logger.Info("[protocol] mssql " + info.Host + ":" + info.Port)
			services.MssqlScan(info)
			continue
		case "ftp":
			logger.Info("[protocol] ftp " + info.Host + ":" + info.Port)
			//services.FtpScan(info)
			continue
		case "DceRpc":
			logger.Info("[protocol] netbios " + info.Host + ":" + info.Port)
			services.Findnet(info)
			continue
			//case "memcache":
			//	logger.Info("[protocol] memcache " + info.Host + ":" + info.Port)
			//	services.MemcachedScan(info)
			//	continue
			//case "mongo":
			//	logger.Info("[protocol] mongo " + info.Host + ":" + info.Port)
			//	services.MongodbScan(info)
			//	continue
			//case "psql":
			//	logger.Info("[protocol] psql " + info.Host + ":" + info.Port)
			//	services.PostgresScan(info)
			//	continue
		}
		switch info.Port {
		case "22":
			logger.Info("[port] ssh " + info.Host)
			services.SshScan(info)
			continue
		case "6379":
			logger.Info("[port] redis " + info.Host + ":" + info.Port)
			services.RedisScan(info)
			continue
		case "445":
			logger.Info("[port] smb-CVE-2020-0796 " + info.Host + ":" + info.Port)
			services.SmbGhostScan(info)
			logger.Info("[port] smb-17010 " + info.Host + ":" + info.Port)
			services.MS17010Scan(info)
			logger.Info("[port] smb-brute " + info.Host + ":" + info.Port)
			services.SmbScan(info)
			continue
		case "3389":
			logger.Info("[port] rdp " + info.Host + ":" + info.Port)
			services.RdpScan(info)
			continue
		case "5432":
			logger.Info("[port] psql " + info.Host + ":" + info.Port)
			services.PostgresScan(info)
			continue
		case "1521":
			logger.Info("[port] oracle " + info.Host + ":" + info.Port)
			services.OracleScan(info)
			continue
		case "139":
			logger.Info("[port] netbios " + info.Host + ":" + info.Port)
			services.NetBIOS(info)
			continue
		case "3306":
			logger.Info("[port] mysql " + info.Host + ":" + info.Port)
			services.MysqlScan(info)
			continue
		case "1433":
			logger.Info("[port] mssql " + info.Host + ":" + info.Port)
			services.MssqlScan(info)
			continue
		case "11211":
			logger.Info("[port] memcache " + info.Host + ":" + info.Port)
			services.MemcachedScan(info)
			continue
		case "27017":
			logger.Info("[port] mongo " + info.Host + ":" + info.Port)
			services.MongodbScan(info)
			continue
		case "21":
			logger.Info("[port] ftp " + info.Host + ":" + info.Port)
			//services.FtpScan(info)
			continue
		case "135":
			logger.Info("[port] netbios " + info.Host + ":" + info.Port)
			services.Findnet(info)
			continue
		case "9000":
			logger.Info("[port] fcgi " + info.Host + ":" + info.Port)
			services.FcgiScan(info)
			continue
		}
	}
}
