package indexer

import (
	"errors"
	"fmt"

	"github.com/aviate-labs/agent-go/principal"
	"github.com/dstarapp/gomorainsc/canisters/planet"
	"github.com/dstarapp/gomorainsc/moraprotocol"
)

type PlanetArticle struct {
	canister principal.Principal
	key      string
	aid      string
}

func (p *PlanetArticle) GetCanister() principal.Principal {
	return p.canister
}

func (p *PlanetArticle) String() string {
	return fmt.Sprintf("%s/%s", p.GetCanister(), p.aid)
}

func (p *PlanetArticle) FetchUpdate(indexer *Indexer) error {
	actor, err := moraprotocol.GetAnonyPlanet(p.canister)
	if err != nil {
		return err
	}

	resp, err := actor.QueryArticle(p.aid)
	if err != nil {
		return err
	}

	if resp.Err != nil {
		return errors.New(*resp.Err)
	}

	if resp.Ok == nil {
		return errors.New("QueryArticle result is nil")
	}

	return p.on_article(resp.Ok.Article, resp.Ok.Content, indexer)
}

func (p *PlanetArticle) on_article(article planet.QueryArticle, content string, indexer *Indexer) error {
	// on next verion, inscription's abstract will be forced to be consistent with the content.

	// second index
	// verify content
	// logs.Info(article.Id)
	// just skip

	return nil
}
