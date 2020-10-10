package storage

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

type Uploader struct {
	partSize   uint32
	upHosts    []string
	recorder   UploadRecorder
	checkSum   bool
	client     *http.Client
	retryTimes int
}

func NewUploaderWithConf(conf *Config, recorder UploadRecorder) *Uploader {
	return &Uploader{
		partSize:   conf.PartSizeMb,
		upHosts:    conf.UpHosts,
		recorder:   recorder,
		checkSum:   conf.CheckSum,
		client:     http.DefaultClient,
		retryTimes: 3,
	}
}

func NewUploader(upHosts []string, partSizeMb uint32, checkSum bool, recorder UploadRecorder) *Uploader {
	return &Uploader{
		partSize: partSizeMb * 1024 * 1024,
		upHosts:  upHosts,
		recorder: recorder,
		checkSum: checkSum,
		client:   http.DefaultClient,
		retryTimes: 3,
	}
}

func (up *Uploader) PutFile(filePath string, key *string, token string,
	mime *string, params map[string]string) (err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()
	if size <= int64(up.partSize) {
		return PutStream(up.client, up.upHosts, f, key, token, uint64(size), up.retryTimes, mime, params)
	}
	return nil
}

func (up *Uploader) PutData(data []byte, key *string, token string,
	mime *string, params map[string]string) (err error) {
	return PutStream(up.client, up.upHosts, bytes.NewReader(data), key, token, uint64(len(data)), up.retryTimes, mime, params)
}

func (up *Uploader) PutStream(reader io.Reader, key *string, token string,
	size uint32, mime *string, params map[string]string) (err error) {
	return nil
}
