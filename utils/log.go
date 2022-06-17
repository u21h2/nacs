package utils

import (
	"fmt"
	"nacs/common"
	"os"
	"strings"
	"sync"
	"time"
)

var Num int64
var End int64
var Results = make(chan *string)
var Start = true
var LogSucTime int64
var LogErrTime int64
var WaitTime int64
var Silent bool
var LogWG sync.WaitGroup

func init() {
	go SaveLog()
}

func LogSuccess(result string) {
	LogWG.Add(1)
	LogSucTime = time.Now().Unix()
	Results <- &result
}

func SaveLog() {
	for result := range Results {
		if Silent == false || strings.Contains(*result, "[+]") || strings.Contains(*result, "[*]") {
			fmt.Println(*result)
		}
		if !common.RunningInfo.NoSave {
			WriteFile(*result, common.RunningInfo.OutputFileName)
		}
		LogWG.Done()
	}
}

func WriteFile(result string, filename string) {
	var text = []byte(result + "\n")
	fl, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Open %s error, %v\n", filename, err)
		return
	}
	_, err = fl.Write(text)
	fl.Close()
	if err != nil {
		fmt.Printf("Write %s error, %v\n", filename, err)
	}
}

func LogError(errinfo interface{}) {
	if WaitTime == 0 {
		fmt.Println(fmt.Sprintf(" %v/%v %v", End, Num, errinfo))
	} else if (time.Now().Unix()-LogSucTime) > WaitTime && (time.Now().Unix()-LogErrTime) > WaitTime {
		fmt.Println(fmt.Sprintf(" %v/%v %v", End, Num, errinfo))
		LogErrTime = time.Now().Unix()
	}
}
