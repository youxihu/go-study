# 📊 IPAnalyzer

IPAnalyzer 是一个使用 Go 编写的高性能 Nginx 日志分析工具，具备自动运行、IP 异常检测、地理定位、白名单判断、钉钉通知和自动封禁功能，适用于生产环境下的安全告警与访问行为分析。

---

## ✨ 功能特性

- ⏱ 心跳分析日志 (心跳间隔 5 分钟)
- 📊 多项目日志支持（如 BBZ、BBX）
- 🧠 异常 IP 访问频次判断（warning / alert / error）
- 🌐 双引擎 IP 地理定位（ip2region + MaxMind）
- 🧾 可配置白名单策略，精准控制告警范围
- 🚫 支持封禁恶意访问 IP（集成 Nginx 封禁）
- 📣 钉钉 Markdown 消息通知，含位置、运营商等信息
- ☁️ 配置从 Nacos 加载，支持远程集中管理

---

## 🗂 项目结构

```
ipAnalyzer/
├── cmd/                  # 启动入口
├── internal/
│   ├── entity/           # 配置结构体定义
│   ├── geolocation/      # IP 定位逻辑
│   ├── ipstat/           # 异常 IP 统计与等级判定
│   ├── logparser/        # Nginx 日志解析模块
│   ├── notice/           # 钉钉消息发送模块
│   ├── processor/        # 日志处理主逻辑
├── pkg/
│   └── nacos/            # Nacos 客户端封装
│   └── nginx/            # Nginx 封禁逻辑
```

---

## ⚙️ 配置示例（通过 Nacos 加载）

```yaml
dbFilepath:
  ip2regionDBPath: pkg/ip2region.xdb              # 本地 IP 离线库（ip2region）
  asnDBPath: pkg/GeoLite2-ASN.mmdb                # MaxMind ASN 数据库（获取运营商信息）
  cityDBPath: pkg/GeoLite2-City.mmdb              # MaxMind City 数据库（国家、省份、城市）

logFilesPath:
  - nginx/nginx-logs/bbz_prod/bbz.com.log
  - nginx/nginx-logs/bbx_prod/bbx.com.log

whiteList:
  - 121.**.**.237
  - 112.**.**.90
  - 120.**.**.226
  - 47.**.**.146
  - 183.**.**.149

dingTalkWebhook: "https://oapi.dingtalk.com/robot/send?access_token=f06c66f2**b7**6e9f6ecd157d175e7bccd93391365b7675e6"

thresholds:
  warning: 400     # 警告等级：只提示
  alert: 600       # 告警等级：发送钉钉
  error: 800       # 严重等级：发送钉钉 + 封禁 IP
```

---

## 🚀 快速开始

1. 安装依赖：

```bash
go mod tidy
```

2. 启动服务：

```bash
go run cmd/main.go
```

3. 每 5 分钟会自动执行一次日志分析、告警判断与地理位置识别。

---

## 📑 日志格式要求（Nginx 示例）

程序要求日志中每行必须包含以下字段（以管道 | 分割）：

```
123.45.67.89 | Status: 200 | Time: [21/Jul/2025:14:30:01 +0800] | ...
```

> ⚠️ 会自动跳过 `Status: 403` 的请求

---

## 🧠 异常判定规则

系统基于近 5 分钟时间窗口内 IP 的访问次数进行等级分类：

| 等级    | 访问次数 | 处理动作            |
|---------|-----------|-----------------|
| warning | ≥ 400     | 发送钉钉提醒          |
| alert   | ≥ 600     | 发送钉钉告警          |
| error   | ≥ 800     | 发送安全告警并 封禁恶意 IP |

---

## 🌍 IP 定位逻辑

1. **优先使用 `ip2region` 本地库**
2. **如果信息不完整，补充使用 `MaxMind City/ASN`**
3. **自动提取 国家、省份、城市、ISP 等信息**
4. **非中国大陆 IP 或不在白名单的中国 IP 会被告警/封禁**

---

## 📣 钉钉告警示例

- 支持 Markdown 格式
- 可直接点击链接查看 IP 详情
- 包含访问频次、地理位置、ISP、项目类型等信息

---
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

### 安全通知: BBX 异常IP 封禁
#### 状态: 已封禁该IP
- 异常时间段: 16:06:07/16:11:07
- 异常期访问: 855次
- IP地址: 106.120.101.58
- 地理位置: 北京北京市
- 运营商: China Networks Inter-Exchange
- [点击查看访问者IP信息](https://iplark.com/116.179.37.68 )

---
### 告警通知: BBZ 异常IP 告警
#### 状态: 待处理
- 异常时间段: 16:06:07/16:11:07
- 异常期访问: 855次
- IP地址: 125.118.66.215
- 地理位置: 浙江省杭州市
- 运营商: China Networks Inter-Exchange
- [点击判断是否介入处理](https://iplark.com/116.179.37.68 )

---

### 事件通知: BBX 异常IP 提醒
#### 状态: 待处理
- 异常时间段: 11:55:07/12:00:07
- 异常期访问: 379次
- IP地址: 13.79.53.166
- 地理位置: Leinster都柏林
- 运营商: MICROSOFT-CORP-MSN-AS-BLOCK
- [点击判断是否介入处理](https://iplark.com/13.79.53.166)