package utils

import (
	"nacs/common"
	"net/url"
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

// CompareTwoUrlEqual 比较是否相同
func CompareTwoUrlEqual(url1, url2 string) bool {
	if strings.EqualFold(url1, url2) {
		return true
	}
	parsedUrl1, _ := url.Parse(url1)
	schema1 := parsedUrl1.Scheme
	host1 := parsedUrl1.Host
	if strings.Contains(host1, ":") {
		host1 = strings.Split(host1, ":")[0]
	}
	port1 := parsedUrl1.Port()
	if port1 == "" {
		port1 = common.ServiceToPortString[schema1]
	}
	path1 := parsedUrl1.Path
	parsedUrl2, _ := url.Parse(url2)
	schema2 := parsedUrl2.Scheme
	host2 := parsedUrl2.Host
	if strings.Contains(host2, ":") {
		host2 = strings.Split(host2, ":")[0]
	}
	port2 := parsedUrl2.Port()
	if port2 == "" {
		port2 = common.ServiceToPortString[schema2]
	}
	path2 := parsedUrl2.Path
	//fmt.Println(schema1, host1, port1, path1)
	//fmt.Println(schema2, host2, port2, path2)
	if !strings.EqualFold(schema1, schema2) {
		return false
	} else if !strings.EqualFold(host1, host2) {
		return false
	} else if !strings.EqualFold(port1, port2) {
		return false
	} else if !strings.EqualFold(path1, path2) {
		if strings.EqualFold(path1[:len(path1)], path2) || strings.EqualFold(path1, path2[:len(path2)]) {
			return true
		} else {
			return false
		}
	}
	return true
}
