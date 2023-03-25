package storage

import (
	"context"

	"github.com/Ghytro/galleryapp/internal/database/objectstore"
)

type ImageDBI interface {
	OpenUploadStream(filename string, opts ...*objectstore.FileUploadOptions) (objectstore.IUploadStream, error)
	OpenDownloadStream(id objectstore.FileID) (objectstore.IDownloadStream, error)
	DeleteContext(ctx context.Context, fileID objectstore.FileID) error
}
