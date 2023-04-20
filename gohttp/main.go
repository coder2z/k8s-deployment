package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"runtime"
)

func main() {
	http.HandleFunc("/my", func(writer http.ResponseWriter, request *http.Request) {
		// 获取机器信息
		// 1. 获取主机名
		host, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}
		// 2. 获取主机IP
		ip := getLocalIp()
		// 3. 获取主机CPU核心数
		cpu := runtime.NumCPU()
		// 将机器信息写入到json文件中
		data := fmt.Sprintf(`{
			"host": "%s",
			"ip":   "%s",
			"cpu":  %d,
		}`, host, ip, cpu)

		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write([]byte(data))
	})

	http.HandleFunc("/debug/pprof", pprof.Index)

	log.Printf("Server is running at %s:8080 \n", getLocalIp())
	if err := http.ListenAndServe(":12500", nil); err != nil {
		log.Fatal("ListenAndServe error", err)
	}
}

func getLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}
