package common

import "fmt"

func PrintBanner() {
	banner := ` _  _     ___     ___     ___   
| \| |   /   \   / __|   / __|  
| .  |   | - |  | (__    \__ \
|_|\_|   |_|_|   \___|   |___/  
             Version: ` + version
	fmt.Println(banner)
}
