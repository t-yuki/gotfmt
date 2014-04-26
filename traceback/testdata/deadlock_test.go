// +build testdata

package testdata

import "testing"

func TestDeadlock(t *testing.T) {
	ch1 := make(chan bool)
	<-ch1
}
