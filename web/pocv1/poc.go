package pocv1

import (
	"fmt"
	"nacs/common"
	"nacs/utils/logger"
	"nacs/web/pocv1/lib"
	"nacs/web/pocv1/poc_struct"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var AllPocs []*lib.Poc

func Poc() {
	if common.RunningInfo.NoPoc {
		return
	}
	logger.Info("Start to send pocs to web services (xray type)")
	LoadPoc()
	logger.Info(logger.LightGreen("Load ") +
		logger.White(strconv.Itoa(len(AllPocs))) +
		logger.LightGreen(" xray poc(s) "))
	InitPocInfo(&poc_struct.PocInfo)
	lib.Inithttp(poc_struct.PocInfo)
	for _, discoverResult := range common.DiscoverResults {
		if discoverResult["protocol"].(string) == "http" || discoverResult["protocol"].(string) == "https" {
			url := discoverResult["uri"].(string)
			RunPoc(url)
		}
	}

}

func RunPoc(url string) {
	var pocInfo = poc_struct.PocInfo
	pocInfo.Target = url // strings.Join(buf[:3], "/")
	Execute(pocInfo)
}

func Execute(PocInfo poc_struct.PocInfoStruct) {
	req, err := http.NewRequest("GET", PocInfo.Target, nil)
	if err != nil {
		errlog := fmt.Sprintf("[-] webpocinit %v %v", PocInfo.Target, err)
		logger.Error(errlog)
		return
	}
	req.Header.Set("User-agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1468.0 Safari/537.36")
	if PocInfo.Cookie != "" {
		req.Header.Set("Cookie", PocInfo.Cookie)
	}
	lib.CheckMultiPoc(req, AllPocs, PocInfo.Num)
}

func InitPocInfo(pocInfo *poc_struct.PocInfoStruct) {
	pocInfo.Num = common.RunningInfo.PocRate
	pocInfo.Timeout = int64(common.RunningInfo.PocTimeout)
	pocInfo.ApiKey = common.RunningInfo.CeyeApi
	pocInfo.CeyeDomain = common.RunningInfo.CeyeDomain

}

// InitPoc
func LoadPoc() {
	err := filepath.Walk(common.RunningInfo.FscanPocPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil {
				return err
			}
			if !info.IsDir() {
				if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
					poc, _ := lib.LoadPocbyPath(path)
					if poc != nil {
						AllPocs = append(AllPocs, poc)
					}
				}
			}
			return nil
		})
	if err != nil {
		fmt.Printf("[-] init xray poc error: %v", err)
	}

}
