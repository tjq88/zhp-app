# zhp-app

`zhp-app` 是一个基于 Gin、GORM、MySQL 和 Redis 的 Go Web 服务。
当前代码已经具备会员注册链路，以及后续继续扩展业务模块所需的基础设施能力。

## 项目结构

```text
cmd/app/main.go                  应用入口，负责启动流程编排
internal/app/router.go           HTTP 路由注册
internal/handler                 HTTP 请求处理和响应封装
internal/middleware              中间件，如上下文注入和访问日志
internal/model                   持久化模型和接口视图模型
internal/service                 业务服务层
pkg/common                       MySQL、Redis 等共享基础设施组件
pkg/config                       基于环境变量的运行时配置
pkg/idgenx                       雪花算法 workerId 注册与 ID 生成器初始化
pkg/logger                       结构化日志初始化
pkg/utils                        小型通用工具函数
```

## 启动流程

应用当前的启动顺序如下：

1. 读取环境变量配置
2. 初始化结构化日志
3. 连接 Redis
4. 通过 Redis 注册雪花算法 `workerId`
5. 连接 MySQL 并校验可用性
6. 组装业务服务和 HTTP 路由
7. 启动 Gin HTTP 服务

## 注册链路说明

当前会员注册接口的处理流程如下：

1. 中间件从请求头读取 `authorization` 和 `tenantCode`，写入 Gin 上下文
2. handler 绑定并校验注册请求 JSON
3. service 层对密码做摘要处理，构造会员模型并输出业务日志
4. model 层通过 GORM 执行数据库写入
5. handler 返回脱敏后的安全响应对象，不暴露密码字段

## 环境变量

### HTTP

- `APP_PORT`：HTTP 监听地址，默认 `:8080`
- `LOG_LEVEL`：日志级别，可选 `debug`、`info`、`warn`、`error`，默认 `info`

### MySQL

- `MYSQL_DSN`：MySQL 连接串

### Redis

- `REDIS_ADDR`：Redis 地址，默认 `127.0.0.1:6379`
- `REDIS_PASSWORD`：Redis 密码
- `REDIS_DB`：Redis 数据库编号
- `REDIS_MASTER_NAME`：Redis Sentinel 模式下的主节点名称，可选
- `REDIS_WORKER_ID_MIN`：workerId 最小值，默认 `0`
- `REDIS_WORKER_ID_MAX`：workerId 最大值，默认 `1023`
- `REDIS_WORKER_ID_LIFE_SECONDS`：workerId 租期秒数，默认 `15`

### 安全相关

- `PWD_KEY`：会员密码摘要使用的密钥

## 日志说明

项目使用 `log/slog` 输出 JSON 结构化日志。

当前已覆盖的日志范围包括：

- 启动阶段和基础设施初始化
- HTTP 访问日志
- 注册请求处理链路
- 数据落库成功与失败

像密码这类敏感信息在日志中会先脱敏后再输出。

## 当前功能范围

已实现：

- 健康检查
- 会员注册

已在结构中预留、后续可继续扩展的业务方向：

- 登录
- 上分 / 下分
- 场馆列表
- 游戏列表
- 进入游戏
- 下注回调
- 派彩回调
- 注单记录
