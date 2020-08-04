package storage

import (
	"fmt"
	"runtime"
	"testing"

	"../../base"
)

func TestVariable(t *testing.T) {
	appName := "test"

	SetAppName(appName)

	want := fmt.Sprintf("QiniuGo/%s (%s; %s; %s) %s", base.Version, runtime.GOOS, runtime.GOARCH, appName, runtime.Version())

	if UserAgent != want {
		t.Fail()
	}
}
