package main

import (
	"GScan/infoscan/api"
	"GScan/infoscan/config"
	"GScan/infoscan/dao/sqlite"
	"GScan/pkg/logger"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

import _ "net/http/pprof"

const debugmod = false

var a *api.Api
var Config *config.Config

func main() {
	// 启动HTTP API服务
	httpApi := api.NewHttpApi(a, sqlite.NewDB(filepath.Join(Config.ResultPath, "data.db")), Config)
	if err := httpApi.Start(":8080"); err != nil {
		log.Fatalf("启动HTTP API服务失败: %v", err)
	}
}

func banner() {
	fmt.Println(`---------------------------------------------------------------
Version: InfoScan 0.4.10
Email:   i@vshex.com
Github:  https://github.com/Ymjie/GScan
---------------------------------------------------------------`)
}

//func debugs() {
//	debug.SetGCPercent(10)
//	go func() {
//		ticker := time.NewTicker(10 * time.Second)
//		for {
//			debug.FreeOSMemory()
//			<-ticker.C
//		}
//	}()
//}

func init() {
	//debugs()
	banner()
	if debugmod {
		go func() {
			err := http.ListenAndServe("0.0.0.0:9991", nil)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
	//配置文件读取
	c, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	Config = c
	os.Mkdir(Config.ResultPath, 0644)
	os.Mkdir(Config.LogPath, 0644)
	//设置日志
	logger.Setallwriterlevel(Config.LogLevel)
	logger.SetStdoutLv(Config.LogPrintingLevel)
	logfile, _ := os.OpenFile(filepath.Join(Config.LogPath, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02 15-04-05"))), os.O_CREATE|os.O_RDWR, 0644)
	logger.SetAllwriter(logfile)
	// 数据库初始化
	DB := sqlite.NewDB(filepath.Join(Config.ResultPath, "data.db"))
	//api 初始化
	a = api.NewApi(DB, Config)
}
func geturl(urlpath string) ([]string, error) {
	// url 列表读取
	f, err := os.Open(urlpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}
