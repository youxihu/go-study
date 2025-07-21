package geolocation

import (
	"fmt"
	"github.com/oschwald/maxminddb-golang"
	"ipAnalyzer/internal/entity"
	"ipAnalyzer/internal/notice"
	"ipAnalyzer/pkg/nginx"
	"log"
	"net"
	"os"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var Ip2RegionSearcher *xdb.Searcher

type ASRecord struct {
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
}

type MaxMindCityRecord struct {
	Country struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Subdivisions []struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
}

func InitIp2Region(dbFile string) error {
	db, err := os.Open(dbFile)
	if err != nil {
		return fmt.Errorf("failed to open ip2region db: %w", err)
	}
	defer db.Close()

	cBuffer, err := xdb.LoadContent(db)
	if err != nil {
		return fmt.Errorf("failed to load ip2region content: %w", err)
	}

	searcher, err := xdb.NewWithBuffer(cBuffer)
	if err != nil {
		return fmt.Errorf("failed to create ip2region searcher: %w", err)
	}

	Ip2RegionSearcher = searcher
	return nil
}

func ipStrToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address")
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return 0, fmt.Errorf("not an IPv4 address")
	}
	return uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3]), nil
}

func GetISPFromMaxMind(ipStr, asnDBPath string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return ""
	}

	asnDB, err := maxminddb.Open(asnDBPath)
	if err != nil {
		log.Printf("Error opening MaxMind ASN DB: %v", err) // Added logging
		return ""
	}
	defer asnDB.Close()

	var asRecord ASRecord
	err = asnDB.Lookup(ip4, &asRecord)
	if err != nil {
		log.Printf("Error looking up ASN in MaxMind DB for %s: %v", ipStr, err)
		return ""
	}

	return asRecord.AutonomousSystemOrganization
}

// GetLocationFromMaxMindCity 尝试从 MaxMind City 数据库获取国家、省份和城市信息。
func GetLocationFromMaxMindCity(ipStr, cityDBPath string) (country, province, city string) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", "", ""
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return "", "", ""
	}

	cityDB, err := maxminddb.Open(cityDBPath)
	if err != nil {
		log.Printf("Error opening MaxMind City DB: %v", err)
		return "", "", ""
	}
	defer cityDB.Close()

	var cityRecord MaxMindCityRecord
	err = cityDB.Lookup(ip4, &cityRecord)
	if err != nil {
		log.Printf("Error looking up City in MaxMind DB for %s: %v", ipStr, err)
		return "", "", ""
	}

	if val, ok := cityRecord.Country.Names["zh-CN"]; ok {
		country = val
	} else if val, ok := cityRecord.Country.Names["en"]; ok {
		country = val
	}

	if len(cityRecord.Subdivisions) > 0 {
		if val, ok := cityRecord.Subdivisions[0].Names["zh-CN"]; ok {
			province = val
		} else if val, ok := cityRecord.Subdivisions[0].Names["en"]; ok {
			province = val
		}
	}

	if val, ok := cityRecord.City.Names["zh-CN"]; ok {
		city = val
	} else if val, ok := cityRecord.City.Names["en"]; ok {
		city = val
	}

	return country, province, city
}

func JudgeWhiteList(ip string, whiteList []string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, entry := range whiteList {
		if strings.Contains(entry, "/") {
			// CIDR 格式
			_, cidr, err := net.ParseCIDR(entry)
			if err != nil {
				continue
			}
			if cidr.Contains(parsedIP) {
				return true
			}
		} else {
			// 单个 IP
			if parsedIP.String() == entry {
				return true
			}
		}
	}
	return false
}

func shouldSendAlert(country string, ip string, whiteList []string) bool {
	// 不是中国大陆，直接告警
	if country != "中国" && country != "CN" { // Ensure "CN" is also considered for China
		return true
	}

	// 中国IP，且不在白名单，才告警
	return !JudgeWhiteList(ip, whiteList)
}

func LookupLocationWithCount(ipStr string, count int, asnDBPath string, cityDBPath string, projectType string, webhookURL string, whiteList []string, threshold entity.ThresholdConfig) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		log.Printf("Invalid IP address: %s", ipStr)
		return
	}

	ip4 := ip.To4()
	if ip4 == nil {
		log.Printf("Not an IPv4 address: %s", ipStr)
		return
	}
	ipStr = ip4.String()

	if Ip2RegionSearcher == nil {
		log.Println("Ip2RegionSearcher not initialized")
		return
	}

	var (
		country  string
		province string
		city     string
		isp      string
	)

	// 首先尝试从 ip2region 获取地区信息
	ipUint32, err := ipStrToUint32(ipStr)
	if err != nil {
		log.Printf("将 IP 转换为 uint32 时出错（用于 ip2region）: %v", err)
		// 如果由于某些原因 ip2region 在此阶段失败，继续尝试使用 MaxMind
	} else {
		region, err := Ip2RegionSearcher.Search(ipUint32)
		if err != nil {
			log.Printf("Error searching ip2region for %s: %v", ipStr, err)
		}

		if region != "" {
			parts := strings.Split(region, "|")
			if len(parts) >= 5 {
				country = parts[0]
				province = parts[2]
				city = parts[3]
				isp = parts[4]

				if province == "0" {
					province = "" // 清空以 MaxMind 填充
				}
				if city == "0" {
					city = "" // 清空以 MaxMind 填充
				}
			}
		}
	}

	// 如果 ip2region 没有提供足够的信息，则尝试使用 MaxMind City DB
	// 特别是在国家为空，或者不是“中国”的情况下省份或城市为空
	if country == "" || (country != "中国" && country != "CN" && (province == "" || city == "")) {
		maxMindCountry, maxMindProvince, maxMindCity := GetLocationFromMaxMindCity(ipStr, cityDBPath)
		if maxMindCountry != "" {
			country = maxMindCountry
		}
		if maxMindProvince != "" {
			province = maxMindProvince
		}
		if maxMindCity != "" {
			city = maxMindCity
		}
	}

	// 总是尝试从 MaxMind 获取 ISP，因为通常更详细
	detailedISP := GetISPFromMaxMind(ipStr, asnDBPath)
	if detailedISP == "" {
		detailedISP = isp
	}

	// 如果仍然为空，则设置默认值
	if detailedISP == "" {
		detailedISP = "Unknown"
	}

	if country == "" {
		country = "Unknown"
	}
	if province == "" {
		province = "Unknown"
	}
	if city == "" {
		city = "Unknown"
	}

	// 构造用于显示/告警的 location 字段
	var location string
	switch {
	case province == "Unknown" && city == "Unknown":
		location = "Unknown"
	case province == "Unknown":
		location = city
	case city == "Unknown":
		location = province
	default:
		location = province + city
	}

	log.Printf("%3d次 %s  %s|%s|%s|%s\n", count, ipStr, country, province, city, detailedISP)

	if shouldSendAlert(country, ipStr, whiteList) {
		if count >= threshold.Error {
			err := nginx.BlockAttackerIP(ipStr)
			if err != nil {
				log.Printf("❌ 封禁 IP 失败: %v", err)
			}
		}
		err := notice.SendDingTalkAlert(webhookURL, ipStr, location, detailedISP, projectType, count, threshold)
		if err != nil {
			log.Printf("⚠️ 发送钉钉告警失败: %v\n", err)
		}
	}
}
