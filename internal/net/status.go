package net

import (
	"errors"
	"net/http"
	"strconv"
)

// HTTP status codes return by qiniu
// See: https://developer.qiniu.com/fusion/kb/1352/the-http-request-return-a-status-code
const (
	StatusPartSuccess = 298

	StatusInvalidCheckSum   = 406
	StatusAccountFrozen     = 419
	StatusImageSourceFailed = 478

	StatusManyVisit             = 573
	StatusCallbackFailed        = 579
	StatusServerOperationFailed = 599

	StatusResourceModified = 608
	StatusResourceDeleted  = 612
	StatusResourceExist    = 614
	StatusBucketMany       = 630
	StatusBucketNotFound   = 631
	StatusInvalidMarker    = 640
	StatusInvalidPartUpCtx = 701
)

var statusText = map[int]string{
	StatusPartSuccess:       "部分操作执行成功",
	StatusInvalidCheckSum:   "上传的数据 CRC32 校验错误",
	StatusAccountFrozen:     "用户账号被冻结",
	StatusImageSourceFailed: "镜像回源失败。 主要指镜像源服务器出现异常",

	StatusManyVisit:             "单个资源访问频率过高",
	StatusCallbackFailed:        "上传成功但是回调失败。 包括业务服务器异常；七牛服务器异常；服务器间网络异常",
	StatusServerOperationFailed: "服务端操作失败",
	StatusResourceModified:      "资源内容被修改",
	StatusResourceDeleted:       "指定资源不存在或已被删除",
	StatusResourceExist:         "目标资源已存在",
	StatusBucketMany:            "已创建的空间数量达到上限，无法创建新空间",
	StatusBucketNotFound:        "指定空间不存在",
	StatusInvalidMarker:         "调用列举资源 (list) 接口时，指定非法的marker参数",
	StatusInvalidPartUpCtx:      "在断点续上传过程中，后续上传接收地址不正确或ctx信息已过期",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}

func BadRequest(code int) error {
	if code < http.StatusInternalServerError || code/100 == 6 {
		if r := StatusText(code); r != "" {
			return errors.New(r)
		} else if r := http.StatusText(code); r != "" {
			return errors.New(r)
		}
		return errors.New("bad request " + strconv.FormatInt(int64(code), 10))
	}
	return nil
}

func StatusError(code int) error {
	if r := StatusText(code); r != "" {
		return errors.New(r)
	} else if r := http.StatusText(code); r != "" {
		return errors.New(r)
	}
	return errors.New("error response " + strconv.FormatInt(int64(code), 10))
}
