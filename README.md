[English](./README.en.md) | 简体中文
# GinChat 即时通讯系统（Go后端）
基于 Go + Gin + WebSocket + MySQL + Redis+ Kafka 实现的高性能 IM 后端，支持高并发消息推送、在线状态管理、多级缓存优化与容器化部署。
### 联系我
qq邮箱:3110940369@qq.com

## 部署
### 1. 克隆项目到本地
### 2. 修改config/application.example.yaml为config/application.yaml,并修改dsn
### 3. 修改docker-compose.example.yml为docker-compose.yml,并修改mysql的dsn和密码
### 4. 需要修改的地方已经用TODO标记
### 5. 如果想直接看到项目效果请键入命令(前提有docker环境): docker compose up -d
### 6. 启动成功后直接在浏览器输入localhost就能看到前端页面(端口是80，所以不需要输入端口)


## 技术栈
- Golang、Gin、GORM
- MySQL、Redis、Kafka
- WebSocket 长连接、JWT 统一鉴权
- Docker、Docker Compose 容器化部署
- Swagger API 文档、Viper 配置管理
- 采用bcrypt加密,安全性比MD5更高

## 核心功能
- 单聊、群聊、消息实时推送
- 好友管理、会话列表、未读消息计数
- 用户在线状态维护
- Redis 缓存好友列表、在线状态、群成员id信息
- kafka使用在多个接口，明显降低接口耗时
- JWT 无状态身份认证，支持 HTTP 与 WebSocket 统一鉴权
- bcrypt加密用户密码
- 支持文件上传、接口文档自动生成
- Docker Compose 一键部署、多环境兼容

## 项目亮点
- 模仿spring实现了三层解耦和简易的依赖注入
- 使用kafka异步入库，死信队列兜底，消费失败重试
- 数据库使用联合唯一索引，消息使用雪花id使消费端存储层幂等，拦截重复消费
- 采用 WebSocket 双协程读写分离架构，Channel 消息队列保证高并发稳定推送、
- 自定义结构体存储ws连接,消息信道，每个连接拥有自己的读写协程，并实现写入消息->写协程串行化，保证并发安全。
- MySQL + Redis 多级存储，缓存热点数据，降低数据库压力
- 完整工程化结构，代码规范，易于维护与扩展
- 容器化编排，部署简单，可快速迁移上线

