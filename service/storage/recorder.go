package storage

type UploadRecorder interface {
	Record(filePath, key string)
}
