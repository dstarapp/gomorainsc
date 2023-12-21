package services

import (
	"github.com/astaxie/beego/logs"
	"github.com/dstarapp/gomorainsc/indexer"
)

var (
	gIndexer *indexer.Indexer
)

func StartIndexer(cfg *indexer.Config) error {
	gIndexer = indexer.NewIndexer(cfg)
	if err := gIndexer.Start(); err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func StopIndexer() {
	if gIndexer != nil {
		gIndexer.Stop()
	}
}

func GetIndexer() *indexer.Indexer {
	return gIndexer
}
