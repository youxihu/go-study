# 📊 IP Analyzer 使用文档

> 一个基于日志文件分析 IP 地理位置和异常行为的 Go 监控工具

---

## 🧩 项目概述

该项目用于分析日志文件中的 IP 数据，并结合地理位置数据库（如 `ip2region` 和 `ASN`）识别异常访问行为，最终通过钉钉 Webhook 发送告警通知。

* 使用 Nacos 获取远程配置
* 使用多线程并发处理多个日志文件
* 检测异常 IP （基于阈值判断频率）
* 分析异常 IP 做地理定位（ASN / Region）
* 确认异常后通过 channel 发送 DingTalk

### ✅ 阈值设定标准

| 行为类型           | 5分钟内访问次数    | 说明                                 | 处理方式 |
| :----------------- | :----------------- | :----------------------------------- | :------- |
| 正常用户行为       | ≤ 100次/5分钟      | 浏览网页、API调用等正常操作          | 忽略     |
| 可疑行为（初步判定） | > 100次/5分钟      | 建议记录并进一步观察                 | 提醒     |
| 异常/攻击行为（高概率）| > 300次/5分钟      | 很可能是爬虫、爆破、CC攻击等         | 告警     |

---

## 🗂️ 目录结构
```
├── cache                       # 项目启动后生成,需正确配置Nacos鉴权文件
│   └── nacos
│       └── config
├── cmd
│   └── main.go
├── internal
│   ├── entity
│   ├── geolocation
│   ├── ipstat
│   ├── logparser
│   ├── notice
│   └── processor
├── pkg
│   ├── GeoLite2-ASN.mmdb        # ASN 地理数据库
│   ├── GeoLite2-City.mmdb       # 城市区分数据库
│   ├── ip2region.xdb            # 国内 IP 定位数据库
│   └── nacos
└── README.md
```
📄 ```nacos_auth.yaml``` Nacos鉴权文件
```
auth:
  host: 127.0.0.1
  port: 8848
  username: aiops-user
  password: aiops-user
  namespace_id: c2n96m6d-3306-543c-b37c-f1887415157
  timeout_ms: 500
  log_dir: cache/nacos/log/
  cache_dir: cache/nacos/
  log_level: debug
  data_id: ip_analyzer
  group: DEFAULT_GROUP
```
📄 ```IP_Analyzer.yaml``` 远程配置文件
```
dbfilepath:
  ip2regionDBPath: pkg/ip2region.xdb
  asnDBPath: pkg/GeoLite2-ASN.mmdb

logFilesPath:
  - /app/nginx/nginx-logs/www.ccdi.gov.cn.log
  - /app/nginx/nginx-logs/www.mod.gov.cn.log

dingTalkWebhook: "https://oapi.dingtalk.com/robot/send?access_token=f06c66f258cd18874151577bccd93391365b7675e6"
```
## ▶️ 使用
```text
1. go build -o ipanalyzer cmd/main.go
2. ./ipanalyzer
3.log:
2025/05/21 09:52:55 正在分析日志文件: /app/nginx/nginx-logs/www.ccdi.gov.cn.log
2025/05/21 09:52:55 [/app/nginx/nginx-logs/www.ccdi.gov.cn.log] 未发现异常 IP
2025/05/21 09:52:55 正在分析日志文件: /app/nginx/nginx-logs/www.mod.gov.cn.log
2025/05/21 09:52:55 [/app/nginx/nginx-logs/www.mod.gov.cn.log] 未发现异常 IP
2025/05/21 09:52:55 所有日志文件分析完成
```
## ⚠️ 告警示例

### 告警通知: BBZ异常IP检测
#### 状态: 待处理
- IP地址: 116.179.37.68
- 地理位置: 山西省阳泉市
- 运营商: CHINA UNICOM China169 Backbone
- 执行时间: 2025-05-20 10:58:53
- [查看详情判断是否处理](https://iplark.com/116.179.37.68 )
