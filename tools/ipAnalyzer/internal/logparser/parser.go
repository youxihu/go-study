package logparser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	IP        string
}

func ParseLogFile(logFilePath string) ([]LogEntry, error) {
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	// 设置为英语环境解析月份
	location, _ := time.LoadLocation("Asia/Shanghai") // 固定为北京时间

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		ipEnd := strings.IndexAny(line, " |")
		if ipEnd == -1 {
			continue
		}
		ip := line[:ipEnd]

		timeStart := strings.Index(line, "Time: [")
		if timeStart == -1 {
			continue
		}
		timeStr := line[timeStart+7 : timeStart+7+20] // 提取 "[16/May/2025:14:08:05 +0800]" 中的时间部分

		// 使用 ParseInLocation 强制用英语解析月份
		timestamp, err := time.ParseInLocation("2/Jan/2006:15:04:05", timeStr, location)
		if err != nil {
			fmt.Printf("Failed to parse timestamp: %s\n", timeStr)
			continue
		}

		entries = append(entries, LogEntry{
			Timestamp: timestamp,
			IP:        ip,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
