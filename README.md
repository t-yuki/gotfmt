gotfmt: Go stack trace formatter library and utilities
======================================================

This is a simple parser library and utilities for [GOTRACEBACK output](http://golang.org/pkg/runtime/) 
Also there is Web UI with demo: http://gotfmt.appspot.com/

Installation
------------

Just type the following to install the program and its dependencies:

    $ go get -u github.com/t-yuki/gotfmt/...

Usage
-----

`gotfmt` formats stack trace texts pretty:

    $ go test -timeout=1us |& gotfmt -f=notest

`gotfmt test` runs `go test` and formats its stack trace texts:

    $ gotfmt test -timeout=1us

`gotfmt -http=:6060` runs Web UI:

    $ gotfmt -http=:6060

Other options:

```
Usage of gotfmt:
  -f="": stack trace filters by comma-separated list
        trimstd:  exclude GOROOT function calls but leave one
        nostd:    exclude GOROOT function calls completely
        notest:   exclude testing function calls
        top:      remove lower function calls
  -h=false: show this help
  -http="": HTTP service address (e.g., ':6060')
  -t="text": output format
        text: pretty formatted text format
        qfix: vim quickfix output format with errorformat: '%f:%l:\ %m'. you should use with 'nostd,notest,top' filters
        json: JSON format
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
