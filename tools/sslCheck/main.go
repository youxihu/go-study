package main

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"time"

	dt "github.com/youxihu/dingtalk/dingtalk"
)

// 配置结构体
type DingTalkConfig struct {
	WebhookURL string   `yaml:"webhook_url"`
	Secret     string   `yaml:"secret"`
	AtMobiles  []string `yaml:"at_mobiles"`
	IsAtAll    bool     `yaml:"is_at_all"`
}

type Config struct {
	Domains  []string       `yaml:"domains"`
	DingTalk DingTalkConfig `yaml:"dingtalk"`
}

// 日志文件路径
const logFilePath = "/notice/ssl_check.log"
const configPath = "/notice/ssl_check.yaml"

// 初始化日志记录器，同时输出到控制台和日志文件
func initLogger() {
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
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
		return fmt.Sprintf("### **告警通知: SSL证书即将过期**\n\n#### 状态: 待处理\n\n#### 证书详情：\n- 检查时间: `%s`\n- 域名: [%s](https://%s)\n-  剩余天数: `%d` 天\n- 过期时间: `%s`\n- 提醒: 请及时更新证书。\n", currentTime.Format(time.DateTime), domain, domain, daysRemaining, expiryDate.Format(time.DateTime))
	}
	return ""
}

// 读取 YAML 配置文件
func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	initLogger()

	// 加载配置文件
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 使用配置中的钉钉参数
	webhookURL := config.DingTalk.WebhookURL
	secret := config.DingTalk.Secret
	atMobiles := config.DingTalk.AtMobiles
	isAtAll := config.DingTalk.IsAtAll

	// 检查每个域名
	for _, domain := range config.Domains {
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
