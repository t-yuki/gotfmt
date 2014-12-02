gotfmt: Go stack trace formatter library and utilities
======================================================

This is a simple parser library and utilities for [GOTRACEBACK output](http://golang.org/pkg/runtime/) 
Also there is Web UI with demo: http://gotfmt.appspot.com/

Installation
------------

Just type the following to install the program and its dependencies:

    $ go get -u github.com/t-yuki/gotfmt

Additionally, you can install `got` command.

    $ go get -u github.com/t-yuki/gotfmt/cmd/got

Usage
-----

`gotfmt` formats stack trace texts pretty:

    $ go test -timeout=1us |& gotfmt -f=notest

`gotfmt test` runs `go test` and formats its stack trace texts. `got` is an alias command:

    $ gotfmt test -timeout=1us
    $ got -timeout=1us

`gotfmt -http=:6060` runs Web UI:

    $ gotfmt -http=:6060

Other options:

```
gotfmt - Go Test formatter utility
Usage of gotfmt:
  -f="trimstd,notest": stack trace filters by comma-separated list
        trimstd:  exclude GOROOT function calls but leave one
        nostd:    exclude GOROOT function calls completely
        notest:   exclude testing function calls
        top:      remove lower function calls
  -h=false: show this help
  -http="": HTTP service address (e.g., ':6060')
  -n=1: repeat the test N times while it passes
  -np=0: similar to a combination of `-n` and `-p` but increment GOMAXPROCS from 1 for each repeat
  -p=0: set GOMAXPROCS
  -t="column": output format
	raw: as-is and no filtering
	column: column formatted text
        qfix: vim quickfix output format with errorformat: '%f:%l:\ %m'. you should use with 'nostd,notest,top' filters
        json: JSON format
Any other flags or arguments will be passed to `go` command.
```

Tips
---
### Command Alias
Since `gotfmt XXX` while `XXX` is not `test` runs `go XXX` command and does nothing, you can use it as go command alias.

```
alias go="gotfmt" # `go test` will runs `go test` and formats its output, other options bypasses gotfmt command.
```

### VIM
`qfix` output format generates quickfix list.

    # .vimrc: auto FileType go set errorformat=%f:%l:\ %m
    $ go test -timeout=3s |& gotfmt -f=nostd,notest,top -t=qfix | vim - -c :cb! -c :copen

Web UI Demo
----
The below is web-based analyzer.
It is running `gotfmt -http=:80` command.

http://gotfmt.appspot.com/

Authors
-------

* [Yukinari Toyota (t-yuki)](https://github.com/t-yuki)
