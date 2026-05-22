# 配置命名约定

本文定义 `zhp-app` 中环境变量命名和 Go 配置字段命名的统一规则。

## 目标

统一命名主要解决三个问题：

1. 避免同一配置在环境变量、Go 结构体、Docker Compose 中出现多套风格
2. 降低新同学阅读和新增配置时的理解成本
3. 让配置项在代码、脚本、部署文件中可以一眼建立映射关系

## 一、环境变量命名规则

环境变量统一使用：

```text
全大写 + 下划线分隔
```

示例：

- `APP_PORT`
- `LOG_LEVEL`
- `MYSQL_DSN`
- `REDIS_ADDR`
- `REDIS_MASTER_NAME`
- `REDIS_WORKER_ID_MAX`
- `PWD_KEY`

### 细则

1. 环境变量必须全大写。
2. 单词之间使用下划线 `_` 分隔。
3. 前缀优先体现配置归属，例如 `APP_`、`MYSQL_`、`REDIS_`。
4. 缩写在环境变量中保持全大写，例如 `DSN`、`DB`、`ID`、`PWD`。
5. 不使用驼峰，例如不要写 `redisMasterName`、`mysqlDsn`。

## 二、Go 配置字段命名规则

Go 配置结构体字段统一使用：

```text
PascalCase（大驼峰）
```

示例：

- `Port`
- `LogLevel`
- `MySQLDSN`
- `PasswordKey`
- `WorkerID`
- `RedisConfig`
- `MasterName`

### 细则

1. 导出的配置字段使用大驼峰命名。
2. 常见缩写遵循 Go 命名习惯，保留大写形式：
   - `ID`
   - `DB`
   - `DSN`
   - `HTTP`
   - `JSON`
   - `URL`
3. 避免语义不清的缩写字段名：
   - 推荐 `PasswordKey`
   - 不推荐 `PwdKey`
4. `Config` 优于 `Conf`：
   - 推荐 `RedisConfig`
   - 不推荐 `RedisConf`

## 三、环境变量与 Go 字段映射规则

统一采用“环境变量表达部署语义，Go 字段表达代码语义”的方式。

示例映射：

| 环境变量 | Go 字段 |
|---|---|
| `APP_PORT` | `Port` |
| `LOG_LEVEL` | `LogLevel` |
| `MYSQL_DSN` | `MySQLDSN` |
| `PWD_KEY` | `PasswordKey` |
| `REDIS_ADDR` | `RedisConfig.Addr` |
| `REDIS_DB` | `RedisConfig.DB` |
| `REDIS_MASTER_NAME` | `RedisConfig.MasterName` |
| `REDIS_WORKER_ID_MIN` | `RedisConfig.WorkerIDMin` |
| `REDIS_WORKER_ID_MAX` | `RedisConfig.WorkerIDMax` |
| `REDIS_WORKER_ID_LIFE_SECONDS` | `RedisConfig.WorkerIDLifeSeconds` |

## 四、当前项目中的统一结论

当前项目后续新增配置时，统一遵循下面的规则：

### 环境变量

- 使用全大写下划线命名
- 保持部署文件、脚本、README 写法一致

### Go 字段

- 使用大驼峰
- 常见缩写保留大写
- 优先完整语义，不滥用缩写

## 五、推荐与不推荐示例

推荐：

- `MYSQL_DSN` -> `MySQLDSN`
- `PWD_KEY` -> `PasswordKey`
- `REDIS_MASTER_NAME` -> `RedisConfig.MasterName`
- `REDIS_WORKER_ID_MAX` -> `RedisConfig.WorkerIDMax`

不推荐：

- `mysqlDsn`
- `redisMasterName`
- `PwdKey`
- `RedisConf`
- `MySqlDsn`

## 六、执行原则

1. 新增配置时，先确定环境变量名，再确定 Go 字段名。
2. 环境变量遵循部署规范，不为了和 Go 字段“看起来一致”而改成驼峰。
3. Go 字段遵循 Go 代码风格，不直接照抄环境变量全大写写法。
4. 文档、脚本、Docker Compose、代码中的配置命名必须同步更新。
