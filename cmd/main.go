package main

import (
	"fmt"

	"github.com/thanos-fil/planner-demo-go/internal/planner"
)

func main() {
	net := planner.MockExampleNetwork()

	cfg := planner.UserCfg{}

	clusters := planner.CreateCandidateClusters(net, cfg)
	fmt.Printf("Created %d clusters\n", len(clusters))
}
