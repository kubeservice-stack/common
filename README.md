# golang common library

[![Go Reference](https://pkg.go.dev/badge/github.com/kubeservice-stack/common.svg)](https://pkg.go.dev/github.com/kubeservice-stack/common) [![Build Status](https://github.com/kubeservice-stack/common/actions/workflows/go.yml/badge.svg)](https://github.com/kubeservice-stack/common/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kubeservice-stack/common)](https://goreportcard.com/report/github.com/kubeservice-stack/common) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/96ffd82a42d7484992d015930fd79f76)](https://app.codacy.com/gh/kubeservice-stack/common/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade) [![Github release](https://img.shields.io/github/v/release/kubeservice-stack/common.svg)](https://github.com/kubeservice-stack/common/releases) [![codecov](https://codecov.io/github/kubeservice-stack/common/branch/main/graph/badge.svg?token=3AX3EHK96Q)](https://codecov.io/github/kubeservice-stack/common) [![Apache-2.0 license](https://img.shields.io/github/license/kubeservice-stack/common)](https://github.com/kubeservice-stack/common/blob/main/LICENSE)
[![Last Commit](https://img.shields.io/github/last-commit/kubeservice-stack/common)](https://github.com/kubeservice-stack/common)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon?ref=badge_shield)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#uncategorized)
[![Codacy Coverage Reporter](https://github.com/kubeservice-stack/common/actions/workflows/codacy-coverage-reporter.yaml/badge.svg?branch=main)](https://github.com/kubeservice-stack/common/actions/workflows/codacy-coverage-reporter.yaml)

[EN README](README.EN.md)

[Common Library](https://github.com/kubeservice-stack/common/) 是一个 Go 语言 `编程工具包`，用于构建`微服务`（或`单体应用`）。解决`分布式系统`和`应用程序架构`中的常见问题，让业务更专注于`交付业务价值`。

此工具包`兼容`主流服务框架：[Kite](https://github.com/koding/kite)、[ServiceComb](https://github.com/go-chassis/go-chassis)、[go-kit](https://github.com/go-kit/kit)、[CloudWeGo/Kitex](https://github.com/cloudwego/kitex)、[gin](https://github.com/gin-gonic/gin)、[beego](https://github.com/beego/beego)

## 动机

Go 已成为主流服务器语言，在 Facebook、Uber、Netflix 以及字节跳动、滴滴等"现代企业"中广泛使用。但仍有大量企业基于 JVM 技术栈处理业务逻辑，这很大程度上得益于 JVM 成熟的`微服务库`和`生态系统`。

为了帮助 Go 开发者达到同样的工程效能，我们构建了这个`全面`的工具包。`Golang Common Library`是一组`包(package)`和`最佳实践`，为任何规模的组织提供`全面`、`健壮`、`可信赖`的微服务开发支持。

## 目标

- 各 package 相互独立，可按需引入
- 减少外部版本依赖，自成生态
- 无业务逻辑，完全开放基础能力
- 不强制绑定特定工具或技术

## 快速开始

```bash
go get github.com/kubeservice-stack/common
```

## 依赖管理

基于 `go mod`，支持多 Go 版本编译，最低支持 `Go 1.22`。

## 包概览

### 缓存 (pkg/cache)

内存缓存模块，支持多种淘汰算法和回调机制。

- **淘汰算法**: LRU（最近最少使用）、LFU（最不常用）、FIFO（先进先出）、ARC（自适应替换）、Simple（简单模式）
- **回调事件**: load 加载回调、淘汰回调、清除回调
- **容量控制**: 支持全局内存最大值限制和最大条目数限制
- **指标统计**: 支持命中率等 metrics 统计
- **查询模式**: 精确匹配（hash）和近似值匹配（多阶二分查找）
- **无 GC 模式**: 高性能无 GC 模式

```go
cache := cache.NewLRUCache(&cache.Setting{...})
cache.Set("key", "value")
val, err := cache.Get("key")
```

### 配置管理 (pkg/config)

统一配置加载，支持 TOML 文件和环境变量。

- 支持 TOML 格式配置文件
- 支持环境变量自动覆盖（通过 `env` 标签）
- 内置配置项：Logging、Metrics、Discovery、Gin、RateLimit、Temporary、Database
- 支持配置默认值设置

```go
config.GlobalCfg.Logging.Level = "debug"
config.LoadGlobalConfig("config.toml")
```

### 服务发现 (pkg/discovery)

基于 etcd 的服务注册与发现。

- 服务注册与健康检查（心跳保活）
- Leader 选举
- Key-Value 存储（Get、List、Put、Delete、Batch）
- Watch 监听（支持单个 key 和前缀匹配）
- 事务支持

```go
factory := discovery.NewDiscoveryFactory("my-service")
disc, _ := factory.CreateDiscovery(config.Discovery{...})
data, _ := disc.Get(ctx, "/services/my-app")
```

### 序列化编解码 (pkg/codec)

统一编解码接口，支持多种序列化格式。

- **MCPack**: 类似 JSON 的二进制编解码格式
- **MsgPack**: 基于 msgpack 协议的序列化

```go
codec := codec.PluginInstance(codec.MSGPACK)
data, _ := codec.Marshal(obj)
codec.Unmarshal(data, &obj)
```

### 数据库 ORM (pkg/orm)

基于 GORM 的数据库连接管理，支持多种数据库和缓存。

- **数据库驱动**: MySQL、PostgreSQL、SQLite3
- **ORM 缓存**: Redis 缓存、内存缓存
- **缓存模式**: 禁用缓存、仅主键缓存、仅查询缓存、全量缓存
- **可追溯**: 内置 OpenTelemetry tracing 插件
- **缓存策略**: 缓存穿透保护、异步写入、自定义 TTL

```go
db, _ := orm.NewDBConn(config.DBConfig{
    DBType:   orm.MYSQL,
    Host:     "localhost",
    Port:     "3306",
    Database: "mydb",
})
defer db.Close()
```

### 错误码管理 (pkg/errno)

统一的错误码定义和管理。

- 支持状态码 + 消息的标准化错误结构
- 预定义常见错误码常量

```go
err := errno.New(404, "resource not found")
```

### 日志模块 (pkg/logger)

基于 Uber Zap 的高性能日志。

- 兼容标准库日志接口，支持 Debug/Info/Warn/Error 级别
- 支持 JSON 和 Text 两种输出格式
- 按模块和角色分类日志（HTTP、Crash 等）
- 终端自动检测（彩色输出）
- 基于 lumberjack 的日志轮转（按大小、按日期）

```go
logger := logger.GetLogger("my-module", "my-role")
logger.Info("operation completed", logger.Int64("duration", 100))
```

### 指标采集 (pkg/metrics)

支持 Prometheus 和 OpenTelemetry 指标采集。

- Prometheus HTTP 接口暴露指标
- OpenTelemetry metrics 集成
- 基于 tally 的内存统计
- Stats 统计工具

### 连接池 (pkg/connpool)

通用连接对象池管理。

- 最大活跃连接数限制
- 空闲连接保留数控制
- 空闲超时自动回收
- 等待超时机制（支持永久等待、不等待）
- 连接创建/销毁/清理回调

```go
pool := connpool.NewConnectionPool(
    100, 10, 5*time.Minute, 30,
    connectFunc, disconnectFunc, clearFunc,
)
conn, _ := pool.Pop()
defer pool.Push(conn)
```

### 有向无环图 (pkg/dag)

DAG（有向无环图）数据结构。

- 顶点管理（增删查）
- 边管理（增删）
- 获取汇点（Sink）、源点（Source）
- 前驱/后继查询
- 基于 OrderMap 的有序遍历

### 分布式锁 (pkg/lock)

轻量级锁接口。

- **文件锁**: 支持 Unix 和 Windows
- **内存锁**: 进程内互斥锁

### 限流器 (pkg/ratelimiter)

多维度限流。

- **令牌桶算法**: 支持 QPS 和 burst 控制
- **动态调整**: 支持运行时更新限流参数
- **多对象限流**: 可对不同 name 独立限流

```go
limiter.TryAccept("api", 100, 50) // QPS=100, burst=50
```

### 任务调度 (pkg/schedule)

轻量级定时任务调度。

- 基于间隔的任务调度
- 任务标签管理
- 支持按函数名、引用、标签移除任务
- 全局默认调度器和独立调度器
- 支持 RunPending、RunAll 等运行模式

```go
schedule.Every(60).Do(myFunc).Tag("cleanup")
schedule.Start()
```

### 任务池 (pkg/workpool)

Goroutine 工作池。

- 最大 worker 数量控制
- 空闲 goroutine 自动回收
- 原子指标统计（创建数、活跃数、完成任务数）
- 支持同步提交（SubmitAndWait）和异步提交（Submit）
- 优雅关闭，处理剩余任务

```go
pool := workpool.NewWorkerPool("my-pool", 10, 5*time.Minute)
defer pool.Stop()
pool.Submit(func() { /* do work */ })
```

### 队列 (pkg/queue)

内存队列实现。

- **RingQueue**: 基于环形缓冲区的无锁队列
- **PriorityQueue**: 优先级队列
- 支持单条和批量 Pop

```go
q := queue.NewRingQueue(1024)
q.Push("message")
msg, ok := q.Pop()
```

### 集合 (pkg/sets)

类型安全的集合操作。

- **基础类型**: ByteSet、IntSet、Int32Set、Int64Set、StringSet
- **集合操作**: 并集、交集、差集、对称差集
- **子集/超集**判断
- 线程安全选项

### 有序集合 (pkg/orders)

保持插入顺序的集合数据结构。

- **OrderedMap**: 有序键值对映射
- **OrderedSet**: 有序元素集合
- 提供 Size、Empty、Values、String 等接口

### 流式二进制 (pkg/stream)

二进制流读写工具。

- **Reader**: 支持 Varint、固定长度整数、Byte、Slice 等读取操作
- **Writer**: 支持对应的写入操作
- 支持 Seek 和 Reset
- 内置编码格式支持

### 存储引擎 (pkg/storage)

轻量级时序数据存储。

- 基于分区（Partition）的数据组织
- 标签（Label）过滤查询
- 时间范围查询（Select by name, labels, start, end）
- 内存分区和分区列表管理

### 临时缓冲 (pkg/temporary)

将 `io.Reader` 转换为 `io.ReadSeeker`。

- 超过阈值自动从内存缓冲切换到临时文件
- 支持同步和异步读取模式
- 支持 ReadCloser 自动关闭
- 实现 io.Reader、io.Seeker、io.Closer 接口

### 位图操作 (pkg/bit)

二进制位操作。

- **Reader**: 从字节流中按位读取
- **Writer**: 向字节流中按位写入

### 追踪 (pkg/tracing)

分布式追踪（OpenTracing + OpenTelemetry）。

- **OpenTracing 桥接**: 兼容 OpenTracing API
- **Jaeger**: Jaeger tracer 集成和远程采样
- **OTLP**: OpenTelemetry Protocol 导出（gRPC/HTTP）
- **工具函数**: StartSpan、DoInSpan、DoWithSpan、DoInSpanWithErr
- **上下文传递**: CopyTraceContext、ContextWithTracer
- **强制采样**: 通过 X-Force-Tracing header 强制采样

### 工具函数 (pkg/utils)

常用工具函数集合。

- **字符串操作**: 格式化、转换等
- **Map 操作**: 类型安全的键值访问
- **Slice 操作**: 切片工具
- **时间操作**: Duration 解析、定时器工具
- **时间戳**: 时间戳处理
- **版本号比较**: 语义化版本比较
- **文件操作**: 文件工具
- **数据对比**: 对象 diff

### 比特掩码 (pkg/bufioutil)

IO 缓冲区优化工具。

### 配置加载器 (pkg/config/loader)

TOML 配置文件加载。

### 自定义 TOML (pkg/config/ltoml)

增强的 TOML 配置解析。

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon?ref=badge_large)

[Apache-2.0](LICENSE)
