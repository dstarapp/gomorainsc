package moraprotocol

import (
	"time"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/identity"
	"github.com/aviate-labs/agent-go/principal"
	"github.com/dstarapp/gomorainsc/canisters/planet"
)

func GetAnonyPlanet(canister principal.Principal) (*planet.Agent, error) {
	return planet.NewAgent(canister, agent.Config{
		IngressExpiry: 300 * time.Second,
	})
}

func GetPlanetWithIdentity(canister principal.Principal, identity identity.Identity) (*planet.Agent, error) {
	return planet.NewAgent(canister, agent.Config{
		IngressExpiry: 300 * time.Second,
		Identity:      identity,
	})
}
