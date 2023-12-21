package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/dstarapp/gomorainsc/canisters/planet"
	"github.com/dstarapp/gomorainsc/inscription"
	"github.com/dstarapp/gomorainsc/utils"
)

const (
	MAX_DEPLOY_ITEM = 1000000
)

// version 0
// content must == abstract
// 1. public article
// 2. createdtime == updatedtime or updatedtime == 0
// 3. tick must be lower case
func (p *Indexer) match_protocol_rule(planet *Planet, article planet.QueryArticle) *inscription.MoraFTItem {

	// this version, just verify #public
	// next version, must be  #lock
	if article.Status.Public == nil {
		return nil
	}

	created, _ := strconv.ParseUint(article.Created.String(), 10, 64)
	updated, _ := strconv.ParseUint(article.Updated.String(), 10, 64)

	//firtst mora inscription deploy time
	if created < 1702043622731 {
		return nil
	}

	if updated != 0 && created != updated {
		// nnsdao deploy upated by author, skip it.
		if article.Id == "0C6KQC7G8NZ4TQY12SZP39046E" {
			//nnsdao deploy skip
			logs.Info("nnsdao", article.Abstract, article.Id, planet.String(), article.Created, article.Updated)
		} else {
			return nil
		}
	}

	// "p": "mora-20"
	abstract := strings.TrimSpace(article.Abstract)
	data := utils.RegSub(abstract, `{(.*)}$`)
	if data == "" {
		return nil
	}

	data = fmt.Sprintf("{%s}", data)

	var protocol inscription.Mora20Protocol
	if err := json.Unmarshal([]byte(data), &protocol); err != nil {
		return nil
	}

	//protocol not mora-20
	if protocol.P != p.cfg.Protocol {
		return nil
	}

	// not lower case
	if strings.ToLower(protocol.Tick) != protocol.Tick {
		return nil
	}
	// protocol.Tick = strings.ToLower(protocol.Tick)

	// tick is done
	ft := p.GetFt(protocol.Tick)
	if ft != nil && ft.Done {
		return nil
	}

	tags := strings.Join(article.Tags, "")

	return &inscription.MoraFTItem{
		Ticker:   protocol.Tick,
		Owner:    article.Author.String(),
		Canister: planet.GetCanister().String(),
		Article:  article.Id,
		Title:    article.Title,
		Thumb:    article.Thumb,
		Content:  article.Abstract,
		Created:  created,
		Updated:  updated,
		Tag:      tags,
		Protocol: &protocol,
	}
}

// 1. public article
// 2. createdtime == updatedtime or updatedtime == 0
// 3. match title (upper case)
// 4. match tags (option)
func (p *Indexer) match_protocol_base(item *inscription.MoraFTItem) bool {
	if !strings.EqualFold(item.Title, item.Ticker) {
		return false
	}

	if p.cfg.VerifyTag && !strings.EqualFold(item.Tag, "$"+item.Protocol.Tick) {
		return false
	}

	return true
}

func (p *Indexer) match_protocol_deploy(item *inscription.MoraFTItem) *inscription.FT {

	if item.Protocol.Op != "deploy" {
		return nil
	}

	max, err := strconv.ParseUint(item.Protocol.Max, 10, 64)
	if err != nil {
		return nil
	}
	if max <= 0 {
		return nil
	}

	lim, err := strconv.ParseUint(item.Protocol.Lim, 10, 64)
	if err != nil {
		return nil
	}
	maxitem := max / lim
	if maxitem > MAX_DEPLOY_ITEM {
		return nil
	}
	// logs.Info(item.Protocol)

	return &inscription.FT{
		Ticker:     item.Protocol.Tick,
		Image:      item.Thumb,
		Limit:      lim,
		Max:        max,
		MaxItem:    maxitem,
		DeployTime: item.Created,
		Deployer:   item.Owner,
		Done:       false,
		Metadata:   item,
	}
}

// inscription rule
// 1. public article
// 2. after deploy time
// 3. createdtime == updatedtime or updatedtime == 0
// 4. match title (upper case)
// 5. match tags (option)
func (p *Indexer) match_inscription(ft *inscription.FT, item *inscription.MoraFTItem) bool {

	if item.Protocol.Op != "mint" {
		return false
	}

	if item.Ticker != ft.Ticker {
		return false
	}

	if item.Protocol.Amt != ft.Metadata.Protocol.Lim {
		return false
	}

	// must after deploy time
	if item.Created <= ft.DeployTime {
		return false
	}

	return true
}

// func mft_ft_key(aitem *inscription.MoraFTItem) []byte {
// 	key := fmt.Sprintf("%d/%s", aitem.Created, aitem.Ticker)
// 	return []byte(key)
// }

func mft_encode_key(aitem *inscription.MoraFTItem) []byte {
	key := fmt.Sprintf("%d/%s/%s", aitem.Created, aitem.Article, aitem.Canister)
	return []byte(key)
}

// return aid, canister
func mft_decode_key(key string) (string, string, error) {
	strs := strings.Split(key, "/")
	if len(strs) != 3 {
		return "", "", errors.New("not a key")
	}

	return strs[1], strs[2], nil
}

// only use for preload $mora
func match_mora_protocol(content string) *inscription.Mora20Protocol {
	abstract := strings.TrimSpace(content)
	data := utils.RegSub(abstract, `{(.+?)}`)
	if data == "" {
		return nil
	}

	data = fmt.Sprintf("{%s}", data)
	var protocol inscription.Mora20Protocol
	if err := json.Unmarshal([]byte(data), &protocol); err != nil {
		return nil
	}
	return &protocol
}
