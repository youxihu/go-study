# [go-study](https://github.com/youxihu)

本仓库为学习与实践项目集合，主要以 Golang 为主，涵盖运维日常、数据迁移脚本、数据库操作工具、AIOps相关 等内容。

## 项目结构说明

### [app/database](https://github.com/youxihu/go-study)

- 包含多个关于数据库操作的项目样例 
- 主要使用 [Ent](https://entgo.io/) 作为 ORM 工具进行数据库建模与访问 
- 实用工具记录: [SQL转Ent](https://old.printlove.cn/tools/sql2gorm/)
---
### [app/database/account](https://github.com/youxihu/go-study/tree/master/app/database/account)

一个迁移脚本项目，用于将 `Account` 表的数据迁移到 `AccountWallet` 表中，迁移逻辑基于 Ent 实现。

---
### [app/database/tabtaba](https://github.com/youxihu/go-study/tree/master/app/database/tabtaba)

项目导入工具，适用于产品经理与商务/电销团队之间的数据对接。  
功能包括：

- 从 Excel 表格中读取订单数据
- 根据不同 `业务项目` 导入至对应的项目数据库
- 支持定期导入任务
- 使用 Ent 进行数据库操作

---
### [tools/ai](https://github.com/youxihu/go-study/tree/master/tools/AI)

预留用于 AIOps 实现的测试脚本，目前尚未开发，仅为初步 demo。

---
### [tools/aliyunsms](https://github.com/youxihu/go-study/tree/master/tools/aliyunsms)

预留用于 AIOps 实现的短信通知通用模板， \
目前用[钉钉机器人](https://github.com/youxihu/dingtalk)代替，因此尚未开发，仅为初步 demo。

---
### [tools/ipAnalyzer](https://github.com/youxihu/go-study/tree/master/tools/ipAnalyzer)

用于 IP 地址信息分析的工具，详细说明请参考其独立的 [README 文档](https://github.com/youxihu/go-study/tree/master/tools/ipAnalyzer#readme)。

---
### [tools/nexus](https://github.com/youxihu/go-study/tree/master/tools/nexus)

提供对私有仓库（Sonatype Nexus）进行删、查的一系列 API 封装与调用示例。

---
### [tools/sslcheck](https://github.com/youxihu/go-study/tree/master/tools/sslCheck)

用于检测业务域名的 SSL 证书有效期，支持：

- 自动检测证书到期时间
- 自动申请并替换即将过期的证书
- 可集成至定时任务实现自动化证书管理


