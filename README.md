# GinChat 即时通讯系统（Go后端）
基于 Go + Gin + WebSocket + MySQL + Redis 实现的高性能 IM 后端，支持高并发消息推送、在线状态管理、多级缓存优化与容器化部署。

## 技术栈
- Golang、Gin、GORM
- MySQL、Redis
- WebSocket 长连接、JWT 统一鉴权
- Docker、Docker Compose 容器化部署
- Swagger API 文档、Viper 配置管理
- MD5 加盐加密、跨域处理

## 核心功能
- 单聊、群聊、消息实时推送
- 好友管理、会话列表、未读消息计数
- 用户在线状态维护
- Redis 缓存好友列表、在线状态、群成员id信息
- JWT 无状态身份认证，支持 HTTP 与 WebSocket 统一鉴权
- 支持文件上传、接口文档自动生成
- Docker Compose 一键部署、多环境兼容

## 项目亮点
- 采用 WebSocket 双协程读写分离架构，Channel 消息队列保证高并发稳定推送
- MySQL + Redis 多级存储，缓存热点数据，降低数据库压力
- 完整工程化结构，代码规范，易于维护与扩展
- 容器化编排，部署简单，可快速迁移上线
