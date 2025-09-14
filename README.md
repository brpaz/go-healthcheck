# Go Healthcheck

> A Golang library that provides Healthchecks for your Go Application.


<div align="center">

![Go version](https://img.shields.io/github/go-mod/go-version/brpaz/go-healthcheck?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/brpaz/go-healthcheck/v2?style=for-the-badge)](https://goreportcard.com/report/github.com/brpaz/go-healthcheck/v2)
[![CI status](https://img.shields.io/github/actions/workflow/status/brpaz/go-healthcheck/ci.yml?style=for-the-badge)](https://github.com/brpaz/go-healthcheck/actions/workflows/ci.yml)
[![Coverage Status](https://img.shields.io/codecov/c/github/brpaz/go-healthcheck?style=for-the-badge)](https://codecov.io/gh/brpaz/go-healthcheck)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Documentation](https://img.shields.io/badge/Documentation-Documentation?style=for-the-badge&logo=mdbook&label=Read&color=%23ccc)](https://brpaz.github.io/go-healthcheck)

</div>

## 🎯 Features

- Built-in Healthchecks that covers the monitoring of the most common use cases, like databases, HTTP endpoints, Redis, Disk space, and more. Check the [list of available checks](https://brpaz.github.io/go-healthcheck/checks/).
- Built-in HTTP handler compatible with the native `http` package that constructs the Healthcheck endpoint following the [RFC Healthcheck](https://inadarei.github.io/rfc-healthcheck/) specification.
- Implement your own custom Healthchecks easily by implementing a simple interface.
- No external dependencies.

## 🚀 Getting Started

### Installation

```bash
go get -u github.com/brpaz/go-healthcheck/v2
```

### Basic Usage

To use the library, simply create a new instance of the Healthcheck service, add your desired checks, and expose the healthcheck endpoint using the provided HTTP handler.

```go
package main

import (
  "net/http"

  "github.com/brpaz/go-healthcheck/v2"
  "github.com/brpaz/go-healthcheck/v2/checks/mockcheck"
)

func main() {
  // Create your healthchecks.
  check1 := mockcheck.New(
      mockcheck.WithName("my-check"),
      mockcheck.WithStatus(checks.StatusPass),
  )

  hc := healthcheck.New(
    healthcheck.WithServiceID("my-service"),
    healthcheck.WithDescription("My Service"),
    healthcheck.WithVersion("1.0.0"),
    healthcheck.WithReleaseID("1.0.0-SNAPSHOT"),
    healthcheck.WithChecks(mycheck),
  )

  http.HandleFunc("/health", healthcheck.HealthHandler(hc))
  http.ListenAndServe(":8080", nil)
}
```

For more information about this package and how to use the provided checks, refer to the [Documentation](https://brpaz.github.io/go-healthcheck).

## 🤝 Contributing

All contributions are welcome! Please read the [CONTRIBUTING](CONTRIBUTING.md) file for details on how to contribute.

## 🫶 Support

If you find this project helpful and would like to support its development, there are a few ways you can contribute:

[![Sponsor me on GitHub](https://img.shields.io/badge/Sponsor-%E2%9D%A4-%23db61a2.svg?&logo=github&logoColor=red&&style=for-the-badge&labelColor=white)](https://github.com/sponsors/brpaz)

<a href="https://www.buymeacoffee.com/Z1Bu6asGV" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>

## 🧑‍🦱 Contacts

-  **Email** - [oss@brunopaz.dev](oss@brunopaz.dev)
-  **Source code**: [https://github.com/brpaz/go-healthcheck](https://github.com/brpaz/go-healthcheck)

## 🗒️ License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Reference

- [Health Check Response Format for HTTP APIs](https://inadarei.github.io/rfc-healthcheck/)
- [health-go](https://github.com/hellofresh/health-go)
