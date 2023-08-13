## User Microservice
a simple user microservice demo, it uses grpc and mysql.

## Components
- gRPC
- GORM
- Viper
- Zap
- Lumberjack

## Getting Started
Clone this repository
```bash
git clone github.com/ramzeng/user-microservice
```
Adjust some configurations in configs/service.yaml
```yaml
service:
  port: 8080
auth:
  secret: "user-microservice"
  accessTokenTimeToLive: 3600
  refreshTokenTimeToLive: 86400
mysql:
  user: root
  password: password
  host: localhost
  port: 3306
  database: user-microservice
  maxIdleConnections: 10
  maxOpenConnections: 100
  connectionMaxIdleTime: 300
  connectionMaxLifetime: 3600
logger:
  channels:
    - name: app
      filename: logs/app.log
      maxSize: 500
      maxBackups: 10
      maxAge: 7
      compress: false
      level: debug
```
Start the service
```bash
go run examples/server.go -path="../configs/service.yaml" -migrate=true
```
Start the client
```bash
go run examples/client.go
```

