package objectstore

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type FileUploadOptions struct {
	ChunkSizeBytes *int32
	Metadata       interface{}
}

type UploadStream struct {
	*gridfs.UploadStream
}

func (s UploadStream) FileID() FileID {
	return FileID{
		ObjectID: s.UploadStream.FileID.(primitive.ObjectID),
	}
}

type IUploadStream interface {
	io.WriteCloser
	FileID() FileID
}

type DownloadStream struct {
	*gridfs.DownloadStream
}

// WithContext взять таймаут из контекста и установить в ReadTimeout
func (s *DownloadStream) WithContext(ctx context.Context) IDownloadStream {
	if deadline, ok := ctx.Deadline(); ok {
		s.SetReadDeadline(deadline)
	}
	return s
}

// GetFileMeta получить информацию о загружаемом файле. В reader ничего нет,
// чтение файла только через методы io.Reader
func (s DownloadStream) GetFileMeta() *File {
	file := s.DownloadStream.GetFile()
	return &File{
		ID: FileID{
			ObjectID: file.ID.(primitive.ObjectID),
		},
		Name:   file.Name,
		Length: int(file.Length),
	}
}

type IDownloadStream interface {
	io.ReadCloser

	// WithContext взять таймаут из контекста и установить в ReadTimeout
	WithContext(ctx context.Context) IDownloadStream
	GetFileMeta() *File
}

type FileID struct {
	primitive.ObjectID
}

func (id FileID) String() string {
	s := id.ObjectID.String()
	return s[10 : len(s)-2]
}

var NilFileID FileID = FileID{primitive.NilObjectID}

type File struct {
	ID     FileID
	Name   string
	Length int
	Reader io.ReadCloser
}
