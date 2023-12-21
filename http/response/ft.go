package response

import "github.com/dstarapp/gomorainsc/inscription"

type BaseResp struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
	Data   any    `json:"data,omitempty"`
}

type FtListResp struct {
	Data    []*inscription.FT `json:"data"`
	Total   int               `json:"total"`
	Hasmore bool              `json:"hasmore"`
}

type FtItemListResp struct {
	Data    []*inscription.MoraFTItem `json:"data"`
	Total   int                       `json:"total"`
	Hasmore bool                      `json:"hasmore"`
}

func NewErrResp(err error) BaseResp {
	return BaseResp{
		Msg:    err.Error(),
		Status: -1,
	}
}

func NewSuccessResp(data any) BaseResp {
	return BaseResp{
		Msg:    "success",
		Status: 0,
		Data:   data,
	}
}
