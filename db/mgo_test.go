package db

import (
	"log"
	"testing"
	"time"
)

type mongoConf struct{}

func (m mongoConf) GetURI() string {
	return "mongodb://127.0.0.1:27017"
}

func (m mongoConf) MaxPoolSize() uint64 {
	return 50
}

func (m mongoConf) MinPoolSize() uint64 {
	return 10
}

func (m mongoConf) MaxConnecting() uint64 {
	return 20
}

func (m mongoConf) ConnIdleTime() time.Duration {
	return 10 * time.Minute
}

func TestConnectMgo(t *testing.T) {
	log.Println(InitMongoDb(mongoConf{}))
}
