package processor

import (
	"ipAnalyzer/internal/ipstat"
	"ipAnalyzer/internal/logparser"
	"log"
	"sync"
	"time"
)

const (
	ipThreshold    = 230             // 异常访问次数阈值
	windowDuration = 5 * time.Minute // 时间窗口
)

func ProcessLogFile(logFile string, wg *sync.WaitGroup, resultChan chan<- []ipstat.IPStat) {
	defer wg.Done()

	log.Printf("正在分析日志文件: %s\n", logFile)

	// 根据文件路径判断项目类型
	projectType := ipstat.GetProjectType(logFile)

	// 1. 解析日志文件
	entries, err := logparser.ParseLogFile(logFile)
	if err != nil {
		log.Printf("解析日志失败 [%s]: %v\n", logFile, err)
		return
	}

	// 2. 获取异常 IP
	anomalyIPs := ipstat.GetAnomalyIPs(entries, ipThreshold, windowDuration)
	if len(anomalyIPs) == 0 {
		log.Printf("[%s] 未发现异常 IP\n", logFile)
		return
	}

	// 3. 将项目类型附加到每个 IPStat 上
	for i := range anomalyIPs {
		anomalyIPs[i].ProjectType = projectType
	}

	// 4. 把结果发送到 channel
	resultChan <- anomalyIPs
}
