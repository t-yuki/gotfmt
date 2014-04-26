// +build testdata

package testdata

import (
	"runtime"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	runtime.Gosched()
	time.Sleep(time.Minute)
}
