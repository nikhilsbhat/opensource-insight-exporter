# OpenSource Insight prometheus exporter


[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/opensource-insight-exporter)](https://goreportcard.com/report/github.com/nikhilsbhat/opensource-insight-exporter)
[![shields](https://img.shields.io/badge/license-MIT-blue)](https://github.com/nikhilsbhat/opensource-insight-exporter/blob/master/LICENSE)
[![shields](https://godoc.org/github.com/nikhilsbhat/opensource-insight-exporter?status.svg)](https://godoc.org/github.com/nikhilsbhat/opensource-insight-exporter)
[![shields](https://img.shields.io/github/v/tag/nikhilsbhat/opensource-insight-exporter.svg)](https://github.com/nikhilsbhat/opensource-insight-exporter/tags)
[![shields](https://img.shields.io/github/downloads/nikhilsbhat/opensource-insight-exporter/total.svg)](https://github.com/nikhilsbhat/opensource-insight-exporter/releases)


prometheus exporter to get an insight on downloads of opensource projects.

## Introduction

prometheus exporter for identifying download metrics of open-source projects hosted on multiple platforms. It supports retrieving metrics from `github` and `terraform-registry` for now.

This exporter interacts with the github API, terraform-registry API, and other platforms to have the metrics collected in one place.

## Requirements

* [Go](https://golang.org/dl/) 1.19 or above . Installing go can be found [here](https://golang.org/doc/install).
* Basic understanding of prometheus exporter and its golang [client](https://github.com/prometheus/client_golang.git) libraries and [building](https://prometheus.io/docs/guides/go-application/) them.


## Installation

* Recommend installing released versions. Release binaries are available on the [releases](https://github.com/nikhilsbhat/opensource-insight-exporter/releases) page and docker from [here](https://github.com/nikhilsbhat/opensource-insight-exporter/pkgs/container/opensource-insight-exporter).
* Can always build it locally by running `go build` against cloned repo.

#### Docker

```bash
docker pull ghcr.io/nikhilsbhat/opensource-insight-exporter:latest
docker pull ghcr.io/nikhilsbhat/opensource-insight-exporter:<github-release-tag>
```