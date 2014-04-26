// +build testdata
// +build signal
// since os/signal.init creates a gorotine and blocks deadlock detector, we use build flag `signal`

package testdata

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestSIGABRT(t *testing.T) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	syscall.Kill(os.Getpid(), syscall.SIGABRT)
	<-ch
}
