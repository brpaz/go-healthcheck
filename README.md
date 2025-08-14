
# go-healthcheck

> Golang library that helps creating Healthchecks endpoints that follow the [IETF RFC Health Check](https://tools.ietf.org/id/draft-inadarei-api-health-check-02.html) format.

![Go version](https://img.shields.io/github/go-mod/go-version/brpaz/go-healthcheck?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/brpaz/go-healthcheck?style=for-the-badge)](https://goreportcard.com/report/github.com/brpaz/go-healthcheck)
[![CI Status](https://github.com/brpaz/go-healthcheck/workflows/CI/badge.svg?style=for-the-badge)](https://github.com/brpaz/go-healthcheck/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/brpaz/go-healthcheck/master.svg?style=for-the-badge)](https://codecov.io/gh/brpaz/go-healthcheck)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)
[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg?style=for-the-badge)](http://commitizen.github.io/cz-cli/)

## Features

This library helps creating Healthchecks endpoints that follows the [IETF RFC Health Check](https://tools.ietf.org/id/draft-inadarei-api-health-check-02.html) format for HTTP APIs specification.

It´s heavily inspired by [health-go](https://github.com/nelkinda/health-go) but the checks are setup in a different way.

The following healthchecks are included by default (PRs welcome):

- **HTTP Check** - A check that checks if an http endpoint is reachable and returns a successful status code.
- **TCP Check** - A check that checks if a TCP endpoint is reachable.
- **SysInfo Check** - A check that returns system information, such as CPU and memory usage.
- **Disk Check** - A check that returns disk usage information.
- **DB Check** - A check that checks if can connect to database.
- **Redis Check** - A check that checks if can connect to a Redis instance.
- **Mock Check** - A simple check that returns success or fail, based on the provided argument. Useful for tests.

## Getting started

### Installation

```shell
go get github.com/brpaz/go-healthcheck
```

### Usage

```go
package main

import "net/http"
import "github.com/brpaz/go-healthcheck"
import "github.com/brpaz/go-healthcheck/pkg/checks/mockcheck"

func main() {

    // 1. Declare your checks
    mycheck := mockcheck.New(
        mockcheck.WithName("my-check"),
    )

    // Initialize the healthcheck service
    hc := healthcheck.New(
        healthcheck.WithServiceID("my-service"),
        healthcheck.WithDescription("My Service"),
        healthcheck.WithVersion("1.0.0"),
        healthcheck.WithReleaseID("1.0.0-SNAPSHOT"),
        healthcheck.WithCheck(mycheck),
    )

    // Serve the healthcheck endpoint using the provided handler
    http.HandleFunc("/health", healthcheck.Handler(hc))
    http.ListenAndServe(":8080", nil)
```

> [!TIP]
> For instructions how to use the specific checks provided in this package, please see [this](docs/checks.md).

## Adding your own checks

If you want to build your own custom checks, it´s very simple.

1. Create a new struct that implements the `Check` interface.
2. Register your new check with the healthcheck service.

You can see examples in the [checks](checks) directory of this project.

## 🤝 Contributing

Contributions, issues and feature requests are welcome!

## Author

👤 **Bruno Paz**

* Website: [https://github.com/brpaz](https://github.com/brpaz)
* Github: [@brpaz](https://github.com/brpaz)

## 📝 License

Copyright [Bruno Paz](https://github.com/brpaz).

This project is [MIT](LICENSE) licensed.


