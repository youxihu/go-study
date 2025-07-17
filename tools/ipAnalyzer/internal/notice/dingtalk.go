package notice

import (
	"bytes"
	"fmt"
	"ipAnalyzer/internal/entity"
	"net/http"
	"time"
)

var (
	title string
	level string
)

// SendDingTalkAlert 发送钉钉告警消息
func SendDingTalkAlert(webhookURL, ip, location string, isp string, project string, count int, threshold entity.ThresholdConfig) error {
	execTime := time.Now().Format("15:04:01")
	actTime := time.Now().Add(-5 * time.Minute).Format("15:04:01")

	var title, level string
	if count >= threshold.Alert {
		title = "告警通知"
		level = "告警"
	} else if count >= threshold.Warning {
		title = "事件通知"
		level = "提醒"
	} else {
		// 不达到阈值不发送通知
		return nil
	}

	msg := fmt.Sprintf(`{
    "msgtype": "markdown",
    "markdown": {
        "title": "异常IP告警",
        "text": "### **%s: %s 异常IP %s**\n#### 状态: 待处理\n- 异常时间段: %s/%s\n- 异常期访问: %d次\n- IP地址: %s\n- 地理位置: %s\n- 运营商: %s\n- [点击判断是否介入处理](https://iplark.com/%s )",
        "at": {
            "isAtAll": true
        }
    }
}`, title, project, level, actTime, execTime, count, ip, location, isp, ip)

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer([]byte(msg)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("dingtalk alert failed with code: %d", resp.StatusCode)
	}

	fmt.Println("✅ 钉钉告警已发送")
	return nil
}
