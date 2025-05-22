package ipstat

import (
	"ipAnalyzer/internal/logparser"
	"strings"
	"time"
)

type IPStat struct {
	IP          string
	Count       int
	ProjectType string
	Level       string
}

type ByCount []IPStat

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].Count > a[j].Count }

func GetErrorLevel(count int) string {
	switch {
	case count >= 500:
		return "alert"
	case count >= 230:
		return "warning"
	default:
		return "info"
	}
}

// GetAnomalyIPs 统计指定时间窗口内 IP 的访问次数，并根据阈值判断是否为异常
func GetAnomalyIPs(entries []logparser.LogEntry, threshold int, windowDuration time.Duration) []IPStat {
	ipStats := make(map[string]int)
	now := time.Now()
	windowStart := now.Add(-windowDuration)

	for _, entry := range entries {
		if entry.Timestamp.Before(windowStart) {
			continue
		}

		ipStats[entry.IP]++
	}

	var anomalyIPs []IPStat

	for ip, count := range ipStats {
		level := GetErrorLevel(count)

		// 只保留达到或超过提醒级别（230）的 IP
		if level == "warning" || level == "alert" {
			anomalyIPs = append(anomalyIPs, IPStat{
				IP:    ip,
				Count: count,
				Level: level,
			})
		}
	}
	//data, _ := json.MarshalIndent(anomalyIPs, "", "  ")
	//fmt.Println("anomalyIPs:\n", string(data))
	return anomalyIPs
}

func GetProjectType(logFile string) string {
	if strings.Contains(logFile, "www.mod.gov.cn") {
		return "mod"
	} else if strings.Contains(logFile, "www.ccdi.gov.cn") {
		return "ccdi"
	}
	return "unknown"
}
