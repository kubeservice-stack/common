# golang common library
[![Go Reference](https://pkg.go.dev/badge/github.com/kubeservice-stack/common.svg)](https://pkg.go.dev/github.com/kubeservice-stack/common) [![Build Status](https://github.com/kubeservice-stack/common/actions/workflows/go.yml/badge.svg)](https://github.com/kubeservice-stack/common/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kubeservice-stack/common)](https://goreportcard.com/report/github.com/kubeservice-stack/common) ![Github release](https://img.shields.io/github/v/release/kubeservice-stack/common) [![codecov](https://codecov.io/github/kubeservice-stack/common/branch/main/graph/badge.svg?token=3AX3EHK96Q)](https://codecov.io/github/kubeservice-stack/common)

[Common Library](https://github.com/kubeservice-stack/common/) 是 一个`编程工具包`,用于在 `Golang` 中构建`微服务`（或`单体`）。解决`分布式系统`和`应用程序架构`中的常见问题，可以让业务更加专注于`交付业务价值`。

此工具包`兼容`目前市场上绝大部分的服务框架: [Kite](https://github.com/koding/kite)、[ServiceComb](https://github.com/go-chassis/go-chassis)、[go-kit](https://github.com/go-kit/kit)、[CloudWeGo/KiteX](https://github.com/cloudwego/kitex)、[gin](https://github.com/gin-gonic/gin)、[beego](https://github.com/beego/beego)

## 动机

`Golang` 已成为服务器语言，但它在 `Facebook`、`Uber`、`Netflix` 和 国内`ByteDance`、`DIDI` 等所谓的“现代企业”公司中的使用度很高。但还有许多这些企业都是基于 `JVM 的堆栈来处理他们的业务逻辑，这在很大程度上归功于直接支持他们的`微服务架构`的`Library库`和`生态系统`。

为了达到更高的成功水平，需要一个`全面`的工具包，以实现`大规模的连贯分布式编程`。`Golang Common Library`就是是一组`包package`和`最佳实践`，它为任何规模的组织构建服务提供了一种`全面`、`健壮`和`可信赖`支持

## 目标

构建相当`完整`、`开箱即用`的`Package集合`

- 各package相互对立，可按需使用
- 减少外部版本依赖，自形生态
- 无业务逻辑，全开放基础能力实现
- 没有特定工具或技术的强制要求

## 依赖管理

基于`go mod`支持多golang语言版本编译： `最小golang`版本支持 `1.12`
