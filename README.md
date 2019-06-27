zaperr
======

[![CircleCI](https://circleci.com/gh/hori-ryota/zaperr.svg?style=svg)](https://circleci.com/gh/hori-ryota/zaperr)
[![Coverage Status](https://coveralls.io/repos/github/hori-ryota/zaperr/badge.svg?branch=master)](https://coveralls.io/github/hori-ryota/zaperr?branch=master)
[![GoDoc](https://godoc.org/github.com/hori-ryota/zaperr?status.svg)](https://godoc.org/github.com/hori-ryota/zaperr)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Error wrapper for [go.uber.org/zap](https://github.com/uber-go/zap)

## Usage

```go
err = zaperr.Wrap(err, "failed to execute something",
    zap.Int("foo", 1),
    zap.String("bar", "baz"),
)

logger.Info("example", zaperr.ToField(err))

// Output:
// {"level":"info","msg":"example","error":{"foo":1,"bar":"baz","error":"failed to execute something: error","errorVerbose":"omitted..."}}
```
