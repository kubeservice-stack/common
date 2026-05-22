# golang common library

[![Go Reference](https://pkg.go.dev/badge/github.com/kubeservice-stack/common.svg)](https://pkg.go.dev/github.com/kubeservice-stack/common) [![Build Status](https://github.com/kubeservice-stack/common/actions/workflows/go.yml/badge.svg)](https://github.com/kubeservice-stack/common/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kubeservice-stack/common)](https://goreportcard.com/report/github.com/kubeservice-stack/common) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/96ffd82a42d7484992d015930fd79f76)](https://app.codacy.com/gh/kubeservice-stack/common/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade) [![Github release](https://img.shields.io/github/v/release/kubeservice-stack/common.svg)](https://github.com/kubeservice-stack/common/releases) [![codecov](https://codecov.io/github/kubeservice-stack/common/branch/main/graph/badge.svg?token=3AX3EHK96Q)](https://codecov.io/github/kubeservice-stack/common) [![Apache-2.0 license](https://img.shields.io/github/license/kubeservice-stack/common)](https://github.com/kubeservice-stack/common/blob/main/LICENSE)
[![Last Commit](https://img.shields.io/github/last-commit/kubeservice-stack/common)](https://github.com/kubeservice-stack/common)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon?ref=badge_shield)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#uncategorized)
[![Codacy Coverage Reporter](https://github.com/kubeservice-stack/common/actions/workflows/codacy-coverage-reporter.yaml/badge.svg?branch=main)](https://github.com/kubeservice-stack/common/actions/workflows/codacy-coverage-reporter.yaml)

[中文文档](README.md)

[Common Library](https://github.com/kubeservice-stack/common/) is a `programming toolkit` for building `microservices` (or `monolithic services`) in `Golang`. It solves common problems in `distributed systems` and `application architecture`, allowing businesses to focus on `delivering value`.

This toolkit is `compatible` with most popular service frameworks: [Kite](https://github.com/koding/kite)、[ServiceComb](https://github.com/go-chassis/go-chassis)、[go-kit](https://github.com/go-kit/kit)、[CloudWeGo/Kitex](https://github.com/cloudwego/kitex)、[gin](https://github.com/gin-gonic/gin)、[beego](https://github.com/beego/beego)

## Motivation

Golang has become a mainstream server-side language, widely used in "modern enterprise" companies like Facebook, Uber, Netflix, and domestic ByteDance, Didi. However, many enterprises still rely on JVM-based stacks for their business logic, largely thanks to the mature `Libraries` and `Ecosystems` that directly support their `Microservice Architectures`.

To achieve the same level of engineering success requires a `comprehensive` toolkit for `distributed programming at scale`. The `Golang Common Library` is a set of `packages` and `best practices` that provide `comprehensive`, `robust`, and `trusted` support for building services for organizations of any size.

## Goal

- Each package is independent and can be imported as needed
- Reduce external version dependencies, self-contained ecosystem
- No business logic, fully open basic capabilities
- No mandatory requirements for specific tools or techniques

## Quick Start

```bash
go get github.com/kubeservice-stack/common
```

## Dependency Management

Based on `go.mod`, supports multiple Go versions. Minimum supported Go version: `1.22`.

## Performance & Security

The project has undergone the following optimizations:

- **Concurrency Safety**: Fixed ratelimiter TOCTOU race condition, connpool infinite loop, workpool busy-wait
- **Performance Optimization**: Cache Len/HasKey/Keys operations optimized from O(n) to O(1), reduced lock contention
- **CI/CD**: Added race detection, dependency caching for faster builds, weekly dependency update strategy

## Packages

### Cache (pkg/cache)

In-memory caching with multiple eviction algorithms and callback mechanisms.

- **Eviction Algorithms**: LRU, LFU, FIFO, ARC (Adjustable Replacement Cache), Simple
- **Callbacks**: Load callback, eviction callback, purge callback
- **Capacity Control**: Global memory limit and max entry count
- **Metrics**: Hit rate statistics
- **Query Modes**: Exact match (hash) and approximate match (multi-stage binary search)
- **GC-free Mode**: High-performance mode without GC overhead

```go
cache := cache.NewLRUCache(&cache.Setting{...})
cache.Set("key", "value")
val, err := cache.Get("key")
```

### Config (pkg/config)

Unified configuration loading from TOML files and environment variables.

- TOML configuration file support
- Environment variable auto-override (via `env` struct tags)
- Built-in config sections: Logging, Metrics, Discovery, Gin, RateLimit, Temporary, Database
- Default value support

```go
config.GlobalCfg.Logging.Level = "debug"
config.LoadGlobalConfig("config.toml")
```

### Discovery (pkg/discovery)

Service registration and discovery based on etcd.

- Service registration and health check (heartbeat)
- Leader election
- Key-Value storage (Get, List, Put, Delete, Batch)
- Watch support (single key and prefix)
- Transaction support

```go
factory := discovery.NewDiscoveryFactory("my-service")
disc, _ := factory.CreateDiscovery(config.Discovery{...})
data, _ := disc.Get(ctx, "/services/my-app")
```

### Codec (pkg/codec)

Unified serialization/Deserialization interface.

- **MCPack**: Binary encode/decode format similar to JSON
- **MsgPack**: msgpack-based serialization

```go
codec := codec.PluginInstance(codec.MSGPACK)
data, _ := codec.Marshal(obj)
codec.Unmarshal(data, &obj)
```

### ORM (pkg/orm)

Database connection management based on GORM with caching.

- **Database Drivers**: MySQL, PostgreSQL, SQLite3
- **ORM Cache**: Redis cache, memory cache
- **Cache Modes**: Disable, OnlyPrimary, OnlySearch, All
- **Tracing**: Built-in OpenTelemetry tracing plugin
- **Cache Strategy**: Cache penetration protection, async write, custom TTL

```go
db, _ := orm.NewDBConn(config.DBConfig{
    DBType:   orm.MYSQL,
    Host:     "localhost",
    Port:     "3306",
    Database: "mydb",
})
defer db.Close()
```

### Errno (pkg/errno)

Standardized error code management.

- Status code + message error structure
- Predefined common error codes

```go
err := errno.New(404, "resource not found")
```

### Logger (pkg/logger)

High-performance logging based on Uber Zap.

- Standard library compatible with Debug/Info/Warn/Error levels
- JSON and Text output formats
- Module and role-based log categorization (HTTP, Crash, etc.)
- Automatic terminal detection (colored output)
- Log rotation via lumberjack (by size and date)

```go
logger := logger.GetLogger("my-module", "my-role")
logger.Info("operation completed", logger.Int64("duration", 100))
```

### Metrics (pkg/metrics)

Prometheus and OpenTelemetry metrics collection.

- Prometheus HTTP metrics endpoint
- OpenTelemetry metrics integration
- Tally-based in-memory statistics
- Stats utility

### Connection Pool (pkg/connpool)

Generic connection object pool.

- Max active connections limit
- Reserved idle connections control
- Idle timeout auto-reclamation
- Wait timeout mechanism (supports forever wait, no wait)
- Connection create/destroy/clear callbacks

```go
pool := connpool.NewConnectionPool(
    100, 10, 5*time.Minute, 30,
    connectFunc, disconnectFunc, clearFunc,
)
conn, _ := pool.Pop()
defer pool.Push(conn)
```

### DAG (pkg/dag)

Directed Acyclic Graph data structure.

- Vertex management (add, delete, get)
- Edge management (add, delete)
- Sink and source vertices discovery
- Predecessor/successor queries
- Ordered traversal via OrderMap

### Lock (pkg/lock)

Lightweight locking interface.

- **File Lock**: Supports Unix and Windows
- **Memory Lock**: In-process mutex

### Rate Limiter (pkg/ratelimiter)

Multi-dimensional rate limiting.

- **Token Bucket Algorithm**: QPS and burst control
- **Dynamic Adjustment**: Runtime parameter updates
- **Per-Object Limiting**: Independent limits per name

```go
limiter.TryAccept("api", 100, 50) // QPS=100, burst=50
```

### Schedule (pkg/schedule)

Lightweight cron-like task scheduling.

- Interval-based task scheduling
- Task tagging
- Remove by function name, reference, or tag
- Global default scheduler and standalone schedulers
- RunPending, RunAll modes

```go
schedule.Every(60).Do(myFunc).Tag("cleanup")
schedule.Start()
```

### Worker Pool (pkg/workpool)

Goroutine worker pool.

- Max worker count control
- Idle goroutine auto-recycling
- Atomic metrics (created, alive, consumed)
- Sync submit (SubmitAndWait) and async submit (Submit)
- Graceful shutdown with remaining task processing

```go
pool := workpool.NewWorkerPool("my-pool", 10, 5*time.Minute)
defer pool.Stop()
pool.Submit(func() { /* do work */ })
```

### Queue (pkg/queue)

In-memory queue implementations.

- **RingQueue**: Lock-free ring buffer queue
- **PriorityQueue**: Priority-based queue
- Single and batch pop support

```go
q := queue.NewRingQueue(1024)
q.Push("message")
msg, ok := q.Pop()
```

### Sets (pkg/sets)

Type-safe set operations.

- **Types**: ByteSet, IntSet, Int32Set, Int64Set, StringSet
- **Operations**: Union, Intersection, Difference, Symmetric Difference
- Subset/Superset checking
- Thread-safety option

### Ordered Collections (pkg/orders)

Insertion-order preserving collections.

- **OrderedMap**: Ordered key-value map
- **OrderedSet**: Ordered element set
- Size, Empty, Values, String interfaces

### Stream (pkg/stream)

Binary stream read/write utilities.

- **Reader**: Varint, fixed-size integer, byte, slice reading
- **Writer**: Corresponding write operations
- Seek and Reset support
- Built-in encoding formats

### Storage (pkg/storage)

Lightweight time-series data storage.

- Partition-based data organization
- Label filtering queries
- Time range queries (Select by name, labels, start, end)
- In-memory partition and partition list management

### Temporary (pkg/temporary)

Converts `io.Reader` to `io.ReadSeeker`.

- Auto-switch from memory buffer to temp file beyond threshold
- Synchronous and asynchronous read modes
- ReadCloser auto-close support
- Implements io.Reader, io.Seeker, io.Closer

### Bitmap (pkg/bit)

Binary bit manipulation.

- **Reader**: Bit-level reading from byte stream
- **Writer**: Bit-level writing to byte stream

### Tracing (pkg/tracing)

Distributed tracing (OpenTracing + OpenTelemetry).

- **OpenTracing Bridge**: OpenTracing API compatibility
- **Jaeger**: Jaeger tracer integration with remote sampling
- **OTLP**: OpenTelemetry Protocol export (gRPC/HTTP)
- **Utilities**: StartSpan, DoInSpan, DoWithSpan, DoInSpanWithErr
- **Context Passing**: CopyTraceContext, ContextWithTracer
- **Force Sampling**: X-Force-Tracing header for forced sampling

### Utils (pkg/utils)

Common utility functions.

- **Strings**: Formatting, conversion
- **Map Operations**: Type-safe key access
- **Slice Operations**: Slice utilities
- **Time**: Duration parsing, timer utilities
- **Timestamp**: Timestamp processing
- **Version Comparison**: Semantic versioning
- **File Operations**: File utilities
- **Data Comparison**: Object diff

### Buffer IO (pkg/bufioutil)

IO buffer optimization utilities.

### Config Loader (pkg/config/loader)

TOML configuration file loader.

### Custom TOML (pkg/config/ltoml)

Enhanced TOML configuration parsing.

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon?ref=badge_large)

[Apache-2.0](LICENSE)
