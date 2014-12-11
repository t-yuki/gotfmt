package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestConvert_data1(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/data1.txt")
	if err != nil {
		t.Fatal(err)
	}
	testConvert(t, data)
}

func TestConvert_data2(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/data2.txt")
	if err != nil {
		t.Fatal(err)
	}
	testConvert(t, data)
}

func testConvert(t *testing.T, data []byte) {
	var out *bytes.Buffer
	*format = "raw"
	out = &bytes.Buffer{}
	Convert(bytes.NewBuffer(data), out)
	if !bytes.Equal(data, out.Bytes()) {
		t.Fatal("want: data == out but:", data, out)
	}

	*format = "text"
	out = &bytes.Buffer{}
	Convert(bytes.NewBuffer(data), out)

	*format = "pretty"
	out = &bytes.Buffer{}
	Convert(bytes.NewBuffer(data), out)

	*format = "json"
	out = &bytes.Buffer{}
	Convert(bytes.NewBuffer(data), out)
	var m map[string]interface{}
	err := json.Unmarshal(out.Bytes(), &m)
	if err != nil {
		t.Fatal(err)
	}

	*format = "qfix"
	out = &bytes.Buffer{}
	Convert(bytes.NewBuffer(data), out)
}
