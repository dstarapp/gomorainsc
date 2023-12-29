package indexer

import (
	"bytes"

	"github.com/astaxie/beego/logs"
	badger "github.com/dgraph-io/badger/v4"
)

type Bucket struct {
	tx   *badger.Txn
	name string
}

func NewBucket(tx *badger.Txn, name string) *Bucket {
	return &Bucket{
		tx:   tx,
		name: name,
	}
}

func (p *Bucket) Exist(key []byte) bool {
	bkey := p.fromkey(key)
	_, err := p.tx.Get(bkey)
	return err == nil
}

func (p *Bucket) Put(key []byte, val []byte) error {
	return p.tx.Set(p.fromkey(key), val)
}

func (p *Bucket) Get(key []byte) []byte {
	bkey := p.fromkey(key)
	item, err := p.tx.Get(bkey)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			logs.Error(string(key), err)
		}
		return nil
	}
	var data []byte
	err = item.Value(func(val []byte) error {
		data = append(data, val...)
		return nil
	})
	if err != nil {
		return nil
	}
	return data
}

func (p *Bucket) Delete(key []byte) error {
	return p.tx.Delete(p.fromkey(key))
}

func (p *Bucket) Scan(fn func([]byte, []byte)) error {
	it := p.tx.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	prefix := BucketPrefix(p.name)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		var data []byte
		err := item.Value(func(val []byte) error {
			data = append(data, val...)
			return nil
		})
		if err != nil {
			return err
		}
		fn(p.tokey(it.Item().Key()), data)
	}
	return nil
}

func (p *Bucket) ScanKey(fn func([]byte) bool) {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false

	it := p.tx.NewIterator(opts)
	defer it.Close()

	prefix := BucketPrefix(p.name)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		if !fn(p.tokey(it.Item().Key())) {
			break
		}
	}
}

func (p *Bucket) fromkey(key []byte) []byte {
	return append(BucketPrefix(p.name), key...)
}

func (p *Bucket) tokey(key []byte) []byte {
	return bytes.Replace(key, BucketPrefix(p.name), []byte(""), -1)
}

// func (p *Bucket) DropAll(db *badger.DB) error {
// 	prefix := append([]byte(p.name), '/')
// 	return db.DropPrefix(prefix)
// }

func BucketPrefix(name string) []byte {
	return append([]byte(name), '/')
}
