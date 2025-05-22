package main

import (
	"ipAnalyzer/internal/processor"
	"time"
)

func main() {
	processor.RunAnalysis()

	// 设置每 5 分钟执行一次
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		processor.RunAnalysis()
	}
}
