package services

import (
	"errors"
	"fmt"

	"github.com/dstarapp/gomorainsc/http/response"
	"github.com/dstarapp/gomorainsc/inscription"
)

func GetTickAllMinters(tick string, verify bool) (map[string]response.MinterInfo, error) {

	index := GetIndexer()

	ft := index.GetFt(tick)
	if ft == nil {
		return nil, errors.New("tick not exist")
	}

	total := 0
	minters := make(map[string]response.MinterInfo)
	err := index.ScanFtItem(tick, func(item *inscription.MoraFTItem) bool {
		if verify {
			if total >= int(ft.MaxItem) {
				return false
			}
			if !item.Verify {
				return true
			}
			total++
		}
		str := fmt.Sprintf("%s/%s", item.Canister, item.Owner)
		if val, ok := minters[str]; ok {
			val.Count++
			minters[str] = val
		} else {
			minters[str] = response.MinterInfo{
				Canister: item.Canister,
				Owner:    item.Owner,
				Count:    1,
			}
		}
		return true
	})
	return minters, err
}
