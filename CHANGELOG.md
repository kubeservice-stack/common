# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased
- No changes yet.

## [v1.0.2](https://github.com/kubeservice-stack/common/compare/v1.0.1...v1.0.2) - 2023-01-05

- [Feat] Add simple cache & arc cache mode

## [v1.0.1](https://github.com/kubeservice-stack/common/compare/v1.0.0...v1.0.1) - 2023-01-05

- [Feat] Add stream reader & writer package 
- [Feat] Add Bitmap reader & writer package 

## [v1.0.0](https://github.com/kubeservice-stack/common/releases/tag/v1.0.0) - 2022-11-22

### Added
#### Cache 内存缓存

- 支持`lru`、`lfu`和`随机`淘汰算法，同时也支持时间过期淘汰；
- 支持`event`回调： load加载回调、加载成功回调、对象淘汰回调；
- 支持全局内存最大值设置（全局内存） 或 全局最大block（对象数）对象设置；
- 支持 metircs 统计 key 的命中率情况分析；
- 支持 Cache 内部数据 匹配查询（hash）， 也支持近似值匹配查询（多阶二分查找）；
- 支持 内存无GC 高性能模式；

#### Config 配置管理

- 支持 `json`、`ini` 和 环境变量 配置
- 支持 默认值 设置与限制
- 支持 内存使用 配置`cli`输出

#### Connect obejct Pool 连接池管理

- 支持 `连接对象 Connect Obeject`、`全局公用对象 Global Obeject`、`数据对象 Data Objetc` 管理
- 支持 对象 `保活`、`最小可用`、`最小空闲` 管理
- 支持 `event`回调： connect前、clear前、销毁前 回调；
- 支持单机百万TPS管理

#### Discovery 注册中心 和 配置中心

- 支持etcd注册、healthcheck、configure下载、心跳保活

#### errno 通用错误库集合

- 支持 多个 errors 压栈传递
- 支持 errors 全局统一管理 和 标准序列号
- 支持 自定义Error 管理

#### logger 日志库（基于 uber zap on logrus）

- 完全兼容golang标准库日志模块，拥有六种日志级别
- 支持扩展的Hook机制
- 支持日志输出格式：JSONFormatter和TextFormatter
- 支持日志轮转：按log大小 和 按日期
- 支持日志清理：按log总大小 和 按日期
- 支持最高性能

#### metrics 收集

- 支持 `prometheus` 数据接口
- 支持 `opentelemetry` 数据接口
- 支持 atomic 内存态 统计

#### queue 队列

- 支持`RingQueue`、`PriorityQueue` 内存队列管理
- 支持无锁高性能队列使用


#### ratelimiter 限流包

- 支持 令牌桶、滑动时间窗口 限流模式
- 支持对对象限流： 可用于全局 和 部分sub对象限流 的 生成、处理 和 销毁

#### temporary 大文件IO交互库

- 支持超过一定大小、将临时缓冲转成临时文件
- 支持异步写源信息：io.Reader转换成io.ReadSeeker
- 支持异步读取源信息后关闭

#### workpool 任务池

- 支持 任务 `保活`、`最小可用`、`最小空闲` 管理
- 支持 `event`回调： 执行前、执行完成后、销毁前回调；
- 支持 `goroutine` 链方式管理；整个栈内存支持全局设置
- 支持 异常处理 `failover` 处理

#### utils 通用库

- 支持`Map、Strings、Time`等基础数据转换
- 支持 `version`版本号比较
- 支持`数据diff`

### Fixed
- Unittest、staticcheck、lint errors
