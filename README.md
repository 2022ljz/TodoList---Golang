# TodoList（Golang + Gin + Gorm）

一个基于 **Gin** + **Gorm** + **MySQL** 的待办事项（Todo List）示例项目。

项目提供：
- 一个前端页面（`/`）用于展示和管理待办事项。
- 一组 RESTful API（`/v1/todo`）用于待办事项的增删改查。

## 项目结构

```text
.
├── main.go                  # 程序入口，路由注册、数据库初始化、自动迁移
├── controller/
│   └── controller.go        # HTTP 处理函数（控制层）
├── dao/
│   └── mysql.go             # MySQL 连接初始化
├── models/
│   └── todo.go              # Todo 数据模型
├── templates/
│   └── index.html           # 前端页面模板
└── static/                  # 前端静态资源（JS/CSS/字体）
```

## 技术栈

- Go
- Gin
- Gorm
- MySQL

## 功能说明

### 页面路由

- `GET /`：返回 `templates/index.html` 页面。

### API 路由（v1）

基础路径：`/v1`

- `POST /todo`：创建待办事项
- `GET /todo`：获取全部待办事项
- `PUT /todo/:id`：更新指定待办事项
- `DELETE /todo/:id`：删除指定待办事项

## 数据模型

`Todo` 结构如下：

- `id`：整数主键（`int`）
- `title`：待办标题（`string`）
- `status`：完成状态（`bool`）

程序启动时会执行自动迁移：

- `AutoMigrate(&models.Todo{})`

## 运行前准备

### 1) 安装依赖

确保本机已安装：
- Go（建议与 `go.mod` 中版本兼容）
- MySQL

### 2) 创建数据库

根据代码中的 DSN，需要在本地 MySQL 创建数据库：

- 数据库名：`projlist`
- 连接地址：`127.0.0.1:3306`
- 用户名：`root`
- 密码：`123456`

> 当前 DSN 写死在 `dao/mysql.go` 中：
> `root:123456@tcp(127.0.0.1:3306)/projlist?charset=utf8mb4&parseTime=True&loc=Local`

可参考 SQL：

```sql
CREATE DATABASE IF NOT EXISTS projlist DEFAULT CHARSET utf8mb4;
```

## 启动项目

在项目根目录执行：

```bash
go run main.go
```

启动成功后，服务监听在：

- `http://127.0.0.1:8090`

可直接访问页面：

- `http://127.0.0.1:8090/`

## API 使用示例

### 1) 创建 Todo

```bash
curl -X POST http://127.0.0.1:8090/v1/todo \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 Gin","status":false}'
```

### 2) 获取列表

```bash
curl http://127.0.0.1:8090/v1/todo
```

### 3) 更新 Todo

```bash
curl -X PUT http://127.0.0.1:8090/v1/todo/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 Gin + Gorm","status":true}'
```

### 4) 删除 Todo

```bash
curl -X DELETE http://127.0.0.1:8090/v1/todo/1
```

## 注意事项

1. 当前数据库连接配置为硬编码，生产环境建议改为环境变量。
2. 项目中未单独拆分 service 层，控制器直接操作 model/dao，适合作为入门示例。
3. 若启动时报数据库连接错误，请先确认 MySQL 服务是否启动、用户名密码和数据库名是否一致。

## 后续可优化方向

- 使用 `.env` 管理配置（端口、DSN、日志级别）。
- 增加 service 层，解耦业务逻辑。
- 增加参数校验与统一错误响应结构。
- 编写单元测试和接口测试。
- 使用 Docker / docker-compose 提供一键启动。
