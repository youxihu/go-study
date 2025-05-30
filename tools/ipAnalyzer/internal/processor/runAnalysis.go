package processor

import (
	"gopkg.in/yaml.v3"
	"ipAnalyzer/internal/entity"
	"ipAnalyzer/internal/geolocation"
	"ipAnalyzer/internal/ipstat"
	"ipAnalyzer/pkg/nacos"
	"log"
	"sync"
)

func RunAnalysis() {
	// 1. 加载 Nacos 认证配置
	nacosConfig, err := nacos.LoadNacosAuth("/home/youxihu/secret/aiops/ipanalyzer/localtest_nacos_auth.yaml")
	if err != nil {
		log.Fatalf("加载 Nacos 认证失败: %v", err)
	}

	// 2. 创建 Nacos 客户端并获取远程配置内容
	configContent, err := nacos.CreateNacosClient(nacosConfig)
	if err != nil {
		log.Fatalf("创建 Nacos 客户端失败: %v", err)
	}

	// 3. 解析配置到结构体
	var appConfig entity.Config
	err = yaml.Unmarshal([]byte(configContent), &appConfig)
	if err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}
	// 4. 初始化 geolocation 模块
	err = geolocation.InitIp2Region(appConfig.DBFilePath.IP2RegionDBPath)
	if err != nil {
		log.Fatalf("InitIp2Region failed: %v", err)
	}

	var wg sync.WaitGroup
	resultChan := make(chan []ipstat.IPStat, len(appConfig.LogFilesPath))

	// 启动多个 goroutine 处理日志文件
	for _, logFile := range appConfig.LogFilesPath {
		wg.Add(1)
		go ProcessLogFile(logFile, &wg, resultChan)
	}

	// 等待所有 goroutine 完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集中间结果并发送通知
	for anomalyIPs := range resultChan {
		for _, item := range anomalyIPs {
			geolocation.LookupLocationWithCount(item.IP, item.Count, appConfig.DBFilePath.ASNDBPath, appConfig.DBFilePath.CITYDBPath, item.ProjectType, appConfig.DingTalkWebhook, appConfig.WhiteList)
		}
	}
	log.Println("所有日志文件分析完成")
}
