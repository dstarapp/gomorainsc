package request

type FtTickReq struct {
	Tick string `json:"tick"`
}

type FtListReq struct {
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Status int    `json:"status"`
	Search string `json:"search"`
}

type FtItemListReq struct {
	Tick string `json:"tick"`
	Page int    `json:"page"`
	Size int    `json:"size"`
}
