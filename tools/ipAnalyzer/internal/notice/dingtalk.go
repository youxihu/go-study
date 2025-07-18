package notice

import (
	"bytes"
	"fmt"
	"ipAnalyzer/internal/entity"
	"log"
	"net/http"
	"time"
)

// 全局复用的 HTTP Client
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// SendDingTalkAlert 发送钉钉告警消息
func SendDingTalkAlert(webhookURL, ip, location string, isp string, project string, count int, threshold entity.ThresholdConfig) error {
	execTime := time.Now().Format("15:04:01")
	actTime := time.Now().Add(-5 * time.Minute).Format("15:04:01")

	var title, level, status, link string

	switch {
	case count >= threshold.Error:
		title = "安全通知"
		level = "封禁"
		status = "已封禁该IP"
		link = "点击查看访问者IP信息"
	case count >= threshold.Alert:
		title = "告警通知"
		level = "告警"
		status = "待处理"
		link = "点击判断是否介入处理"
	case count >= threshold.Warning:
		title = "事件通知"
		level = "提醒"
		status = "需关注"
		link = "点击判断是否介入处理"
	default:
		return nil
	}

	markdown := fmt.Sprintf(`### **%s: %s 异常IP %s**
#### 状态: %s
- 异常时间段: %s/%s
- 异常期访问: %d次
- IP地址: %s
- 地理位置: %s
- 运营商: %s
- [%s](https://iplark.com/%s )`,
		title, project, level, status,
		actTime, execTime, count, ip, location, isp,
		link, ip,
	)

	msg := fmt.Sprintf(`{
    "msgtype": "markdown",
    "markdown": {
        "title": "异常IP告警",
        "text": "%s",
        "at": {
            "isAtAll": true
        }
    }
}`, markdown)

	err := sendDingTalkMessage(webhookURL, msg)
	if err != nil {
		log.Printf("❌ 钉钉告警发送失败: %v", err)
		return err
	}

	log.Printf("✅ 钉钉告警已发送: %s", ip)
	return nil
}

// sendDingTalkMessage 发送消息到钉钉 Webhook
func sendDingTalkMessage(webhookURL, message string) error {
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送钉钉消息失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("钉钉告警失败，状态码: %d", resp.StatusCode)
	}

	return nil
}
