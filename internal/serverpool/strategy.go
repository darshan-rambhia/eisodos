package serverpool

import "log"

// LBStrategy represents the type of load balancing strategy
type LBStrategy int

const (
	RoundRobin LBStrategy = iota
	LeastConnected
)

// ParseStrategy converts a string strategy to LBStrategy
func ParseStrategy(strategy string) LBStrategy {
	switch strategy {
	case "round-robin":
		return RoundRobin
	case "least-connected":
		return LeastConnected
	default:
		log.Printf("Unknown strategy %s, defaulting to round-robin", strategy)
		return RoundRobin
	}
}
