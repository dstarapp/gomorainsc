package indexer

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/astaxie/beego/logs"
	badger "github.com/dgraph-io/badger/v4"
	"github.com/dstarapp/gomorainsc/inscription"
)

const (
	BUCKET_PROTOCOL_SYSTEM        = "p_system"
	BUCKET_PROTOCOL_PREINDEX      = "p_preindex"
	BUCKET_PROTOCOL_PREINDEX_NEW  = "p_preindex_new"
	BUCKET_PROTOCOL_INDEX         = "p_index"
	BUCKET_PROTOCOL_PREFT         = "p_preft"
	BUCKET_PROTOCOL_FT            = "p_ft"
	BUCKET_PROTOCOL_TICK          = "pr_tick"
	BUCKET_PROTOCOL_DROPINDEX     = "p_dropindex"
	BUCKET_PROTOCOL_WHITELIST     = "p_whitelist"
	BUCKET_PROTOCOL_BLACKLIST     = "p_blacklist"
	BUCKET_PROTOCOL_BLACKCANISTER = "p_blackcanister"
)

var (
	INSCRIPTION_ID = []byte("inscription_id")
)

type IndexerDB struct {
	db *badger.DB
}

func NewIndexerDB(dbdir string) (*IndexerDB, error) {
	db, err := badger.Open(badger.DefaultOptions(dbdir))
	if err != nil {
		return nil, err
	}

	return &IndexerDB{
		db: db,
	}, nil
}

func (p *IndexerDB) Close() error {
	if p.db != nil {
		return p.db.Close()
	}

	return nil
}

// func (p *IndexerDB) read_insc_id() (uint64, error) {
// 	res := uint64(0)
// 	err := p.db.View(func(tx *badger.Txn) error {
// 		sys := NewBucket(tx, BUCKET_PROTOCOL_SYSTEM)
// 		if sys == nil {
// 			return nil
// 		}
// 		data := sys.Get(INSCRIPTION_ID)
// 		if data == nil {
// 			return nil
// 		}

// 		val, err := strconv.ParseUint(string(data), 10, 64)
// 		if err != nil {
// 			return err
// 		}
// 		res = val
// 		return nil
// 	})
// 	return res, err
// }

// func (p *IndexerDB) write_insc_id(id uint64) error {
// 	err := p.db.Update(func(tx *badger.Txn) error {
// 		sys := NewBucket(tx, BUCKET_PROTOCOL_SYSTEM)
// 		data := fmt.Sprintf("%d", id)
// 		return sys.Put(INSCRIPTION_ID, []byte(data))
// 	})
// 	return err
// }

// func (p *IndexerDB) write_insc_id_bytx(tx *badger.Txn, id uint64) error {
// 	sys := NewBucket(tx, BUCKET_PROTOCOL_SYSTEM)
// 	data := fmt.Sprintf("%d", id)
// 	return sys.Put(INSCRIPTION_ID, []byte(data))
// }

func (p *IndexerDB) LoadFT() (map[string]*inscription.FT, error) {
	if p.db == nil {
		return nil, errors.New("db is nil")
	}

	ftmaps := make(map[string]*inscription.FT)
	err := p.db.View(func(tx *badger.Txn) error {
		nbucket := NewBucket(tx, BUCKET_PROTOCOL_FT)
		bucket := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
		nbucket.ScanKey(func(key []byte) bool {
			val := bucket.Get(key)
			if val == nil {
				logs.Info("skip not exist in preindex", string(key))
				return true
			}
			var item inscription.FT
			if err := json.Unmarshal(val, &item); err != nil {
				return true
			}
			if _, ok := ftmaps[item.Ticker]; !ok {
				ftmaps[item.Ticker] = &item
			}
			return true
		})
		return nil
	})

	for _, ft := range ftmaps {
		ft.ConfirmCount, ft.UnconfirmCount = p.CalcCount(ft)
	}

	return ftmaps, err
}

