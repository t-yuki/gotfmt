gotracetools: Go stack trace parser library and utilities
=========================================================

This is a simple parser library and utilities for [GOTRACEBACK output](http://golang.org/pkg/runtime/) 

Installation
------------

Just type the following to install the program and its dependencies:

    $ go get -u github.com/t-yuki/gotracetools/...

Usage
-----

`gotb` formats stack trace texts pretty:

    $ go test -timeout=1us 2>&1 | gotb

Or it can make quickfix list:

    # .vimrc: auto FileType go set errorformat=%f:%l:\ %m
    $ go test -timeout=3s 2>&1 | gotb -q | vim - -c :cb! -c :copen

Authors
-------

* [Yukinari Toyota (t-yuki)](https://github.com/t-yuki)
