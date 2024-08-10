package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/xiaoqidun/qqwry"
)

var (
	logs   string
	ipdata string
	domain string
)

func init() {
	flag.StringVar(&logs, "l", "./1.log", "日志")
	flag.StringVar(&ipdata, "i", "./qqwry.dat", "IP数据库")
	flag.StringVar(&domain, "d", "", "包含")
	flag.Parse()
}

func main() {
	bts, err := os.ReadFile(logs)
	if err != nil {
		panic(err)
	}

	// 从文件加载IP数据库
	if err := qqwry.LoadFile(ipdata); err != nil {
		panic(err)
	}

	// 正则匹配ipv4地址
	reg := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}).*accepted.*` + domain)

	var cache []string

	lines := strings.Split(string(bts), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		clientIps := reg.FindStringSubmatch(line)
		if len(clientIps) == 2 {
			if !slices.Contains(cache, clientIps[1]) {
				cache = append(cache, clientIps[1])
			}
		}
	}

	lenStr := strconv.Itoa(len(strconv.Itoa(len(cache))))
	for i, v := range cache {
		// 从内存或缓存查询IP
		location, err := qqwry.QueryIP(v)
		if err != nil {
			fmt.Printf("错误：%v\n", err)
			continue
		}

		fmt.Printf("%0"+lenStr+"d: %s"+strings.Repeat(" ", 15-len(v))+" 国家 %s, 省份 %s, 城市 %s, 区县 %s, 运营商：%s\n",
			i+1,
			v,
			location.Country,
			location.Province,
			location.City,
			location.District,
			location.ISP,
		)
	}
}
