package storage

import (
	"github.com/longbai/go-sdk/base/auth"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDirectUpload(t *testing.T) {
	up := NewUploader([]string{"http://up.qiniu.com"}, 4, false, nil)
	cd := auth.NewCredentials("zIV2FhwbOS0wKj3K1TmMRWOgyDSx6uTmMwu3mIRL", "pYP1XF8pWfdY-99Jx0gIDPEgKFcqIQDt_kH5N9tt")
	policty := PutPolicy{
		Scope: "test11:hello",
	}
	token := policty.UploadToken(cd)
	key := "hello"
	err := up.PutData([]byte("hello"), &key, token, nil, nil)
	require.NoError(t, err, "upload no error")
}
