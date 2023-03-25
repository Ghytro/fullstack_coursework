package storage

import (
	"bytes"
	"io"
	"os"

	"github.com/Ghytro/galleryapp/internal/database/objectstore"
)

// DownloadStratTempFile стратегия загрузки, скачивающая файл на диск для дальнейшего чтения.
// Хорошо подходит для больших файлов, однако не такая быстрая как чтение из памяти
func DownloadStratTempFile(stream objectstore.IDownloadStream) (io.ReadCloser, error) {
	tmp, err := os.CreateTemp("/tmp", "galleryapp-*")
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(tmp, stream); err != nil {
		os.Remove(tmp.Name())
		return nil, err
	}
	return &tmpFileStream{
		File: tmp,
	}, nil
}

// tmpFileStream стрим обертка вокруг временного файла,
// который при вызове Close() удаляет temp file
//
// thoughts: может демон чистильщик или более
// умное хранилище будет лучше чтоб избежать костыля
// с bytes.Buffer?
type tmpFileStream struct {
	*os.File
}

// Close прекратить чтение из файла и удалить его
func (s *tmpFileStream) Close() error {
	if err := s.File.Close(); err != nil {
		return err
	}
	return os.Remove(s.File.Name())
}

// DownloadStratRAM стратегия загрузки, загружающая файл из бд в оперативную память
func DownloadStratRAM(stream objectstore.IDownloadStream) (io.ReadCloser, error) {
	var b bytes.Buffer
	if _, err := io.Copy(&b, stream); err != nil {
		return nil, err
	}
	return &closingBytesBuffer{
		Buffer: &b,
	}, nil
}

// closingBytesBuffer костыль чтобы работало чтение из файла
type closingBytesBuffer struct {
	*bytes.Buffer
}

func (b *closingBytesBuffer) Close() error {
	return nil
}
