package storage

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/longbai/go-sdk/internal/log"
	"github.com/longbai/go-sdk/internal/net"
)

//POST /put/<Fsize>/key/<EncodedKey>/mimeType/<EncodedMimeType>/crc32/<Crc32>/x:user-var/<EncodedUserVarVal>/x-qn-meta-<metaKey>/<EncodedMetaValue>
//Authorization: UpToken <UpToken>
//Content-Type: application/octet-stream
//
//<FileContent>

func buildPutUrl(host string, key *string, size uint64, mime *string, params map[string]string) string {
	var buf *bytes.Buffer
	if strings.HasPrefix(host, "http") {
		buf = bytes.NewBufferString(host)
	} else {
		buf = bytes.NewBufferString("http://")
		buf.WriteString(host)
	}
	buf.WriteString("/put/")
	buf.WriteString(strconv.FormatInt(int64(size), 10))
	if key != nil {
		buf.WriteString("/key/")
		buf.WriteString(net.Base64EncodeToString(*key))
	}
	if mime != nil {
		buf.WriteString("/mime/")
		buf.WriteString(net.Base64EncodeToString(*mime))
	}
	return buf.String()
}

func PutStream(client net.HttpClient, hosts []string, reader io.Reader, key *string, token string,
	size uint64, retry int, mime *string, params map[string]string) (re error) {

	h := net.NewHostList(hosts)
	for i := 0; i < retry; i++ {
		url := buildPutUrl(h.Next(), key, size, mime, params)
		req, err := http.NewRequest(http.MethodPost, url, reader)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "UpToken "+token)
		resp, err := client.Do(req)
		re = err
		if err != nil {
			log.Info("connection error", err)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		if re = net.BadRequest(resp.StatusCode); re != nil {
			return
		}
		re = net.StatusError(resp.StatusCode)
	}
	return
}
