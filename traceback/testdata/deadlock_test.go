// +build testdata

package testdata

import "testing"

func TestDeadlock(t *testing.T) {
	ch1 := (chan bool)(nil)
	go func() {
		ch1 <- true
	}()
	<-ch1
}
