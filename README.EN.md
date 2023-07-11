# golang common library
[![Go Reference](https://pkg.go.dev/badge/github.com/kubeservice-stack/common.svg)](https://pkg.go.dev/github.com/kubeservice-stack/common) [![Build Status](https://github.com/kubeservice-stack/common/actions/workflows/go.yml/badge.svg)](https://github.com/kubeservice-stack/common/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kubeservice-stack/common)](https://goreportcard.com/report/github.com/kubeservice-stack/common) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/96ffd82a42d7484992d015930fd79f76)](https://app.codacy.com/gh/kubeservice-stack/common/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade) [![Github release](https://img.shields.io/github/v/release/kubeservice-stack/common.svg)](https://github.com/kubeservice-stack/common/releases) [![codecov](https://codecov.io/github/kubeservice-stack/common/branch/main/graph/badge.svg?token=3AX3EHK96Q)](https://codecov.io/github/kubeservice-stack/common) [![Apache-2.0 license](https://img.shields.io/github/license/kubeservice-stack/common)](https://github.com/kubeservice-stack/common/blob/main/LICENSE)
[![Last Commit](https://img.shields.io/github/last-commit/kubeservice-stack/common)](https://github.com/kubeservice-stack/common)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon?ref=badge_shield)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#uncategorized)

[中文版README](README.md)

[Common Library](https://github.com/kubeservice-stack/common/) is a `programming toolkit` for building `microservices` (or `allInOne service`) in `Golang`. Solving common problems in `fractionated systems` and `application architecture` is more cumbersome, allowing more problems to be solved in `traffic payment business value`.

This toolkit is `compatible` with most of the service frameworks currently on the market.
Like: [Kite](https://github.com/koding/kite)、[ServiceComb](https://github.com/go-chassis/go-chassis)、[go-kit](https://github.com/go-kit/kit)、[CloudWeGo/KiteX](https://github.com/cloudwego/kitex)、[gin](https://github.com/gin-gonic/gin)、[beego](https://github.com/beego/beego)

## Motivation

`Golang` has become a mainstream language, but it is heavily used in so-called "modern enterprise" companies such as `Facebook`, `Uber`, `Netflix` and domestic `ByteDance`, `DIDI`. But there are also many of these businesses that handle their business logic on a JVM-based stack, thanks in large part to the `Libraries` and `Ecosystems` that directly support their `Microservice Architectures`.

To reach higher levels of success requires a `comprehensive` toolkit for `coherent distributed programming at scale`. The `Golang Common Library` is a set of `packages` and `best practices` that provide a `comprehensive`, `robust` and `trusted` support for building services for organizations of any size


## Goal

Build a fairly `complete`, `out-of-the-box` `Package Collection`

- Each package is opposite to each other and can be used as needed
- Reduce external version dependencies, self-shaped ecology
- No business logic, fully open basic capabilities
- No mandatory requirements for specific tools or techniques

## Dependency Management

Based on `go.mod`, it supports multi-golang language version compilation: `Minimum golang` version supports `1.12`

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkubeservice-stack%2Fcommon?ref=badge_large)
