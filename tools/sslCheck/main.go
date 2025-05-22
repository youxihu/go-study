package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	dt "github.com/youxihu/dingtalk/dingtalk"
)

// 定义域名数组
var domains = []string{
	"www.mps.gov.cn",
	"www.mod.gov.cn",
	"www.ccdi.gov.cn",
}

// 日志文件路径
const logFilePath = "/notice/ssl_check.log"

// 初始化日志记录器，同时输出到控制台和日志文件
func initLogger() {
	// 打开或创建日志文件
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}

	// 多写入器：同时写入到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 设置默认日志输出
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // 包含日期、时间、文件行号
}

// 获取 SSL 证书的过期时间
func getCertificateExpiry(domain string) (time.Time, error) {
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return time.Time{}, err
	}
	defer conn.Close()
	cert := conn.ConnectionState().PeerCertificates[0]
	expiryDate := cert.NotAfter
	return expiryDate, nil
}

// 检查 SSL 证书是否将在 15 天内过期
func checkExpiry(domain string) string {
	expiryDate, err := getCertificateExpiry(domain)
	if err != nil {
		log.Printf("无法获取域名 %s 的证书信息: %v\n", domain, err)
		return ""
	}

	currentTime := time.Now()
	daysRemaining := int(expiryDate.Sub(currentTime).Hours() / 24)

	if daysRemaining <= 15 {
		return fmt.Sprintf("### **告警通知: SSL证书即将过期**\n\n#### 状态: 待处理\n\n#### 证书详情：\n- 检查时间: `%s`\n- 域名: [%s](https://%s)\n- 剩余天数: `%d` 天\n- 过期时间: `%s`\n- 提醒: 请及时更新证书。\n", currentTime.Format(time.DateTime), domain, domain, daysRemaining, expiryDate.Format(time.DateTime))
	}
	return ""
}

func main() {
	// 初始化日志记录
	initLogger()

	// 测试群钉钉配置
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=606800eb6413d83a0b42a"
	secret := "SEC7e7bc8109e21133b3fe90d853d584ad311b2f6"

	atMobiles := []string{"19****42"} // 艾特指定手机号
	isAtAll := false                  // 是否艾特所有人

	// 循环检查每个域名的证书
	for _, domain := range domains {
		warning := checkExpiry(domain)
		if warning != "" {
			title := "证书过期提醒"
			text := warning

			err := dt.SendDingDingNotification(webhookURL, secret, title, text, atMobiles, isAtAll)
			if err != nil {
				log.Printf("发送通知失败 for domain %s: %v\n", domain, err)
			} else {
				log.Printf("通知已发送成功 for domain %s!\n", domain)
			}
		} else {
			log.Printf("域名 %s 的 SSL 证书正常。\n", domain)
		}
	}
}
