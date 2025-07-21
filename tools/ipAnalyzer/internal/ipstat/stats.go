package ipstat

import (
	"ipAnalyzer/internal/entity"
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

func GetErrorLevel(count int, thresholds entity.ThresholdConfig) string {
	switch {
	case count >= thresholds.Error:
		return "error"
	case count >= thresholds.Alert:
		return "alert"
	case count >= thresholds.Warning:
		return "warning"
	default:
		return "info"
	}
}

func GetAnomalyIPs(entries []logparser.LogEntry, thresholds entity.ThresholdConfig, windowDuration time.Duration) []IPStat {
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
		level := GetErrorLevel(count, thresholds)

		if level == "warning" || level == "alert" || level == "error" {
			anomalyIPs = append(anomalyIPs, IPStat{
				IP:    ip,
				Count: count,
				Level: level,
			})
		}
	}

	return anomalyIPs
}

func GetProjectType(logFile string) string {
	if strings.Contains(logFile, "www.51bbz.com.log") {
		return "BBZ"
	} else if strings.Contains(logFile, "online.biaobiaoxing.com.log") {
		return "BBX"
	}
	return "unknown"
}
