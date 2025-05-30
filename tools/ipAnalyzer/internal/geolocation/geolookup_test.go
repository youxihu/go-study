package geolocation

import (
	"testing"
)

func TestLookupLocationWithCount(t *testing.T) {
	// 初始化 ip2region 数据库（路径请替换为你本地的 ip2region.xdb）
	err := InitIp2Region("/home/youxihu/mywork/myproject/ipAnalyzer/pkg/ip2region.xdb")
	if err != nil {
		t.Fatalf("初始化 ip2region 失败: %v", err)
	}

	// MaxMind 数据库路径（请确保这些文件存在）
	asnDBPath := "/home/youxihu/mywork/myproject/ipAnalyzer/pkg/GeoLite2-ASN.mmdb"
	cityDBPath := "/home/youxihu/mywork/myproject/ipAnalyzer/pkg/GeoLite2-City.mmdb"

	// 测试 IP 列表
	testIPs := []string{
		"8.8.8.8",              // Google DNS - 美国
		"1.1.1.1",              // Cloudflare - 美国
		"114.114.114.114",      // DNSPod - 中国
		"180.76.9.138",         // Baidu - 中国
		"223.5.5.5",            // 阿里云 DNS - 中国
		"119.29.29.29",         // 腾讯云 DNS - 中国
		"invalid.ip.address",   // 无效 IP
		"2001:4860:4860::8888", // IPv6 不处理
	}

	// 模拟参数
	projectType := "TestProject"
	webhookURL := ""
	whiteList := []string{}

	// 对每个 IP 执行 LookupLocationWithCount
	for _, ip := range testIPs {
		t.Logf("正在测试 IP: %s", ip)
		LookupLocationWithCount(ip, 1, asnDBPath, cityDBPath, projectType, webhookURL, whiteList)
	}
}
