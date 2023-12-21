package indexer

import (
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/principal"
	"github.com/dstarapp/gomorainsc/canisters/planet"
	"github.com/dstarapp/gomorainsc/moraprotocol"
)

type Planet struct {
	idx      int
	canister principal.Principal
}

func NewPlanet(pid principal.Principal) *Planet {
	return &Planet{
		canister: pid,
	}
}

func (p *Planet) String() string {
	return p.canister.String()
}

func (p *Planet) GetCanister() principal.Principal {
	return p.canister
}

func (p *Planet) SetIdx(idx int) {
	p.idx = idx
}

func (p *Planet) FetchUpdate(indexer *Indexer) error {
	actor, err := moraprotocol.GetAnonyPlanet(p.canister)
	if err != nil {
		return err
	}
	page := uint(1)
	size := uint(50)

	total := 0
	for {
		// first version : base on #Article, next version will be #Inscription
		resp, err := actor.QueryArticles(planet.QueryArticleReq{
			Atype: &planet.ArticleType{Article: &idl.Null{}},
			Page:  idl.NewNat(page),
			Size:  idl.NewNat(size),
			Sort:  planet.QuerySort{TimeAsc: &idl.Null{}},
		})
		if err != nil {
			return err
		}

		for _, article := range resp.Data {
			p.on_article(article, indexer)
		}

		// if p.canister.String() == "q3g4f-dqaaa-aaaan-qd6qq-cai" {
		// total += len(resp.Data)
		// tv, _ := strconv.Atoi(resp.Total.String())
		// logs.Info(p, p.GetCanister().String(), page, total, tv)
		// }

		tv, _ := strconv.Atoi(resp.Total.String())
		total = tv

		if !resp.Hasmore {
			break
		}

		// total += len(resp.Data)
		page++
	}
	logs.Info("[%d] planet FetchUpdate Finish planet = %v, total = %d", p.idx, p, total)

	return nil
}

func (p *Planet) on_article(article planet.QueryArticle, indexer *Indexer) error {

	//pre index
	ftitem := indexer.match_protocol_rule(p, article)
	if ftitem == nil {
		return nil
	}

	// if ftitem.Protocol.Op == "deploy" {
	// 	logs.Info(ftitem.Article, ftitem.Canister, ftitem.Content, ftitem.Created)
	// }

	indexer.PreProtocol(ftitem)
	return nil
}