func (p *IndexerDB) check_pre_protocol(item *inscription.MoraFTItem) bool {
	if p.db == nil {
		return false
	}
	has := false
	err := p.db.View(func(tx *badger.Txn) error {
		key := mft_encode_key(item)
		preb := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
		if preb.Get(key) != nil {
			has = true
		}
		return nil
	})
	if err != nil {
		return false
	}
	return has
}

func (p *IndexerDB) pre_protocol(item *inscription.MoraFTItem) error {
	if p.check_pre_protocol(item) {
		return nil
	}

	err := p.db.Update(func(tx *badger.Txn) error {
		bucket := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
		key := mft_encode_key(item)
		data, _ := json.Marshal(item)
		bucket.Put(key, data)

		nbucket := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX_NEW)
		nbucket.Put(key, []byte{1})
		return nil
	})
	return err
}

// func (p *IndexerDB) DropProtocol(item *inscription.MoraFTItem) error {
// 	err := p.db.Batch(func(tx *badger.Txn) error {
// 		key := mft_encode_key(item)

// 		dropb, err := tx.CreateBucketIfNotExists(BUCKET_PROTOCOL_DROPINDEX)
// 		if err != nil {
// 			return err
// 		}
// 		err = dropb.Put(key, []byte{0})
// 		if err != nil {
// 			return err
// 		}

// 		preb, err := tx.CreateBucketIfNotExists(BUCKET_PROTOCOL_PREINDEX)
// 		if err != nil {
// 			return err
// 		}
// 		return preb.Delete(key)
// 	})
// 	return err
// }

//	func (p *IndexerDB) IndexProtocol(item *inscription.MoraFTItem) error {
//		err := p.db.Batch(func(tx *badger.Txn) error {
//			key := mft_encode_key(item)
//			indexb, err := tx.CreateBucketIfNotExists(BUCKET_PROTOCOL_INDEX)
//			if err != nil {
//				return err
//			}
//			data, _ := json.Marshal(item)
//			err = indexb.Put(key, data)
//			if err != nil {
//				return err
//			}
//			preb, err := tx.CreateBucketIfNotExists(BUCKET_PROTOCOL_PREINDEX)
//			if err != nil {
//				return err
//			}
//			return preb.Delete(key)
//		})
//		return err
//	}
func (p *IndexerDB) PreIndexFT(ft *inscription.FT, oft *inscription.FT) error {

	if p.db == nil {
		return errors.New("db is nil")
	}

	// ft.Metadata.Verify = true

	err := p.db.Update(func(tx *badger.Txn) error {
		item := ft.Metadata
		key := mft_encode_key(item)

		// update pre index data
		{
			indexb := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
			data, _ := json.Marshal(ft)
			if err := indexb.Put(key, data); err != nil {
				return err
			}
		}

		// add ft
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_FT)
			err := ftb.Put(mft_encode_key(item), []byte{0})
			if err != nil {
				return err
			}
			if oft != nil {
				okey := mft_encode_key(oft.Metadata)
				if err := ftb.Delete(okey); err != nil {
					return err
				}
			}
		}

		// if err := p.write_insc_id_bytx(tx, p.inscId); err != nil {
		// 	return err
		// }
		return nil
	})
	return err
}

