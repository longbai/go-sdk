package internal

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"os"
	"time"
)

type reqIdKey struct{}

// WithReqId 把reqId加入context中
func WithReqId(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, reqIdKey{}, reqId)
}

// ReIdFromContext 从context中获取reqId
func ReqIdFromContext(ctx context.Context) (reqId string, ok bool) {
	reqId, ok = ctx.Value(reqIdKey{}).(string)
	return
}

func genReqId() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], uint32(os.Getpid()))
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}
