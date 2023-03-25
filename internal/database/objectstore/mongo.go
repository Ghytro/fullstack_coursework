package objectstore

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const DefaultFileDBName = "mydb"
const DefaultFileCollectionName = "fs.files"

type FileDB struct {
	*gridfs.Bucket
}

func NewFileDB(ctx context.Context, url string) *FileDB {
	mongoConn, err := mongo.Connect(ctx, options.Client().ApplyURI(url).SetMaxPoolSize(5))
	if err != nil {
		panic(err)
	}
	// таймаут на подключение
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	if err := mongoConn.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	bucket, err := gridfs.NewBucket(mongoConn.Database(DefaultFileDBName))
	if err != nil {
		panic(err)
	}
	return &FileDB{bucket}
}

func (db *FileDB) OpenUploadStream(fileName string, opts ...*FileUploadOptions) (IUploadStream, error) {
	unwrapedOpts := lo.Map(opts, func(o *FileUploadOptions, _ int) *options.UploadOptions {
		return &options.UploadOptions{
			ChunkSizeBytes: o.ChunkSizeBytes,
			Metadata:       o.Metadata,
		}
	})
	stream, err := db.Bucket.OpenUploadStream(fileName, unwrapedOpts...)
	if err != nil {
		return nil, err
	}
	return &UploadStream{
		UploadStream: stream,
	}, nil
}

func (db *FileDB) OpenDownloadStream(fileID FileID) (IDownloadStream, error) {
	stream, err := db.Bucket.OpenDownloadStream(fileID.ObjectID)
	if err != nil {
		return nil, err
	}
	return &DownloadStream{
		DownloadStream: stream,
	}, nil
}

func (db *FileDB) DeleteContext(ctx context.Context, fileID FileID) error {
	return db.Bucket.DeleteContext(ctx, fileID.ObjectID)
}
