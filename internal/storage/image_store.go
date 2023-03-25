package storage

import (
	"context"
	"io"

	"github.com/Ghytro/galleryapp/internal/database/objectstore"
)

type ImageStore struct {
	db ImageDBI
}

func NewImageStore(db *objectstore.FileDB) *ImageStore {
	return &ImageStore{db: db}
}

func (s *ImageStore) UploadFile(ctx context.Context, file *objectstore.File) (objectstore.FileID, error) {
	uploadStream, err := s.db.OpenUploadStream(file.Name)
	if err != nil {
		return objectstore.NilFileID, err
	}
	defer uploadStream.Close()

	fileID := uploadStream.FileID()

	if _, err := io.Copy(uploadStream, file.Reader); err != nil {
		return objectstore.NilFileID, err
	}
	return fileID, err
}

// DownloadFile скачать из объектного хранилища файл по id.
// downloadStrategy - коллбек выкачивающий файл из стрима куда-то и
// возвращающий стрим откуда полученный файл можно читать.
//
// WARNING: ресивер файла сам несет ответственность за закрытие стрима
// во избежание утечек при использовании DownloadStratTempFile
func (s *ImageStore) DownloadFile(
	ctx context.Context,
	fileID objectstore.FileID,
	downloadStrategy func(stream objectstore.IDownloadStream) (io.ReadCloser, error),
) (*objectstore.File, error) {
	downloadStream, err := s.db.OpenDownloadStream(fileID)
	if err != nil {
		return nil, err
	}
	meta := downloadStream.GetFileMeta()
	file, err := downloadStrategy(downloadStream.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return &objectstore.File{
		ID:     meta.ID,
		Name:   meta.Name,
		Reader: file,
	}, nil
}

func (s *ImageStore) DeleteFile(ctx context.Context, fileID objectstore.FileID) error {
	return s.db.DeleteContext(ctx, fileID)
}
