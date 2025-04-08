package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"time"
)

type DBInerface interface {
	GetURI() string
	MaxPoolSize() uint64
	MinPoolSize() uint64
	MaxConnecting() uint64
	ConnIdleTime() time.Duration
}

var mgoDBInstance *mongo.Client
var mgoClientOpt *options.ClientOptions

func InitMongoDb(in DBInerface) error {
	mgoClientOpt = options.Client().ApplyURI(in.GetURI()).
		SetMaxPoolSize(in.MaxPoolSize()).
		SetMinPoolSize(in.MinPoolSize()).
		SetMaxConnecting(in.MaxConnecting()).
		SetMaxConnIdleTime(in.ConnIdleTime())

	var err error
	mgoDBInstance, err = mongo.NewClient(mgoClientOpt)
	if err != nil {
		return fmt.Errorf("newClient %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = mgoDBInstance.Connect(ctx)
	if err != nil {
		return fmt.Errorf("Connect %v", err)
	}
	
	err = mgoDBInstance.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("ping %v %s", err, in.GetURI())
	}

	return nil
}

func MgoCollectionDo(database string, collectionName string, f func(col *mongo.Collection) error) error {
	col := mgoDBInstance.Database(database).Collection(collectionName)
	return f(col)
}

func MgoGetFile(database string, fileName string) ([]byte, error) {
	db := mgoDBInstance.Database(database)
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return nil, err
	}

	downStream, err := bucket.OpenDownloadStreamByName(fileName)
	if err != nil {
		return nil, err
	}

	defer downStream.Close()

	bs, err := ioutil.ReadAll(downStream)
	return bs, err
}

func MgoStoreFile(database string, fileName string, file []byte) error {
	db := mgoDBInstance.Database(database)
	bucket, err := gridfs.NewBucket(db)
	if err != nil {
		return err
	}

	up, err := bucket.OpenUploadStream(fileName)
	if err != nil {
		return err
	}

	defer up.Close()

	_, err = up.Write(file)
	if err != nil {
		return err
	}
	return nil
}
