package nginx

import (
	"fmt"
	"log"
	"os/exec"
)

func BlockAttackerIP(ip string) error {
	cmd := exec.Command("sudo", "/home/youxihu/tools/ipanalyzer/block_ip.sh", ip)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to block IP %s: %v, output: %s", ip, err, output)
	}

	log.Printf("✅ IP %s 封禁成功: %s", ip, output)
	return nil
}
