package notice

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

var (
	title string
	level string
)

// SendDingTalkAlert 发送钉钉告警消息
func SendDingTalkAlert(webhookURL, ip, location string, isp string, project string, count int) error {
	timestamp := time.Now().Format(time.DateTime)

	// 判断级别并设置对应内容
	if count >= 500 {
		title = "告警通知"
		level = "告警"
	} else if count >= 230 {
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
        "text": "### **%s: %s 异常IP %s**\n#### 状态: 待处理\n- 五分钟内访问: %d次\n- IP地址: %s\n- 地理位置: %s\n- 运营商: %s\n- 执行时间: %s\n- [查看详情判断是否处理](https://iplark.com/%s )",
        "at": {
            "isAtAll": true
        }
    }
}`, title, project, level, count, ip, location, isp, timestamp, ip)

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
