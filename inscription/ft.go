package inscription

import (
	"fmt"
	"strings"
)

type FT struct {
	Ticker         string      `json:"ticker"`
	Description    string      `json:"description,omitempty"`
	Image          string      `json:"image,omitempty"`
	Limit          uint64      `json:"limit"`
	Max            uint64      `json:"max"`
	MaxItem        uint64      `json:"max_item,omitempty"`
	DeployTime     uint64      `json:"deploy_time"`
	Deployer       string      `json:"deployer,omitempty"`
	Canister       string      `json:"canister,omitempty"`
	Index          string      `json:"index,omitempty"`
	FullCount      int         `json:"full_count"` // after full, run spider count
	ConfirmCount   int         `json:"confirm_count"`
	UnconfirmCount int         `json:"unconfirm_count"`
	Minters        int         `json:"minters"`
	Done           bool        `json:"done"`
	LastUpdated    uint64      `json:"last_updated"`
	Metadata       *MoraFTItem `json:"metadata"`
}

type MoraFTItem struct {
	ID       uint64          `json:"id,omitempty"`
	Ticker   string          `json:"ticker"`
	Owner    string          `json:"owner"`
	Canister string          `json:"canister"`
	Article  string          `json:"article"`
	Title    string          `json:"title"`
	Thumb    string          `json:"thumb"`
	Created  uint64          `json:"created"`
	Updated  uint64          `json:"updated"`
	Content  string          `json:"content"`
	Tag      string          `json:"tag"`
	Verify   bool            `json:"verify,omitempty"`
	MintTx   string          `json:"mint_tx,omitempty"`
	Protocol *Mora20Protocol `json:"protocol,omitempty"`
}

type Mora20Protocol struct {
	P    string `json:"p"`
	Op   string `json:"op"`
	Tick string `json:"tick"`
	Amt  string `json:"amt,omitempty"`
	Max  string `json:"max,omitempty"`
	Lim  string `json:"lim,omitempty"`
}

type FTS []*FT

func (p FTS) Len() int {
	return len(p)
}

func (p FTS) Less(i, j int) bool {
	if p[i].Done {
		if !p[j].Done {
			return true
		}
	} else if p[j].Done {
		return false
	}
	return p[i].DeployTime < p[j].DeployTime
}

func (p FTS) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type FTItems []*MoraFTItem

func (p FTItems) Len() int {
	return len(p)
}

func (p FTItems) Less(i, j int) bool {
	k1 := fmt.Sprintf("%d/%s/%s", p[i].Created, p[i].Article, p[i].Canister)
	k2 := fmt.Sprintf("%d/%s/%s", p[j].Created, p[j].Article, p[j].Canister)
	return strings.Compare(k1, k2) < 0
}

func (p FTItems) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
