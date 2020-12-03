package util

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestBytesToMB(t *testing.T) {
	res := BytesToMB(1024. * 1024.)
	equals(t, 1., res)
}

func TestMBToBytes(t *testing.T) {
	res := MBToBytes(7.)
	equals(t, int64(7 * 1024 * 1024), res)

	res = MBToBytes(0.)
	equals(t, int64(0), res)
}

func TestAlmostEqualFloats(t *testing.T) {
	eqTheshold := 0.1

	res := AlmostEqualFloats(5.7e-15, 0., eqTheshold)
	equals(t, res, true)

	res = AlmostEqualFloats(7.4, 7.5, eqTheshold)
	equals(t, res, true)

	res = AlmostEqualFloats(4.4, 4.51, eqTheshold)
	equals(t, res, false)

	res = AlmostEqualFloats(0., 2.9, eqTheshold)
	equals(t, res, false)
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
