package controllers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dstarapp/gomorainsc/http/request"
	"github.com/dstarapp/gomorainsc/http/response"
	"github.com/dstarapp/gomorainsc/http/services"
	"github.com/dstarapp/gomorainsc/inscription"
	"github.com/kataras/iris/v12"
)

type FtController struct {
	Ctx iris.Context
}

func (p *FtController) PostTick() response.BaseResp {
	var req request.FtTickReq
	if err := p.Ctx.ReadJSON(&req); err != nil {
		return response.NewErrResp(err)
	}

	indexer := services.GetIndexer()
	ft := indexer.GetFt(req.Tick)
	if ft == nil {
		return response.NewErrResp(fmt.Errorf("%s is not exist", req.Tick))
	}
	return response.NewSuccessResp(ft)
}

func (p *FtController) PostList() response.BaseResp {
	var req request.FtListReq
	if err := p.Ctx.ReadJSON(&req); err != nil {
		return response.NewErrResp(err)
	}

	indexer := services.GetIndexer()
	fts := indexer.GetSortFts()

	hasmore := false
	total := 0
	page, size := checkPageSize(req.Page, req.Size)
	start := (page - 1) * size
	res := make(inscription.FTS, 0)
	for _, ft := range fts {
		if req.Search != "" && !strings.Contains(ft.Ticker, req.Search) {
			continue
		}
		if req.Status > 0 {
			//select progress
			total := ft.ConfirmCount + ft.UnconfirmCount
			if req.Status == 1 && total >= int(ft.MaxItem) {
				continue
			}
			//select completed
			if req.Status == 2 && total < int(ft.MaxItem) {
				continue
			}
		}
		if total >= start && total < start+size {
			res = append(res, ft)
		}
		total = total + 1
	}
	if total >= start+size {
		hasmore = true
	}

	resp := response.FtListResp{
		Data:    res,
		Total:   total,
		Hasmore: hasmore,
	}
	return response.NewSuccessResp(resp)
}

func (p *FtController) PostListMinters() response.BaseResp {
	var req request.FtListMinterReq
	if err := p.Ctx.ReadJSON(&req); err != nil {
		return response.NewErrResp(err)
	}

	index := services.GetIndexer()
	ft := index.GetFt(req.Tick)
	if ft == nil {
		return response.NewErrResp(fmt.Errorf("%s is not exist", req.Tick))
	}

	if req.Verify && !ft.Done {
		return response.NewSuccessResp([]any{})
	}

	res, err := services.GetTickAllMinters(ft.Ticker, req.Verify)
	if err != nil {
		return response.NewErrResp(err)
	}

	datas := make([]response.MinterInfo, 0)
	for _, v := range res {
		datas = append(datas, v)
	}

	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Count > datas[j].Count
	})
	return response.NewSuccessResp(datas)
}

func checkPageSize(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}
