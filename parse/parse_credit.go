package parse

import (
	"nacs/common"
	"strings"
)

func ParseUser(InputInfo *common.InputInfoStruct) {
	if InputInfo.UsernameAdd == "" {
		return
	}
	var usernames []string
	if InputInfo.UsernameAdd != "" {
		usernames = strings.Split(InputInfo.UsernameAdd, ",")
	}

	for name := range common.Userdict {
		for _, username := range usernames {
			common.Userdict[name] = append(common.Userdict[name], username)
		}
		common.Userdict[name] = Reverse(common.Userdict[name])
	}
}

// ParsePass 解析密码:将-pwd和-pwdf中的密码提取出来，替换到Passwords中
func ParsePass(InputInfo *common.InputInfoStruct) {
	if InputInfo.PasswordAdd != "" {
		passs := strings.Split(InputInfo.PasswordAdd, ",")
		for _, pass := range passs {
			if pass != "" {
				common.Passwords = append(common.Passwords, pass)
			}
		}
	}
	common.Passwords = Reverse(common.Passwords)
}

func Reverse(input []string) []string {
	inputLen := len(input)
	output := make([]string, inputLen)

	for i, n := range input {
		j := inputLen - i - 1
		output[j] = n
	}

	return output
}
