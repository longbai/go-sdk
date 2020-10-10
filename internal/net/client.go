package net

import (
	"fmt"
	"github.com/longbai/go-sdk/base/auth"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/longbai/go-sdk/base"
)

var userAgent = fmt.Sprintf(
	"QiniuGO/%s (%s; %s; %s) %s", base.Version, runtime.GOOS, runtime.GOARCH, runtime.Version(), appName())

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func appName() string {
	p, _ := os.Executable()
	app, _ := filepath.EvalSymlinks(p)
	return path.Base(app)
}

type Client struct {
}

type CredentialedClient struct {
	mac *auth.Credentials
}
