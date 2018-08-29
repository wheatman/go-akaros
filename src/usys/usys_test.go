package usys_test

import (
	"testing"
	"usys"
)

func TestUsys(t *testing.T) {
	usys.Call(usys.USYS_FUNC, 1)
	usys.Call(usys.USYS_FUNC, 1,2)
	usys.Call(usys.USYS_FUNC, 1,2,3)
	usys.Call(usys.USYS_FUNC, 1,2,3,4)
	usys.Call(usys.USYS_FUNC, 1,2,3,4,5)
	usys.Call(usys.USYS_FUNC, 1,2,3,4,5,6)
	usys.Call(usys.USYS_FUNC, 1,2,3,4,5,6,7)
	usys.Call(usys.USYS_FUNC, 1,2,3,4,5,6,7,8)
	usys.Call(usys.USYS_FUNC, 1,2,3,4,5,6,7,8,9)
	usys.Call(usys.USYS_FUNC, 1,2,3,4,5,6,7,8,9,10)
	usys.Call(usys.USYS_FUNC, 1,2,3,4,5,6,7,8,9,10,11)
	usys.Call(usys.USYS_FUNC, 1,2,3,4,5,6,7,8,9,10,11,12)
	ret := usys.Call(usys.USYS_FUNC, 1, 2, 3, 4, 5, 6, 7)
	expect := int64(0x0007060504030201)
	if ret != expect {
		t.Errorf("usys.Call returned the wrong value for its arguments got %x, expected %x", ret, expect)
	}
}
