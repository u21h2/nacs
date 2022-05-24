package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func GetPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path1, _ := filepath.Abs(file)
	filename := filepath.Dir(path1)

	var path string
	//if strings.Contains(filename, "/") {
	//	tmp := strings.Split(filename, `/`)
	//	tmp[len(tmp)-1] = ``
	//	path = strings.Join(tmp, `/`)
	//} else if strings.Contains(filename, `\`) {
	//	tmp := strings.Split(filename, `\`)
	//	tmp[len(tmp)-1] = ``
	//	path = strings.Join(tmp, `\`)
	//}
	if strings.Contains(filename, "/") {
		path = filename + "/"
	} else if strings.Contains(filename, `\`) {
		path = filename + `\`
	}
	return path
}

func RemoveDuplicate(old []string) []string {
	result := []string{}
	temp := map[string]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicateInt(old []int) []int {
	result := []int{}
	temp := map[int]struct{}{}
	for _, item := range old {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func CheckErrs(err error) bool {
	if err == nil {
		return false
	}
	errs := []string{
		"closed by the remote host", "too many connections",
		"i/o timeout", "EOF", "A connection attempt failed",
		"established connection failed", "connection attempt failed",
		"Unable to read", "is not allowed to connect to this",
		"no pg_hba.conf entry",
		"No connection could be made",
		"invalid packet size",
		"bad connection",
	}
	for _, key := range errs {
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower(key)) {
			return true
		}
	}
	return false
}
