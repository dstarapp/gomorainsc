// Package dao provides a client for the "dao" canister.
// Do NOT edit this file. It was automatically generated by https://github.com/aviate-labs/agent-go.
package dao

import (
	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/aviate-labs/agent-go/principal"
)

type CreatePlanetSetting struct {
	Avatar string              `ic:"avatar"`
	Code   string              `ic:"code"`
	Desc   string              `ic:"desc"`
	Name   string              `ic:"name"`
	Owner  principal.Principal `ic:"owner"`
}

type CreatePlanetResp struct {
	Err *string `ic:"Err,variant"`
	Ok  *struct {
		Id principal.Principal `ic:"id"`
	} `ic:"Ok,variant"`
}

type CanisterInfo struct {
	Id          principal.Principal `ic:"id"`
	InitArgs    []byte              `ic:"initArgs"`
	LaunchTrail principal.Principal `ic:"launchTrail"`
	ModuleHash  []byte              `ic:"moduleHash"`
	Owner       principal.Principal `ic:"owner"`
}

// Agent is a client for the "dao" canister.
type Agent struct {
	a          *agent.Agent
	canisterId principal.Principal
}

// NewAgent creates a new agent for the "dao" canister.
func NewAgent(canisterId principal.Principal, config agent.Config) (*Agent, error) {
	a, err := agent.New(config)
	if err != nil {
		return nil, err
	}
	return &Agent{
		a:          a,
		canisterId: canisterId,
	}, nil
}

// CanisterAccount calls the "canisterAccount" method on the "dao" canister.
func (a Agent) CanisterAccount() (*string, error) {
	var r0 string
	if err := a.a.Query(
		a.canisterId,
		"canisterAccount",
		[]any{},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// CreatePlanet calls the "createPlanet" method on the "dao" canister.
func (a Agent) CreatePlanet(arg0 CreatePlanetSetting) (*CreatePlanetResp, error) {
	var r0 CreatePlanetResp
	if err := a.a.Call(
		a.canisterId,
		"createPlanet",
		[]any{arg0},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// InitTrail calls the "initTrail" method on the "dao" canister.
func (a Agent) InitTrail(arg0 idl.Int) (*idl.Int, error) {
	var r0 idl.Int
	if err := a.a.Call(
		a.canisterId,
		"initTrail",
		[]any{arg0},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// QueryAgreePayee calls the "queryAgreePayee" method on the "dao" canister.
func (a Agent) QueryAgreePayee() (*string, error) {
	var r0 string
	if err := a.a.Query(
		a.canisterId,
		"queryAgreePayee",
		[]any{},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// QueryCanisterCount calls the "queryCanisterCount" method on the "dao" canister.
func (a Agent) QueryCanisterCount() (*idl.Int, error) {
	var r0 idl.Int
	if err := a.a.Query(
		a.canisterId,
		"queryCanisterCount",
		[]any{},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// QueryCanisterIds calls the "queryCanisterIds" method on the "dao" canister.
func (a Agent) QueryCanisterIds() (*[]string, error) {
	var r0 []string
	if err := a.a.Query(
		a.canisterId,
		"queryCanisterIds",
		[]any{},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// QueryCanisterPids calls the "queryCanisterPids" method on the "dao" canister.
func (a Agent) QueryCanisterPids() (*[]principal.Principal, error) {
	var r0 []principal.Principal
	if err := a.a.Query(
		a.canisterId,
		"queryCanisterPids",
		[]any{},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// QueryCanisters calls the "queryCanisters" method on the "dao" canister.
func (a Agent) QueryCanisters() (*[]CanisterInfo, error) {
	var r0 []CanisterInfo
	if err := a.a.Query(
		a.canisterId,
		"queryCanisters",
		[]any{},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// QueryTrailPids calls the "queryTrailPids" method on the "dao" canister.
func (a Agent) QueryTrailPids() (*[]principal.Principal, error) {
	var r0 []principal.Principal
	if err := a.a.Query(
		a.canisterId,
		"queryTrailPids",
		[]any{},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// QueryTrailWasmHash calls the "queryTrailWasmHash" method on the "dao" canister.
func (a Agent) QueryTrailWasmHash() (*string, error) {
	var r0 string
	if err := a.a.Query(
		a.canisterId,
		"queryTrailWasmHash",
		[]any{},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// QueryWasmHash calls the "queryWasmHash" method on the "dao" canister.
func (a Agent) QueryWasmHash() (*string, error) {
	var r0 string
	if err := a.a.Query(
		a.canisterId,
		"queryWasmHash",
		[]any{},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}

// SetAgreePayee calls the "setAgreePayee" method on the "dao" canister.
func (a Agent) SetAgreePayee(arg0 []byte) error {
	if err := a.a.Call(
		a.canisterId,
		"setAgreePayee",
		[]any{arg0},
		[]any{},
	); err != nil {
		return err
	}
	return nil
}

// SetInviter calls the "setInviter" method on the "dao" canister.
func (a Agent) SetInviter(arg0 principal.Principal) error {
	if err := a.a.Call(
		a.canisterId,
		"setInviter",
		[]any{arg0},
		[]any{},
	); err != nil {
		return err
	}
	return nil
}

// SetOwner calls the "setOwner" method on the "dao" canister.
func (a Agent) SetOwner(arg0 principal.Principal) error {
	if err := a.a.Call(
		a.canisterId,
		"setOwner",
		[]any{arg0},
		[]any{},
	); err != nil {
		return err
	}
	return nil
}

// SetTrailWasm calls the "setTrailWasm" method on the "dao" canister.
func (a Agent) SetTrailWasm(arg0 []byte) error {
	if err := a.a.Call(
		a.canisterId,
		"setTrailWasm",
		[]any{arg0},
		[]any{},
	); err != nil {
		return err
	}
	return nil
}

// SetUserRouter calls the "setUserRouter" method on the "dao" canister.
func (a Agent) SetUserRouter(arg0 principal.Principal) error {
	if err := a.a.Call(
		a.canisterId,
		"setUserRouter",
		[]any{arg0},
		[]any{},
	); err != nil {
		return err
	}
	return nil
}

// SetWasm calls the "setWasm" method on the "dao" canister.
func (a Agent) SetWasm(arg0 []byte) error {
	if err := a.a.Call(
		a.canisterId,
		"setWasm",
		[]any{arg0},
		[]any{},
	); err != nil {
		return err
	}
	return nil
}

// UpgradePlanet calls the "upgradePlanet" method on the "dao" canister.
func (a Agent) UpgradePlanet(arg0 principal.Principal) (*bool, error) {
	var r0 bool
	if err := a.a.Call(
		a.canisterId,
		"upgradePlanet",
		[]any{arg0},
		[]any{&r0},
	); err != nil {
		return nil, err
	}
	return &r0, nil
}