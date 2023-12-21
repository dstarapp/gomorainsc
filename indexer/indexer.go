package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/aviate-labs/agent-go/principal"
	"github.com/dstarapp/gomorainsc/inscription"
	"github.com/dstarapp/gomorainsc/moraprotocol"
	"github.com/eapache/channels"
)

type Config struct {
	Protocol  string   `toml:"protocol"`
	VerifyTag bool     `toml:"verify_tag"`
	Num       int      `toml:"num"`
	PreNum    int      `toml:"pre_num"`
	TickMin   int      `toml:"tick_min"` // ticker minute
	DbPath    string   `toml:"db_path"`
	MoraFile  string   `toml:"mora_file"`
	BlackIDs  []string `toml:"black_ids"`
}

type Indexer struct {
	cfg         *Config
	indexwait   sync.WaitGroup
	closechan   chan bool
	db          *IndexerDB
	fts         map[string]*inscription.FT
	ftmtx       sync.RWMutex
	preitemchan *channels.InfiniteChannel
	preitemwait *sync.WaitGroup
}

func NewIndexer(cfg *Config) *Indexer {
	return &Indexer{
		cfg:       cfg,
		closechan: make(chan bool),
		fts:       make(map[string]*inscription.FT),
	}
}

func (p *Indexer) Start() error {
	db, err := NewIndexerDB(p.cfg.DbPath)
	if err != nil {
		return err
	}
	p.db = db

	if p.cfg.MoraFile != "" {
		if err := p.PreImportMora(p.cfg.MoraFile); err != nil {
			return err
		}
	}

	p.fts, err = p.db.LoadFT()
	if err != nil {
		return err
	}
	p.PrintFts()

	p.indexwait.Add(1)
	go p.do_ticker()
	return nil
}

func (p *Indexer) Stop() {
	close(p.closechan)
	p.indexwait.Wait()

	for _, ft := range p.fts {
		p.db.UpdateFTCache(ft)
	}

	if err := p.db.Close(); err != nil {
		logs.Error(err)
	}
}

func (p *Indexer) GetFt(tick string) *inscription.FT {
	p.ftmtx.RLock()
	defer p.ftmtx.RUnlock()
	ft := p.fts[tick]
	return ft
}

func (p *Indexer) do_ticker() {
	defer p.indexwait.Done()

	if err := p.do_task(); err != nil {
		logs.Error("fisrt spider_planets", err)
	}

	ticker := time.NewTicker(time.Minute * time.Duration(p.cfg.TickMin))
	for {
		select {
		case <-ticker.C:
			if err := p.do_task(); err != nil {
				logs.Error("spider_planets", err)
			}
		case <-p.closechan:
			logs.Info("stop close!")
			return
		}
	}
}

func (p *Indexer) do_task() error {
	p.init_pre_index()
	if err := p.do_spider_planets(); err != nil {
		return err
	}
	p.wait_pre_index()

	if err := p.do_indexer_articles(); err != nil {
		return err
	}

	return nil
}

func (p *Indexer) init_pre_index() {

	p.preitemchan = channels.NewInfiniteChannel()
	p.preitemwait = &sync.WaitGroup{}

	for i := 0; i < p.cfg.PreNum; i++ {
		p.preitemwait.Add(1)
		go p.do_pre_index()
	}
}

func (p *Indexer) do_pre_index() {
	defer p.preitemwait.Done()

	for it := range p.preitemchan.Out() {
		item := it.(*inscription.MoraFTItem)
		if item == nil {
			continue
		}
		if err := p.db.pre_protocol(item); err != nil {
			logs.Error(err)
		}
	}
}

func (p *Indexer) wait_pre_index() {
	p.preitemchan.Close()
	p.preitemwait.Wait()

	p.preitemchan = nil
	p.preitemwait = nil
}

func (p *Indexer) PreProtocol(item *inscription.MoraFTItem) {
	if p.preitemchan != nil {
		p.preitemchan.In() <- item
	}
}

func (p *Indexer) do_spider(waiter *sync.WaitGroup, planetchan chan *Planet) {
	defer waiter.Done()

	for planet := range planetchan {
		for i := 0; i < 3; i++ {
			if err := planet.FetchUpdate(p); err == nil {
				break
			} else {
				logs.Error("planet FetchUpdate error, id = %v, error = %v", planet, err)
			}
		}
	}
}

func (p *Indexer) do_spider_planets() error {
	logs.Info("do_spider_planets start!")
	start := time.Now()
	planets, err := moraprotocol.GetAllPlantCanisterIds()
	if err != nil {
		return err
	}

	planetwait := sync.WaitGroup{}
	planetchan := make(chan *Planet, p.cfg.Num)
	for i := 0; i < p.cfg.Num; i++ {
		planetwait.Add(1)
		go p.do_spider(&planetwait, planetchan)
	}

	logs.Info("total planets: ", len(planets))
	for idx, pid := range planets {
		if p.checkBlack(pid) {
			continue
		}
		planet := NewPlanet(pid)
		planet.SetIdx(idx)
		planetchan <- planet
	}

	close(planetchan)
	planetwait.Wait()

	logs.Info("do_spider_planets done! ", time.Since(start).Seconds(), "s")
	return nil
}

