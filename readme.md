myapp/
├── api/                    # API 接口定义文件（如 OpenAPI/Swagger, gRPC proto 文件）
├── build/                  # 打包与持续集成文件（如 Dockerfile, CI/CD 脚本）
├── cmd/                    # 应用程序的入口目录（每个子目录对应一个可执行文件）
│   └── server/
│       └── main.go         # Web 服务的主入口
├── configs/                # 配置文件模板或默认配置（如 config.yaml, env 文件）
├── deployments/            # 部署配置和模板（如 Kubernetes YAML 模板, docker-compose）
├── docs/                   # 设计文档、API 文档
├── internal/               # 私有应用程序和库代码（Go 编译器会强制确保外部无法 import 该目录下的包）
│   ├── handler/            # Web 请求处理层（Controller）
│   ├── middleware/         # Web 中间件（如 JWT 鉴权、跨域处理）
│   ├── model/              # 数据库实体定义及操作
│   └── service/            # 业务逻辑层
├── pkg/                    # 公共库代码（允许外部项目引用的工具或基础库）
├── scripts/                # 用于构建、安装、分析等的脚本
├── test/                   # 外部测试代码和测试数据
├── web/                    # Web 静态资源（前端 SPA 源码、HTML、CSS、JS 等）
├── go.mod                  # 依赖管理文件
└── go.sum                  # 依赖校验和文件



核心目录职责解析

- **`cmd/`**: 存放项目的主要入口。例如，`cmd/server/main.go` 负责启动 Web 服务，`cmd/migration/main.go` 负责数据库迁移脚本。
- **`internal/`**: 存放私有代码。**非常关键**，放在此目录下的包只能被当前项目内部代码引用，Go 编译器会严格限制外部引用，确保业务代码的安全性和隔离性。
- **`pkg/`**: 存放可复用的公共库。如果你的 Web 项目包含独立的工具包（如密码哈希、通用 HTTP 响应封装）且计划开源或给其他项目调用，可放在此目录。
- **`api/`**: 集中存放接口相关的文件。例如 OpenAPI (Swagger) 规范或 `protobuf` 文件，便于前后端联调及统一规范。
- **`web/`**: 专门用于放置前端资源（如 React/Vue 源码、静态页面或打包后的 dist 目录），将前后端代码分离整合在同一个仓库时使用。
- **`configs/`**: 存放配置文件，例如 `config.yaml` 或数据库连接串的配置信息。注意不要将生产环境的敏感密码直接提交到版本控制中。 [[1](https://tonybai.com/2023/10/05/the-official-guide-of-organizing-go-project/), [2](https://www.tizi365.com/topic/11736.html), [3](https://www.cnblogs.com/pingyeaa/p/19058303)]

常用命名规范小贴士

- **单数优先**：目录名称推荐全部使用单数形式，例如 `handler`、`model`、`service`。
- **全小写**：目录和包名应尽量简短，使用全小写字母，避免使用下划线 `_` 或驼峰命名（如使用 `userapi` 而非 `user_api` 或 `userApi`）。





功能列表

1、注册

2、登陆

3、上分

4、下分

5、场馆列表

6、游戏列表

7、进入游戏

8、下注回调

9、派彩回调

10、注单记录

