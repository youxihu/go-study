package main

import (
	"fmt"
	"github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"math/rand"
	"strings"
	"time"
)

func main() {
	accessKeyId := "LTAI5t****84DPSGYno"     // 你的 AccessKeyId
	accessKeySecret := "Eap*****AwVIn358BnQ" // 你的 AccessKeySecret
	endpoint := "dysmsapi.a****.com"
	c := &client.Config{AccessKeyId: &accessKeyId, AccessKeySecret: &accessKeySecret, Endpoint: &endpoint}

	// 创建客户端
	newClient, err := dysmsapi20170525.NewClient(c)
	if err != nil {
		panic(fmt.Sprintf("Failed to create SMS client: %v", err))
	}

	phoneNumber := "19*****2" // 收短信的手机号码
	templateCode := "SMS_27***39"
	signName := "国安部" // 短信签名
	code := fmt.Sprintf("{\"code\":\"%s\"}", GenerateSmsCode(6))

	// 构造发送请求
	request := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &phoneNumber,
		TemplateCode:  &templateCode,
		SignName:      &signName,
		TemplateParam: &code,
	}

	// 发送短信
	sms, err := newClient.SendSms(request)
	if err != nil {
		panic(fmt.Sprintf("Failed to send SMS: %v", err))
	}

	// 打印返回结果
	fmt.Printf("SMS Send Result: %+v\n", sms)
}

// GenerateSmsCode 生成验证码; length 代表验证码的长度
func GenerateSmsCode(length int) string {
	// 生成一个数字数组
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().Unix()) // 程序启动时调用一次即可

	var sb strings.Builder
	for i := 0; i < length; i++ {
		// 随机选择数字并加入字符串构建器
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(len(numeric))])
	}
	return sb.String()
}
