package geolocation

import (
	"fmt"
	"github.com/oschwald/maxminddb-golang"
	"ipAnalyzer/internal/notice"
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
		return ""
	}
	defer asnDB.Close()

	var asRecord ASRecord
	err = asnDB.Lookup(ip4, &asRecord)
	if err != nil {
		return ""
	}

	return asRecord.AutonomousSystemOrganization
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

func shouldSendAlert(country, isp string, ip string, whiteList []string) bool {
	// 如果不是中国大陆，强制发告警
	if country != "中国" && country != "CN" {
		return true
	}

	// 判断是否是中国三大运营商之一
	isp = strings.ToLower(isp)
	isMainISP := strings.Contains(isp, "chinanet") ||
		strings.Contains(isp, "unicom") ||
		strings.Contains(isp, "mobile")

	// 如果是三大运营商，并且在白名单中，则不发告警
	if isMainISP && JudgeWhiteList(ip, whiteList) {
		return false
	}
	return true
}

func LookupLocationWithCount(ipStr string, count int, asnDBPath string, projectType string, webhookURL string, whiteList []string) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return
	}
	ipStr = ip4.String()

	if Ip2RegionSearcher == nil {
		log.Println("Ip2RegionSearcher not initialized")
		return
	}

	ipUint32, err := ipStrToUint32(ipStr)
	if err != nil {
		return
	}

	region, err := Ip2RegionSearcher.Search(ipUint32)

	if err != nil || region == "" {
		return
	}

	parts := strings.Split(region, "|")
	if len(parts) < 4 {
		return
	}

	country := parts[0]
	province := parts[2]
	city := parts[3]
	isp := parts[4]

	if province == "0" {
		province = city
	}
	if city == province {
		city = ""
	}

	detailedISP := GetISPFromMaxMind(ipStr, asnDBPath)
	if detailedISP == "" {
		detailedISP = isp
	}

	log.Printf("%3d次 %s  %s|%s|%s|%s\n", count, ipStr, country, province, city, detailedISP)

	if shouldSendAlert(country, detailedISP, ipStr, whiteList) {
		err := notice.SendDingTalkAlert(webhookURL, ipStr, province+city, detailedISP, projectType, count)
		if err != nil {
			log.Printf("⚠️ 发送钉钉告警失败: %v\n", err)
		}
	}
}
