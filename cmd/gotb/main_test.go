package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestConvert_data1(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/data1.txt")
	if err != nil {
		t.Fatal(err)
	}
	out := &bytes.Buffer{}
	convert(bytes.NewBuffer(data), out)
}

func TestConvert_data2(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/data2.txt")
	if err != nil {
		t.Fatal(err)
	}
	out := &bytes.Buffer{}
	convert(bytes.NewBuffer(data), out)
}
