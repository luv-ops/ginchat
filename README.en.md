English | [简体中文](./README.md)

# GinChat Instant Messaging System (Go Backend)
A high-performance IM backend built with Go + Gin + WebSocket + MySQL + Redis + Kafka. It supports high-concurrency message pushing, online status management, multi-level caching optimization, and containerized deployment.

### Contact Me
Email: 3110940369@qq.com

## Deployment
### 1. Clone the project to your local machine
### 2. Rename `config/application.example.yaml` to `config/application.yaml`, and modify the DSN
### 3. Rename `docker-compose.example.yml` to `docker-compose.yml`, and modify the MySQL DSN and password
### 4. Places that need modification are marked with `TODO`
### 5. If you want to see the project in action directly, enter the command (Docker environment required): `docker compose up -d`
### 6. After successful startup, enter `localhost` in your browser to see the frontend page (port 80, so no port number is needed)

## Tech Stack
- Golang, Gin, GORM
- MySQL, Redis, Kafka
- WebSocket long connections, JWT unified authentication
- Docker & Docker Compose containerized deployment
- Swagger API documentation, Viper configuration management
- bcrypt encryption (higher security than MD5)

## Core Features
- One-on-one chat, group chat, real-time message pushing
- Friend management, conversation list, unread message count
- User online status maintenance
- Redis caching for friend lists, online status, and group member IDs
- Kafka applied to multiple APIs, significantly reducing interface latency
- JWT stateless authentication, supporting unified auth for both HTTP and WebSocket
- bcrypt password encryption
- File upload support, automated API documentation generation
- Docker Compose one-click deployment, multi-environment compatibility

## Project Highlights
- Three-layer decoupling and simple dependency injection inspired by Spring
- Kafka asynchronous processing with dead-letter queue fallback and consumption failure retry
- Database composite unique indexes and Snowflake IDs for message idempotency at the storage layer, preventing duplicate consumption
- WebSocket dual-coroutine read-write separation architecture, utilizing Channel message queues to ensure stable high-concurrency pushing
- Custom structs to store WS connections and message channels; each connection has its own read/write coroutines, implementing write message -> write coroutine serialization to ensure concurrency safety
- MySQL + Redis multi-level storage, caching hot data to reduce database pressure
- Complete engineering structure, clean code standards, easy to maintain and extend
- Containerized orchestration, simple deployment, quick migration to production
