package cmd_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/t-yuki/gotfmt/cmd"
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
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cmd.RegisterFlags(flags)
	var out *bytes.Buffer

	out = &bytes.Buffer{}
	cmd.ParseFlags(flags, []string{"-t=raw"})
	cmd.Convert(bytes.NewBuffer(data), out)
	if !bytes.Equal(data, out.Bytes()) {
		t.Fatalf("want: data == out but: %s != %s", data, out)
	}

	out = &bytes.Buffer{}
	cmd.ParseFlags(flags, []string{"-t=text"})
	cmd.Convert(bytes.NewBuffer(data), out)

	out = &bytes.Buffer{}
	cmd.ParseFlags(flags, []string{"-t=pretty"})
	cmd.Convert(bytes.NewBuffer(data), out)

	out = &bytes.Buffer{}
	cmd.ParseFlags(flags, []string{"-t=json"})
	cmd.Convert(bytes.NewBuffer(data), out)
	var m map[string]interface{}
	err := json.Unmarshal(out.Bytes(), &m)
	if err != nil {
		t.Fatal(err)
	}

	out = &bytes.Buffer{}
	cmd.ParseFlags(flags, []string{"-t=qfix"})
	cmd.Convert(bytes.NewBuffer(data), out)
}
