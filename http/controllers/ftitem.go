package controllers

import (
	"fmt"

	"github.com/dstarapp/gomorainsc/http/request"
	"github.com/dstarapp/gomorainsc/http/response"
	"github.com/dstarapp/gomorainsc/http/services"
	"github.com/dstarapp/gomorainsc/inscription"
	"github.com/kataras/iris/v12"
)

type FtItemController struct {
	Ctx iris.Context
}

func (p *FtItemController) PostList() response.BaseResp {
	var req request.FtItemListReq
	if err := p.Ctx.ReadJSON(&req); err != nil {
		return response.NewErrResp(err)
	}

	indexer := services.GetIndexer()
	ft := indexer.GetFt(req.Tick)
	if ft == nil {
		return response.NewErrResp(fmt.Errorf("tick %s not exist", req.Tick))
	}

	hasmore := false
	total := 0
	page, size := checkPageSize(req.Page, req.Size)
	start := (page - 1) * size
	res := make([]*inscription.MoraFTItem, 0)

	indexer.ScanFtItem(ft.Ticker, func(item *inscription.MoraFTItem) bool {
		if total >= start && total < start+size {
			res = append(res, item)
		}
		total = total + 1
		return len(res) < size
	})

	if ft.ConfirmCount+ft.UnconfirmCount > start+size {
		hasmore = true
	}

	resp := response.FtItemListResp{
		Data:    res,
		Total:   ft.ConfirmCount + ft.UnconfirmCount,
		Hasmore: hasmore,
	}
	return response.NewSuccessResp(resp)
}