func (p *IndexerDB) UpdateFTCache(ft *inscription.FT) error {
	if p.db == nil {
		return errors.New("db is nil")
	}

	// ft.Metadata.Verify = true
	err := p.db.Update(func(tx *badger.Txn) error {
		item := ft.Metadata
		key := mft_encode_key(item)
		// update pre index data
		{
			indexb := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
			data, _ := json.Marshal(ft)
			if err := indexb.Put(key, data); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (p *IndexerDB) UpdateItemCache(item *inscription.MoraFTItem) error {
	if p.db == nil {
		return errors.New("db is nil")
	}

	// ft.Metadata.Verify = true
	err := p.db.Update(func(tx *badger.Txn) error {
		key := mft_encode_key(item)
		// update pre index data
		{
			indexb := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
			data, _ := json.Marshal(item)
			if err := indexb.Put(key, data); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (p *IndexerDB) PreIndexFTItem(item *inscription.MoraFTItem) error {
	if p.db == nil {
		return errors.New("db is nil")
	}

	// p.inscMtx.Lock()
	// p.inscId++
	// p.inscMtx.Unlock()

	// item.ID = p.inscId
	// item.Verify = true

	err := p.db.Update(func(tx *badger.Txn) error {
		key := mft_encode_key(item)

		// update index FTItem inscription
		{
			indexb := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
			data, _ := json.Marshal(item)
			if err := indexb.Put(key, data); err != nil {
				return err
			}
		}

		// index ticker list
		{
			bucketname := fmt.Sprintf("%s_%s", string(BUCKET_PROTOCOL_TICK), item.Ticker)
			indexb := NewBucket(tx, bucketname)
			err := indexb.Put(key, []byte{0})
			if err != nil {
				return err
			}
		}

		// if err := p.write_insc_id_bytx(tx, p.inscId); err != nil {
		// 	return err
		// }
		return nil
	})
	return err
}

func (p *IndexerDB) ScanPreProtocolNew(fn func(item *inscription.MoraFTItem) bool) error {
	err := p.db.View(func(tx *badger.Txn) error {
		nbucket := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX_NEW)
		bucket := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
		nbucket.ScanKey(func(key []byte) bool {
			val := bucket.Get(key)
			if val == nil {
				return true
			}
			var item inscription.MoraFTItem
			if err := json.Unmarshal(val, &item); err != nil {
				return true
			}
			return fn(&item)
		})
		return nil
	})

	return err
}

func (p *IndexerDB) FlushPreProtocolNew() error {
	return p.db.DropPrefix(BucketPrefix(BUCKET_PROTOCOL_PREINDEX_NEW))
}

func (p *IndexerDB) IndexFT(ft *inscription.FT) error {

	if p.db == nil {
		return errors.New("db is nil")
	}

	ft.Metadata.Verify = true

	err := p.db.Update(func(tx *badger.Txn) error {
		item := ft.Metadata
		key := mft_encode_key(item)

		// update pre index data
		{
			indexb := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
			data, _ := json.Marshal(ft)
			if err := indexb.Put(key, data); err != nil {
				return err
			}
		}

		// index ft
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_FT)
			err := ftb.Put(mft_encode_key(item), []byte{1})
			if err != nil {
				return err
			}
		}

		// put into global index
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_INDEX)
			err := ftb.Put(mft_encode_key(item), []byte{1})
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

func (p *IndexerDB) IndexFTItem(item *inscription.MoraFTItem) error {
	if p.db == nil {
		return errors.New("db is nil")
	}

	item.Verify = true
	err := p.db.Update(func(tx *badger.Txn) error {
		key := mft_encode_key(item)

		// update index FTItem inscription
		{
			indexb := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
			data, _ := json.Marshal(item)
			if err := indexb.Put(key, data); err != nil {
				return err
			}
		}

		// index ticker list
		{
			bucketname := fmt.Sprintf("%s_%s", string(BUCKET_PROTOCOL_TICK), item.Ticker)
			indexb := NewBucket(tx, bucketname)
			err := indexb.Put(key, []byte{1})
			if err != nil {
				return err
			}
		}

		// put into global index
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_INDEX)
			err := ftb.Put(mft_encode_key(item), []byte{1})
			if err != nil {
				return err
			}
		}

		// if err := p.write_insc_id_bytx(tx, p.inscId); err != nil {
		// 	return err
		// }
		return nil
	})
	return err
}

func (p *IndexerDB) DeleteFTItem(item *inscription.MoraFTItem) error {
	if p.db == nil {
		return errors.New("db is nil")
	}

	err := p.db.Update(func(tx *badger.Txn) error {
		key := mft_encode_key(item)

		// update index FTItem inscription
		{
			indexb := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
			indexb.Delete(key)
		}

		// index ticker list
		{
			bucketname := fmt.Sprintf("%s_%s", string(BUCKET_PROTOCOL_TICK), item.Ticker)
			indexb := NewBucket(tx, bucketname)
			indexb.Delete(key)
		}

		// put into global index
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_INDEX)
			ftb.Delete(key)
		}
		return nil
	})
	return err
}

func (p *IndexerDB) ScanFtItem(tick string, fn func(*inscription.MoraFTItem) bool) error {
	err := p.db.View(func(tx *badger.Txn) error {
		bucketname := fmt.Sprintf("%s_%s", string(BUCKET_PROTOCOL_TICK), tick)
		indexb := NewBucket(tx, bucketname)
		bucket := NewBucket(tx, BUCKET_PROTOCOL_PREINDEX)
		indexb.ScanKey(func(key []byte) bool {
			val := bucket.Get(key)
			if val == nil {
				return true
			}
			var item inscription.MoraFTItem
			if err := json.Unmarshal(val, &item); err != nil {
				return true
			}
			return fn(&item)
		})
		return nil
	})

	return err
}

func (p *IndexerDB) CalcMinters(ft *inscription.FT) (int, error) {
	minters := make(map[string]int)
	count := 0
	err := p.ScanFtItem(ft.Ticker, func(item *inscription.MoraFTItem) bool {
		count++
		if val, ok := minters[item.Owner]; ok {
			minters[item.Owner] = val + 1
		} else {
			minters[item.Owner] = 1
		}
		return count < int(ft.MaxItem)
	})
	return len(minters), err
}

func (p *IndexerDB) CalcCount(ft *inscription.FT) (int, int) {
	confirm_count := 0
	unconfirm_count := 0
	p.ScanFtItem(ft.Ticker, func(item *inscription.MoraFTItem) bool {
		if item.Verify {
			confirm_count++
		} else {
			unconfirm_count++
		}
		return true
	})
	return confirm_count, unconfirm_count
}

func (p *IndexerDB) PutWhiteList(owner string) error {
	if p.db == nil {
		return errors.New("db is nil")
	}

	err := p.db.Update(func(tx *badger.Txn) error {
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_WHITELIST)
			ftb.Put([]byte(owner), []byte{1})
		}
		return nil
	})
	return err
}

func (p *IndexerDB) PutBlackList(owner string) error {
	if p.db == nil {
		return errors.New("db is nil")
	}

	err := p.db.Update(func(tx *badger.Txn) error {
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_BLACKLIST)
			ftb.Put([]byte(owner), []byte{1})
		}
		return nil
	})
	return err
}

