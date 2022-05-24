package check

import (
	"fmt"
	"nacs/utils/logger"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/panjf2000/ants"
	"nacs/web/poc/internal/common/errors"
	common_structs "nacs/web/poc/pkg/common/structs"
	nuclei_structs "nacs/web/poc/pkg/nuclei/structs"
	xray_structs "nacs/web/poc/pkg/xray/structs"
)

var (
	EmptyLinks = []string{}

	Ticker  *time.Ticker
	Pool    *ants.PoolWithFunc
	Verbose bool

	WaitGroup sync.WaitGroup

	OutputChannel chan common_structs.Result

	ResultPool = sync.Pool{
		New: func() interface{} {
			return new(common_structs.PocResult)
		},
	}
)

// 初始化协程池
func InitCheck(threads, rate int, verbose bool) {
	var err error

	rateLimit := time.Second / time.Duration(rate)
	Ticker = time.NewTicker(rateLimit)
	Pool, err = ants.NewPoolWithFunc(threads, check)
	if err != nil {
		logger.Error("Initialize goroutine pool error: " + err.Error())
	}

	Verbose = verbose
}

// 将任务放入协程池
func Start(targets []string, xrayPocMap map[string]xray_structs.Poc, nucleiPocMap map[string]nuclei_structs.Poc, outputChannel chan common_structs.Result) {
	// 设置outputChannel
	OutputChannel = outputChannel

	for _, target := range targets {
		for _, poc := range xrayPocMap {
			WaitGroup.Add(1)
			Pool.Invoke(&xray_structs.Task{
				Target: target,
				Poc:    poc,
			})
		}
		for _, poc := range nucleiPocMap {
			WaitGroup.Add(1)
			Pool.Invoke(&nuclei_structs.Task{
				Target: target,
				Poc:    poc,
			})
		}
	}
}

// 等待协程池
func Wait() {
	WaitGroup.Wait()
}

// 释放协程池
func End() {
	Pool.Release()
}

// 核心代码，poc检测
func check(taskInterface interface{}) {
	var (
		oRequest *http.Request = nil

		isVul   bool
		err     error
		pocName string
	)

	defer WaitGroup.Done()
	<-Ticker.C

	switch taskInterface.(type) {
	case *xray_structs.Task:
		task, ok := taskInterface.(*xray_structs.Task)
		if !ok {
			wrappedErr := errors.Newf(errors.ConvertInterfaceError, "Can't convert task interface: %#v", err)
			//utils.ErrorP(wrappedErr)
			logger.DebugError(wrappedErr)
			return
		}
		target, poc := task.Target, task.Poc

		pocName = poc.Name
		if poc.Transport != "tcp" && poc.Transport != "udp" {
			oRequest, _ = http.NewRequest("GET", target, nil)
		}

		isVul, err = executeXrayPoc(oRequest, target, &poc)
		if err != nil {
			//utils.ErrorP(err)
			logger.DebugError(err)
			return
		}

		pocResult := ResultPool.Get().(*common_structs.PocResult)
		pocResult.Str = fmt.Sprintf("%s (%s)", target, pocName)
		pocResult.Success = isVul
		pocResult.URL = target
		pocResult.PocName = poc.Name
		pocResult.PocLink = poc.Detail.Links
		pocResult.PocAuthor = poc.Detail.Author
		pocResult.PocDescription = poc.Detail.Description

		OutputChannel <- pocResult

	case *nuclei_structs.Task:
		var (
			desc    string
			author  string
			authors []string
		)

		task, ok := taskInterface.(*nuclei_structs.Task)
		if !ok {
			wrappedErr := errors.Newf(errors.ConvertInterfaceError, "Can't convert task interface: %#v", err)
			//utils.ErrorP(wrappedErr)
			logger.DebugError(wrappedErr)
			return
		}
		target, poc := task.Target, task.Poc
		authors, ok = poc.Info.Authors.Value.([]string)
		if !ok {
			author = "Unknown"
		} else {
			author = strings.Join(authors, ", ")
		}

		results, isVul, err := executeNucleiPoc(target, &poc)
		if err != nil {
			//utils.ErrorP(err)
			logger.DebugError(err)
			return
		}

		for _, r := range results {
			if r.ExtractorName != "" {
				desc = r.TemplateID + ":" + r.ExtractorName
			} else if r.MatcherName != "" {
				desc = r.TemplateID + ":" + r.MatcherName
			}

			pocResult := ResultPool.Get().(*common_structs.PocResult)
			pocResult.Str = fmt.Sprintf("%s (%s) ", r.Matched, r.TemplateID)
			pocResult.Success = isVul
			pocResult.URL = r.Matched
			pocResult.PocName = r.TemplateID
			pocResult.PocLink = EmptyLinks
			pocResult.PocAuthor = author
			pocResult.PocDescription = desc

			OutputChannel <- pocResult
		}
	}

}

func PutPocResult(result *common_structs.PocResult) {
	ResultPool.Put(result)
}
