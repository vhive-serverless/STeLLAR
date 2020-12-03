package connection

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestSetupExternalConnection(t *testing.T) {
	Initialize("www.google.com", "")
	assert(t, Singleton.ListAPIs() == nil, "External connection: ListAPIs() should return nil.")
	assert(t, Singleton.DeployFunction == nil, "External connection: DeployFunction should be nil.")
	assert(t, Singleton.RemoveFunction == nil, "External connection: RemoveFunction should be nil.")
	assert(t, Singleton.UpdateFunction == nil, "External connection: UpdateFunction should be nil.")
}

func TestSetupFileConnection(t *testing.T) {
	Initialize("vhive", "../../../../endpoints")
	equals(t, len(Singleton.ListAPIs()), 4)
	equals(t, 60., Singleton.ListAPIs()[0].ImageSizeMB)
	equals(t, int64(128), Singleton.ListAPIs()[0].FunctionMemoryMB)
	equals(t, "test1", Singleton.ListAPIs()[0].GatewayID)
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
