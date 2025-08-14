
# go-healthcheck

> Golang library that helps creating Healthchecks endpoints that follow the [IETF RFC Health Check](https://tools.ietf.org/id/draft-inadarei-api-health-check-02.html) Response Format for HTTP APIs specification.

![Go version](https://img.shields.io/github/go-mod/go-version/brpaz/go-healthcheck?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/brpaz/go-healthcheck?style=for-the-badge)](https://goreportcard.com/report/github.com/brpaz/go-healthcheck)
[![CI Status](https://github.com/brpaz/go-healthcheck/workflows/CI/badge.svg?style=for-the-badge)](https://github.com/brpaz/go-healthcheck/actions)
[![Coverage Status](https://img.shields.io/codecov/c/github/brpaz/go-healthcheck/master.svg?style=for-the-badge)](https://codecov.io/gh/brpaz/go-healthcheck)


[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](LICENSE)
[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg?style=for-the-badge)](http://commitizen.github.io/cz-cli/)

## Features

This library helps creating Healthchecks endpoints that follows the [IETF RFC Health Check](https://tools.ietf.org/id/draft-inadarei-api-health-check-02.html) Response Format for HTTP APIs specification.

It´s heavily inspired by [health-go](https://github.com/nelkinda/health-go) but the checks are setup in a different way. It also doesnt include any HTTP handler by default. It´s up to you to use the healthcheck library to build the formatted healthcheck response and then adapt to your handler of choice.

It includes the following Healthchecks by default:

* Sysinfo (Uptime, Memory Usage, Load Average, etc)
* Database
* Url
* TCP

## Usage

```go
package main

import (
    "github.com/brpaz/go-healthcheck"
    "github.com/brpaz/go-healthcheck/checks"
)

func main() {
    health := healthcheck.New("myservice", "Some Test service", "1.0.0", "1.0.0-SNAPSHOT")
    health.AddCheckProvider(checks.NewSysInfoChecker())

    result := health.Get()

    // TODO use the result in your HTTP handler to send the response to the health endpoint.
}
```

For instructions how to use the specific checks provided in this package, please see [this](docs/checks.md).


## Creating new checks.

It´s very simple to create a new check. Just create a struct that implements the `Check provider` interface and register it in the healthcheck struct.
You can see examples in the [checks](checks) directory of this project.

## Run tests

```sh
make tests
```


## 🤝 Contributing

Contributions, issues and feature requests are welcome!

## Author

👤 **Bruno Paz**

* Website: [https://github.com/brpaz](https://github.com/brpaz)
* Github: [@brpaz](https://github.com/brpaz)

## 📝 License

Copyright [Bruno Paz](https://github.com/brpaz).

This project is [MIT](LICENSE) licensed.


