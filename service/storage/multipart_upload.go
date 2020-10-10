package storage

import (
	"github.com/longbai/go-sdk/internal/net"
	"io"
)

func PutMultiPart(client net.HttpClient, hosts []string, reader io.ReaderAt, key *string, token string,
	size uint64, retry int, mime *string, params map[string]string, partSize int64) error {
	//uploadParts := makeUploadParts(int64(size), partSize)
	//return putMultiPart()
	return nil
}

func putMultiPart(client net.HttpClient, hosts []string, reader io.ReaderAt, key *string, token string,
	size uint64, retry int, mime *string, params map[string]string, partSize int64) error {
	return nil
}

func makeUploadParts(fsize, partSize int64) []int64 {
	partCnt := partNumber(fsize, partSize)
	uploadParts := make([]int64, partCnt)
	for i := 0; i < partCnt-1; i++ {
		uploadParts[i] = partSize
	}
	uploadParts[partCnt-1] = fsize - (int64(partCnt)-1)*partSize
	return uploadParts
}

func checkUploadParts(fsize int64, uploadParts []int64) bool {
	var partSize int64 = 0
	for _, size := range uploadParts {
		partSize += size
	}
	return fsize == partSize
}

func partNumber(fsize, partSize int64) int {
	return int((fsize + partSize - 1) / partSize)
}
