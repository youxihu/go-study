package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

const (
	sshHost = "192.168.0.214:2345"
	sshUser = "youxihu"
	sshPass = "1"
)

func runSSHCommand(command string) (string, error) {
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", sshHost, config)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func parseCommandOutput(command, output string) gin.H {
	// 判断命令是否是 `free -h`
	if command == "free -h" {
		return parseFreeOutput(output)
	}

	// 对于其他命令，直接返回原始输出
	return gin.H{
		"result": output,
	}
}

func parseFreeOutput(output string) gin.H {
	parsedData := make(map[string]interface{})

	lines := strings.Split(output, "\n")
	if len(lines) > 1 {
		// 解析内存部分
		memFields := strings.Fields(lines[1])
		if len(memFields) > 6 {
			parsedData["Mem"] = map[string]string{
				"total": memFields[1],
				"used":  memFields[2],
				"free":  memFields[3],
			}
		}
	}

	return parsedData
}

func main() {
	r := gin.Default()

	r.POST("/ssh", func(c *gin.Context) {
		command := c.PostForm("cmd")
		if command == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "缺少命令"})
			return
		}

		output, err := runSSHCommand(command)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 根据命令判断是否需要格式化输出
		parsedOutput := parseCommandOutput(command, output)

		// 返回结果
		c.JSON(http.StatusOK, gin.H{
			"command":       command,
			"parsed_output": parsedOutput,
		})
	})

	fmt.Println("API 服务器启动: http://localhost:8080/ssh?cmd=uptime")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("API 服务器启动失败: %v", err)
	}
}
