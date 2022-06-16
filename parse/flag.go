package parse

import (
	"flag"
	"nacs/common"
)

func Flag(Info *common.InputInfoStruct) {
	flag.StringVar(&Info.Host, "h", "", "IP address of the host you want to scan,for example: 192.168.11.11 | 192.168.11.11-255 | 192.168.11.11,192.168.11.12")
	flag.StringVar(&Info.HostFile, "hf", "", "File of IP address")

	flag.StringVar(&Info.Mode, "m", "all", "Select scan mode")
	flag.StringVar(&Info.OutputFileName, "o", "output.txt", "")

	flag.BoolVar(&Info.NoProbe, "np", false, "no probe")
	flag.BoolVar(&Info.Silent, "silent", false, "")
	flag.BoolVar(&Info.NoSave, "ns", false, "not save file")
	flag.BoolVar(&Info.NoColor, "nc", false, "no color")
	flag.IntVar(&Info.LogLevel, "loglevel", 3, "log level")

	flag.IntVar(&Info.LiveTop, "top", 10, "show live len top")

	flag.StringVar(&Info.PortsOnly, "po", "", "only use these ports")
	flag.StringVar(&Info.PortsAdd, "pa", "", "add these ports")

	flag.IntVar(&Info.Thread, "t", 1000, "Number of concurrent threads")

	flag.StringVar(&Info.Proxy, "proxy", "", "Use proxy scan, support http/socks5 protocol [e.g. --proxy socks5://127.0.0.1:1080]")

	flag.IntVar(&Info.Timeout, "timeout", 3, "Response timeout time, the default is 3 seconds")
	flag.StringVar(&Info.DiscoverMode, "mode", "", "Specify the protocol [e.g. -m mysql/-m http]")
	flag.StringVar(&Info.DiscoverType, "type", "", "Specify the type [e.g. --type tcp/--type udp]")

	flag.StringVar(&Info.OutJson, "json", "output.json", "json")

	flag.StringVar(&Info.CeyeKey, "ceyekey", "", "ceye.io api key")
	flag.StringVar(&Info.CeyeDomain, "ceyedomain", "", "ceye.io api subdomain")

	flag.BoolVar(&Info.NoBrute, "nobrute", false, "no brute")
	flag.BoolVar(&Info.NoPoc, "nopoc", false, "no poc")

	flag.IntVar(&Info.PocRate, "pocrate", 20, "Request rate(per second)")
	flag.BoolVar(&Info.PocDebug, "pocdebug", false, "print failed pocs")
	flag.IntVar(&Info.PocThread, "pocthread", 10, "poc thread")
	flag.IntVar(&Info.PocTimeout, "poctimeout", 20, "poc timeout")
	flag.StringVar(&Info.NucleiPocPath, "nucleipocpath", "pocs/nuclei/**", "Nuclei poc path") //  pocs/nuclei/**
	flag.StringVar(&Info.FscanPocPath, "fscanpocpath", "pocs/xrayv1/", "Fscan poc path")      // pocs/fscan/**
	flag.BoolVar(&Info.NoNuclei, "nonuclei", false, "no nuclei")

	flag.StringVar(&Info.Command, "command", "whoami", "exec command (ssh)")
	flag.StringVar(&Info.SSHKey, "sshkey", "", "sshkey file (id_rsa)")
	flag.BoolVar(&Info.BruteDebug, "brutedebug", false, "print failed attempts")
	flag.IntVar(&Info.BruteTimeout, "brutetimeout", 1, "brute timeout time, the default is 1 seconds")
	flag.StringVar(&Info.BruteSocks5Proxy, "bruteproxy", "", "brute proxy")
	flag.IntVar(&Info.BruteThread, "brutethread", 10, "brute thread")
	flag.StringVar(&Info.UsernameAdd, "usernameadd", "", "username add: split by [,]")
	flag.StringVar(&Info.PasswordAdd, "passwordadd", "", "password add: split by [,]")

	flag.StringVar(&Info.RedisFile, "redisfile", "id_rsa.pub", "redis file to write sshkey file (as: -rf id_rsa.pub) ")
	// redis反弹shell的目标
	flag.StringVar(&Info.RedisShell, "redishell", "", "redis shell to write cron file (as: -rs 192.168.1.1:6666) ")

	flag.StringVar(&Info.DirectUrl, "u", "", "url split by [,] i.e., http(s), ssh, ftp, ...")
	flag.StringVar(&Info.DirectUrlFile, "uf", "", "url file")
	flag.BoolVar(&Info.DirectUrlForce, "uforce", false, "always exploit")

	flag.Parse()
}
