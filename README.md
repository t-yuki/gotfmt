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

    $ go test -timeout=1us |& gotfmt

`gotfmt test` runs `go test` and formats its stack trace texts:

    $ gotfmt test -timeout=1us

`gotfmt -http=:6060` runs Web UI:

    $ gotfmt -http=:6060

Tips
---
Since `gotfmt XXX` while `XXX` is not `test` runs `go XXX` command and does nothing, you can use it as go command alias.

```
alias go="gotfmt" # `go test` will runs `go test` and formats its output, other options bypasses gotfmt command.
```

Alos the command can make quickfix list style:

    # .vimrc: auto FileType go set errorformat=%f:%l:\ %m
    $ go test -timeout=3s |& gotb -q | vim - -c :cb! -c :copen

Web UI Demo
----
The below is web-based analyzer.
It is running `gotfmt -http=:80` command.

http://gotfmt.appspot.com/

Authors
-------

* [Yukinari Toyota (t-yuki)](https://github.com/t-yuki)
