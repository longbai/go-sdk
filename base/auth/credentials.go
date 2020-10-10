package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/longbai/go-sdk/base"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"sort"
	"strings"

	"github.com/longbai/go-sdk/internal/dump"
)

//  七牛鉴权类，用于生成API签名
// AK/SK可以从 https://portal.qiniu.com/user/key 获取。
type Credentials struct {
	accessKey string
	secretKey []byte
}

func (cred *Credentials) Validate() bool {
	return len(cred.accessKey) > 8 && len(cred.secretKey) > 8
}

// 构建Credentials对象
func NewCredentials(accessKey, secretKey string) *Credentials {
	return &Credentials{accessKey, []byte(secretKey)}
}

// Sign 对数据进行签名，一般用于私有空间下载用途
func (cred *Credentials) Sign(data []byte) (token string) {
	h := hmac.New(sha1.New, cred.secretKey)
	h.Write(data)

	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("%s:%s", cred.accessKey, sign)
}

// SignWithData 对数据进行签名，一般用于上传凭证的生成用途
func (cred *Credentials) SignWithData(b []byte) (token string) {
	encodedData := base64.URLEncoding.EncodeToString(b)
	sign := cred.Sign([]byte(encodedData))
	return fmt.Sprintf("%s:%s", sign, encodedData)
}

func (cred *Credentials) APISignerV2() APISigner {
	return &v2{cred}
}

func (cred *Credentials) APISignerV1() APISigner {
	return &v1{cred}
}

type APISigner interface {
	Sign(req *http.Request) error
	Token(req *http.Request) (string, error)
	Verify(req *http.Request) (bool, error)

	collectData(req *http.Request) (data []byte, err error)
	includeBody(req *http.Request) bool
}

type v2 struct {
	*Credentials
}

func (signer *v2) Verify(req *http.Request) (bool, error) {
	return verify(signer, req, v2Prefix)
}

func (signer *v2) Token(req *http.Request) (string, error) {
	return token(signer, req, signer.Credentials)
}

func (signer *v2) Sign(req *http.Request) error {
	return sign(signer, req, v2Prefix)
}

type (
	xQiniuHeaderItem struct {
		HeaderName  string
		HeaderValue string
	}
	xQiniuHeaders []xQiniuHeaderItem
)

func (headers xQiniuHeaders) Len() int {
	return len(headers)
}

func (headers xQiniuHeaders) Less(i, j int) bool {
	if headers[i].HeaderName < headers[j].HeaderName {
		return true
	} else if headers[i].HeaderName > headers[j].HeaderName {
		return false
	} else {
		return headers[i].HeaderValue < headers[j].HeaderValue
	}
}

func (headers xQiniuHeaders) Swap(i, j int) {
	headers[i], headers[j] = headers[j], headers[i]
}

func (signer *v2) collectData(req *http.Request) (data []byte, err error) {
	u := req.URL
	//write method path?query
	s := fmt.Sprintf("%s %s", req.Method, u.Path)
	if u.RawQuery != "" {
		s += "?"
		s += u.RawQuery
	}

	//write host and post
	s += "\nHost: " + req.Host + "\n"

	//write content type
	contentType := req.Header.Get(cTypeHeaderKey)
	if contentType == "" {
		contentType = base.ContentTypeForm
		req.Header.Set(cTypeHeaderKey, contentType)
	}
	s += fmt.Sprintf("%s: %s\n", cTypeHeaderKey, contentType)

	xQiniuHeaders := make(xQiniuHeaders, 0, len(req.Header))
	for headerName, headerValues := range req.Header {
		if len(headerName) > len(xQiniuHeaderPrefix) && strings.HasPrefix(headerName, xQiniuHeaderPrefix) {
			for _, headerValue := range headerValues {
				xQiniuHeaders = append(xQiniuHeaders, xQiniuHeaderItem{
					HeaderName:  textproto.CanonicalMIMEHeaderKey(headerName),
					HeaderValue: headerValue,
				})
			}
		}
	}
	if len(xQiniuHeaders) > 0 {
		sort.Sort(xQiniuHeaders)
		for _, xQiniuHeader := range xQiniuHeaders {
			s += fmt.Sprintf("%s: %s\n", xQiniuHeader.HeaderName, xQiniuHeader.HeaderValue)
		}
	}
	s += "\n"
	return appendBody(signer, s, req)
}

func (signer *v2) includeBody(req *http.Request) bool {
	contentType := req.Header.Get("Content-Type")
	return req.Body != nil && (contentType == base.ContentTypeForm || contentType == base.ContentTypeJson)

}

type v1 struct {
	*Credentials
}

func (signer *v1) Verify(req *http.Request) (bool, error) {
	return verify(signer, req, "QBox ")
}

func (signer *v1) Token(req *http.Request) (string, error) {
	return token(signer, req, signer.Credentials)
}

func (signer *v1) Sign(req *http.Request) error {
	return sign(signer, req, v1Prefix)
}

func (signer *v1) collectData(req *http.Request) (data []byte, err error) {
	u := req.URL
	s := u.Path
	if u.RawQuery != "" {
		s += "?"
		s += u.RawQuery
	}
	s += "\n"
	return appendBody(signer, s, req)
}

func (signer *v1) includeBody(req *http.Request) bool {
	return req.Body != nil && req.Header.Get("Content-Type") == base.ContentTypeForm
}

func sign(signer APISigner, req *http.Request, prefix string) error {
	token, err := signer.Token(req)
	if err != nil {
		return err
	}
	req.Header.Add(authorizationHeaderKey, prefix+token)
	return nil
}

func verify(signer APISigner, req *http.Request, prefix string) (bool, error) {
	auth := req.Header.Get(authorizationHeaderKey)
	if auth == "" {
		return false, nil
	}

	token, err := signer.Token(req)
	if err != nil {
		return false, err
	}

	return auth == prefix+token, nil
}

func token(s APISigner, req *http.Request, c *Credentials) (string, error) {
	data, err := s.collectData(req)
	if err != nil {
		return "", err
	}
	token := c.Sign(data)
	return token, nil
}

func appendBody(signer APISigner, s string, req *http.Request) (data []byte, err error) {
	data = []byte(s)
	if signer.includeBody(req) {
		s2, rErr := dump.BytesFromRequest(req)
		if rErr != nil {
			err = rErr
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(s2))
		data = append(data, s2...)
	}
	return
}

const authorizationHeaderKey = "Authorization"
const v1Prefix = "QBox "
const v2Prefix = "Qiniu "
const cTypeHeaderKey = "Content-Type"
const xQiniuHeaderPrefix = "X-Qiniu-"