func (p *IndexerDB) PutBlackCanister(canister string) error {
	if p.db == nil {
		return errors.New("db is nil")
	}

	err := p.db.Update(func(tx *badger.Txn) error {
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_BLACKCANISTER)
			ftb.Put([]byte(canister), []byte{1})
		}
		return nil
	})
	return err
}

func (p *IndexerDB) CheckWhite(owner string) bool {

	has := false
	err := p.db.View(func(tx *badger.Txn) error {
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_WHITELIST)
			if ftb.Get([]byte(owner)) != nil {
				has = true
			}
		}
		return nil
	})
	if err != nil {
		return false
	}
	return has
}

func (p *IndexerDB) CheckBlack(owner string) bool {
	has := false
	err := p.db.View(func(tx *badger.Txn) error {
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_BLACKLIST)
			if ftb.Get([]byte(owner)) != nil {
				has = true
			}
		}
		return nil
	})
	if err != nil {
		return false
	}
	return has
}

func (p *IndexerDB) CheckBlackCanister(canister string) bool {
	has := false
	err := p.db.View(func(tx *badger.Txn) error {
		{
			ftb := NewBucket(tx, BUCKET_PROTOCOL_BLACKCANISTER)
			if ftb.Get([]byte(canister)) != nil {
				has = true
			}
		}
		return nil
	})
	if err != nil {
		return false
	}
	return has
}
