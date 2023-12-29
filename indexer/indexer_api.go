package indexer

import (
	"sort"

	"github.com/dstarapp/gomorainsc/inscription"
)

func (p *Indexer) GetSortFts() []*inscription.FT {
	fts := make(inscription.FTS, 0)
	for _, ft := range p.fts {
		fts = append(fts, ft)
	}
	sort.Sort(fts)
	return fts
}

func (p *Indexer) ScanFtItem(tick string, fn func(*inscription.MoraFTItem) bool) error {
	return p.db.ScanFtItem(tick, fn)
}

func (p *Indexer) IndexFtItem(item *inscription.MoraFTItem) error {
	return p.db.IndexFTItem(item)
}

func (p *Indexer) IndexFt(ft *inscription.FT) error {
	return p.db.IndexFT(ft)
}

func (p *Indexer) GetDB() *IndexerDB {
	return p.db
}
