package selector

import "context"

type SelectMode int

const (
	RandomSelect SelectMode = iota
	RoundRobin
	WeightedRoundRobin
	ConsistentHash
)

type Selector interface {
	Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string
	UpdateServer(servers map[string]string)
}

func NewSelector(selectMode SelectMode, servers map[string]string) Selector {
	switch selectMode {
	case RandomSelect:
		return newRandomSelector(servers)
	case RoundRobin:
		return newRoundRobinSelector(servers)
	case WeightedRoundRobin:
		return newWeightedRoundRobinSelector(servers)
	case ConsistentHash:
		return newConsistentHashSelector(servers)
	default:
		return newRandomSelector(servers)
	}
}
