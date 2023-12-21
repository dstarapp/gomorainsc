package moraprotocol

import (
	"errors"
	"time"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/principal"
	"github.com/dstarapp/gomorainsc/canisters/dao"
)

const (
	MORA_DAO_CANISTER_ID = "53i5d-faaaa-aaaan-qda6a-cai"
)

func GetAllPlantCanisterIds() ([]principal.Principal, error) {
	actor, err := GetAnonyDao()
	if err != nil {
		return nil, err
	}
	pids, err := actor.QueryCanisterPids()
	if err != nil {
		return nil, err
	}
	if pids == nil {
		return nil, errors.New("can not get planets")
	}
	return *pids, nil
}

func GetAnonyDao() (*dao.Agent, error) {
	id, _ := principal.Decode(MORA_DAO_CANISTER_ID)
	return dao.NewAgent(id, agent.Config{
		IngressExpiry: 300 * time.Second,
	})
}