func (p *Indexer) do_indexer_articles() error {

	logs.Info("do_indexer_articles start")
	//mark fullspider count
	if err := p.db.ScanPreProtocolNew(p.visitFtItemNew); err != nil {
		return err
	}
	for _, ft := range p.fts {
		ft.FullCount++
		ft.LastUpdated = uint64(time.Now().UnixMilli())
		ft.Minters, _ = p.db.CalcMinters(ft)
		p.db.UpdateFTCache(ft)
	}
	p.PrintFts()
	logs.Info("do_indexer_articles done!")
	return p.db.FlushPreProtocolNew()
}

func (p *Indexer) visitFtItemNew(item *inscription.MoraFTItem) bool {
	// not match base, skip it
	if !p.match_protocol_base(item) {
		return true
	}

	switch item.Protocol.Op {
	case "deploy":
		if _, err := p.createFT(item); err != nil {
			logs.Error("create FT", err)
		}
	case "mint":
		if _, err := p.createFTItem(item); err != nil {
			logs.Error("create FTItem", err)
		}
	}
	return true
}

func (p *Indexer) createFT(item *inscription.MoraFTItem) (bool, error) {
	ft := p.match_protocol_deploy(item)
	if ft == nil || ft.Metadata == nil {
		return false, fmt.Errorf("%s/%s is not match protocol deploy", item.Canister, item.Article)
	}

	if ft == nil {
		return false, errors.New("not protocol deploy: " + item.Content)
	}

	tick := strings.ToLower(ft.Ticker)
	oft, ok := p.fts[tick]
	if ok {
		if oft.Metadata.Article == ft.Metadata.Article && oft.Metadata.Canister == ft.Metadata.Canister {
			return false, nil
		}
		if ft.DeployTime > oft.DeployTime {
			return false, fmt.Errorf("tick exist: %s/%s", ft.Metadata.Canister, ft.Metadata.Article)
		}
		if oft.Done {
			return false, errors.New("tick is exist and done: " + item.Content)
		}
		ft.FullCount = oft.FullCount
		ft.ConfirmCount = oft.ConfirmCount
		ft.UnconfirmCount = oft.UnconfirmCount
	}

	if err := p.db.PreIndexFT(ft, oft); err != nil {
		return false, err
	}

	p.fts[tick] = ft
	//call FT manager canister to deploy canister

	return true, nil
}

func (p *Indexer) createFTItem(item *inscription.MoraFTItem) (bool, error) {

	ft, ok := p.fts[item.Ticker]
	if !ok || ft == nil {
		//just skip
		return false, nil
	}

	if ft.Done {
		return false, nil
	}

	if !p.match_inscription(ft, item) {
		return false, nil
	}

	if err := p.db.PreIndexFTItem(item); err != nil {
		return false, err
	}

	ft.UnconfirmCount++

	return true, nil
}

func (p *Indexer) checkBlack(pid principal.Principal) bool {
	for _, str := range p.cfg.BlackIDs {
		if pid.String() == str {
			return true
		}
	}
	return false
}

func (p *Indexer) PrintFts() {
	logs.Info("--------------------- FT Data ---------------------")
	for _, ft := range p.fts {
		logs.Info("tick = %s, done = %v, max = %d, total = %d, minters = %d",
			ft.Ticker, ft.Done, int(ft.MaxItem), ft.ConfirmCount+ft.UnconfirmCount, ft.Minters)
	}
	logs.Info("--------------------- FT Data ---------------------")
}

// import mora from txt
// keep consistency
// https://mora.app/planet/qvsfp-6aaaa-aaaan-qdbua-cai/0C2G8XYQ1Y3DVFX0PJT58QW2DT
func (p *Indexer) PreImportMora(file string) error {
	body, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var items inscription.FTItems
	if err := json.Unmarshal(body, &items); err != nil {
		return err
	}

	tick := "mora"
	for idx, item := range items {
		if item.Title == "" {
			item.Title = "MORA"
		}
		item.Protocol = match_mora_protocol(item.Content)
		if item.Protocol == nil {
			logs.Error(idx, item.Content, item.Article)
			continue
		}
		if item.Protocol.Op == "deploy" {
			if _, err := p.createFT(item); err != nil {
				logs.Error(idx, err)
				return err
			}
			tick = item.Protocol.Tick
		} else if item.Protocol.Op == "mint" {
			if _, err := p.createFTItem(item); err != nil {
				logs.Error(idx, err)
				return err
			}
			item.MintTx = fmt.Sprintf("%d", idx-1)
			p.db.IndexFTItem(item)
		}
	}

	ft := p.GetFt(tick)
	if ft != nil {
		ft.Canister = "hbjjz-kaaaa-aaaan-qiocq-cai"
		ft.Index = "3kq54-gyaaa-aaaan-qiyxa-cai"
		ft.Done = true
		ft.FullCount = 1
		ft.LastUpdated = 1702283841534
		ft.ConfirmCount = ft.UnconfirmCount
		ft.UnconfirmCount = 0
		ft.Minters, _ = p.db.CalcMinters(ft)
		p.db.IndexFT(ft)
	}

	return nil
}

// func (p *Indexer) do_indexer(waiter *sync.WaitGroup, articlechan chan *PlanetArticle) {
// 	defer waiter.Done()
// 	for article := range articlechan {
// 		for i := 0; i < 3; i++ {
// 			if err := article.FetchUpdate(p); err == nil {
// 				continue
// 			} else {
// 				logs.Error("article FetchUpdate error, id = %v, error = %v", article, err)
// 			}
// 		}
// 	}
// }
