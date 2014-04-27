// +build testdata

package testdata

import "testing"

func TestCastNil(t *testing.T) {
	obj := interface{}(nil)
	go func() { obj.(error).Error() }()
}

func TestCastInvalid(t *testing.T) {
	obj := interface{}(t)
	go func() { obj.(error).Error() }()
}
