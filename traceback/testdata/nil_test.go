// +build testdata

package testdata

import "testing"

func TestNil(t *testing.T) {
	obj := error(nil)
	go func() { obj.Error() }()
}
