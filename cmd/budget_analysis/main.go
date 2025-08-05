package main

import (
	"fmt"

	"github.com/thanos-fil/planner-demo-go/internal/graph"
	"github.com/thanos-fil/planner-demo-go/internal/neighbourhood"
)

func main() {
	fmt.Println("=== Budget Impact Analysis ===")
	fmt.Println()

	// Create the same network as in the detailed demo
	net := createDemoNetwork()
	displayNetworkInfo(net)

	// Test different budget scenarios
	startNode := graph.ID(0)
	budgets := []float64{3.0, 5.0, 8.0, 12.0}

	fmt.Printf("\n=== Budget Scenario Analysis (Start Node: %d) ===\n", startNode)

	for _, budget := range budgets {
		fmt.Printf("\n--- Budget: %.1f ---\n", budget)
		result := neighbourhood.Neighbourhood(net, startNode, budget)
		analyzeScenario(net, result, budget, startNode)
	}

	// Show the impact of different starting nodes
	fmt.Printf("\n=== Starting Node Impact Analysis (Budget: 8.0) ===\n")
	budget := 8.0
	startNodes := []graph.ID{0, 1, 3, 5}

	for _, start := range startNodes {
		fmt.Printf("\n--- Starting from Node %d ---\n", start)
		result := neighbourhood.Neighbourhood(net, start, budget)
		analyzeScenario(net, result, budget, start)
	}
}

func createDemoNetwork() *graph.Network {
	net := graph.New()

	// Same network as visualization demo
	nodes := []*graph.Node{
		{ID: 0, Score: 0.8, Length: 2.0}, // High risk, short
		{ID: 1, Score: 0.3, Length: 1.5}, // Low risk, short
		{ID: 2, Score: 0.6, Length: 3.0}, // Medium risk, medium
		{ID: 3, Score: 0.9, Length: 1.0}, // Very high risk, very short
		{ID: 4, Score: 0.2, Length: 4.0}, // Very low risk, long
		{ID: 5, Score: 0.7, Length: 2.5}, // High risk, medium
		{ID: 6, Score: 0.4, Length: 1.8}, // Medium risk, short
	}

	for _, node := range nodes {
		net.AddNode(node)
	}

	// Same topology: 1-0-2-5-6 with 3 and 4 connected to 0
	net.AddUndirectedEdge(0, 1)
	net.AddUndirectedEdge(0, 2)
	net.AddUndirectedEdge(0, 3)
	net.AddUndirectedEdge(0, 4)
	net.AddUndirectedEdge(2, 5)
	net.AddUndirectedEdge(5, 6)

	return net
}

func displayNetworkInfo(net *graph.Network) {
	fmt.Println("Network Topology:")
	fmt.Println("        3(0.9/1.0)")
	fmt.Println("         |")
	fmt.Println("    1(0.3/1.5)---0(0.8/2.0)---2(0.6/3.0)---5(0.7/2.5)---6(0.4/1.8)")
	fmt.Println("                  |")
	fmt.Println("               4(0.2/4.0)")
	fmt.Println()
	fmt.Println("Format: NodeID(LoF/Length)")
}

func analyzeScenario(net *graph.Network, result map[neighbourhood.PipeID]struct{}, budget float64, startNode graph.ID) {
	var totalRisk, totalLength float64
	var visitedNodes []int

	for nodeID := range result {
		node := net.Nodes[graph.ID(nodeID)]
		totalRisk += node.Score
		totalLength += node.Length
		visitedNodes = append(visitedNodes, int(nodeID))
	}

	// Calculate project metrics
	estimatedCost := totalLength*500 + 10000
	avgCoF := 27500.0
	bre := totalRisk * avgCoF
	roi := bre / estimatedCost

	fmt.Printf("Results:\n")
	fmt.Printf("  Nodes included: %v\n", visitedNodes)
	fmt.Printf("  Total pipes: %d\n", len(result))
	fmt.Printf("  Total length: %.1fm\n", totalLength)
	fmt.Printf("  Total risk: %.2f\n", totalRisk)
	fmt.Printf("  Avg risk: %.3f\n", totalRisk/float64(len(result)))
	fmt.Printf("  Cost: €%.0f\n", estimatedCost)
	fmt.Printf("  ROI: %.2f\n", roi)

	// Show which high-risk nodes are included
	highRiskIncluded := 0
	for nodeID := range result {
		if net.Nodes[graph.ID(nodeID)].Score >= 0.7 {
			highRiskIncluded++
		}
	}
	fmt.Printf("  High-risk nodes (≥0.7): %d/%d\n", highRiskIncluded, countHighRiskNodes(net))

	// Visual representation
	fmt.Printf("  Visual: ")
	for i := 0; i < 7; i++ {
		if _, included := result[neighbourhood.PipeID(i)]; included {
			if graph.ID(i) == startNode {
				fmt.Printf("[%d*] ", i)
			} else {
				fmt.Printf("[%d] ", i)
			}
		} else {
			fmt.Printf("(%d) ", i)
		}
	}
	fmt.Printf("  (* = start node)\n")
}

func countHighRiskNodes(net *graph.Network) int {
	count := 0
	for _, node := range net.Nodes {
		if node.Score >= 0.7 {
			count++
		}
	}
	return count
}
