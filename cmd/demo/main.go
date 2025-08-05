package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/thanos-fil/planner-demo-go/internal/graph"
	"github.com/thanos-fil/planner-demo-go/internal/neighbourhood"
	"github.com/thanos-fil/planner-demo-go/internal/planner"
)

// PipeID is an alias for graph.ID for consistency
type PipeID = graph.ID

func main() {
	fmt.Println("=== Pipes Network Risk Assessment and Planning Demo ===")
	fmt.Println()

	// Create a realistic large network
	network := createLargePipeNetwork(1000) // 1000 pipes for demo
	fmt.Printf("Created network with %d pipes\n", len(network.Nodes))

	// Configuration for planning
	cfg := planner.UserCfg{
		TargetCount:         10,    // Want 10 clusters/projects
		MaxLength:           150.0, // Maximum length per project
		OvershootFactor:     1.2,   // Allow 20% overshoot
		LongestPathFraction: 0.8,   // Use 80% of max length for longest paths
	}

	// Analyze the network risk
	fmt.Println("\n=== Network Risk Analysis ===")
	analyzeNetworkRisk(network)

	// Find high-risk neighborhoods
	fmt.Println("\n=== High-Risk Neighborhoods Analysis ===")
	highRiskAreas := findHighRiskNeighborhoods(network, cfg.MaxLength)

	// Create candidate clusters/projects
	fmt.Println("\n=== Project Planning ===")
	projects := createProjectsFromHighRiskAreas(network, highRiskAreas, cfg)

	// Display results
	fmt.Println("\n=== Recommended Projects ===")
	displayProjects(network, projects)

	// Calculate potential risk reduction
	fmt.Println("\n=== Risk Reduction Analysis ===")
	calculateRiskReduction(network, projects)
}

// createLargePipeNetwork creates a realistic pipe network with varying LoF values
func createLargePipeNetwork(numPipes int) *graph.Network {
	net := graph.New()
	rand.Seed(time.Now().UnixNano())

	// Create pipes with realistic characteristics
	for i := 0; i < numPipes; i++ {
		// Simulate realistic pipe properties
		length := 10.0 + rand.Float64()*90.0 // 10-100m pipes

		// LoF (Likelihood of Failure) - higher values mean more likely to fail
		// Simulate different risk profiles:
		// - Most pipes: low risk (0.1-0.3)
		// - Some pipes: medium risk (0.3-0.6)
		// - Few pipes: high risk (0.6-0.9)
		var lof float64
		riskCategory := rand.Float64()
		switch {
		case riskCategory < 0.7: // 70% low risk
			lof = 0.1 + rand.Float64()*0.2
		case riskCategory < 0.9: // 20% medium risk
			lof = 0.3 + rand.Float64()*0.3
		default: // 10% high risk
			lof = 0.6 + rand.Float64()*0.3
		}

		net.AddNode(&graph.Node{
			ID:     graph.ID(i),
			Score:  lof,
			Length: length,
		})
	}

	// Create realistic network topology (grid-like with some randomness)
	gridSize := int(math.Sqrt(float64(numPipes)))

	// Add grid connections
	for i := 0; i < numPipes; i++ {
		row := i / gridSize
		col := i % gridSize

		// Connect to right neighbor
		if col < gridSize-1 {
			rightNeighbor := row*gridSize + col + 1
			if rightNeighbor < numPipes {
				net.AddUndirectedEdge(graph.ID(i), graph.ID(rightNeighbor))
			}
		}

		// Connect to bottom neighbor
		if row < gridSize-1 {
			bottomNeighbor := (row+1)*gridSize + col
			if bottomNeighbor < numPipes {
				net.AddUndirectedEdge(graph.ID(i), graph.ID(bottomNeighbor))
			}
		}

		// Add some random cross-connections (15% chance)
		if rand.Float64() < 0.15 {
			randomNeighbor := rand.Intn(numPipes)
			if randomNeighbor != i {
				net.AddUndirectedEdge(graph.ID(i), graph.ID(randomNeighbor))
			}
		}
	}

	return net
}

// analyzeNetworkRisk provides an overview of the network's risk profile
func analyzeNetworkRisk(net *graph.Network) {
	var totalRisk, totalLength float64
	var highRiskPipes, mediumRiskPipes, lowRiskPipes int

	for _, node := range net.Nodes {
		totalRisk += node.Score
		totalLength += node.Length

		switch {
		case node.Score >= 0.6:
			highRiskPipes++
		case node.Score >= 0.3:
			mediumRiskPipes++
		default:
			lowRiskPipes++
		}
	}

	avgRisk := totalRisk / float64(len(net.Nodes))

	fmt.Printf("Network Overview:\n")
	fmt.Printf("  Total Pipes: %d\n", len(net.Nodes))
	fmt.Printf("  Total Length: %.1f km\n", totalLength/1000)
	fmt.Printf("  Average Risk (LoF): %.3f\n", avgRisk)
	fmt.Printf("  Risk Distribution:\n")
	fmt.Printf("    High Risk (≥0.6): %d pipes (%.1f%%)\n", highRiskPipes, float64(highRiskPipes)/float64(len(net.Nodes))*100)
	fmt.Printf("    Medium Risk (0.3-0.6): %d pipes (%.1f%%)\n", mediumRiskPipes, float64(mediumRiskPipes)/float64(len(net.Nodes))*100)
	fmt.Printf("    Low Risk (<0.3): %d pipes (%.1f%%)\n", lowRiskPipes, float64(lowRiskPipes)/float64(len(net.Nodes))*100)
}

// findHighRiskNeighborhoods identifies areas with concentrated high-risk pipes
func findHighRiskNeighborhoods(net *graph.Network, maxProjectLength float64) []PipeID {
	var highRiskSeeds []PipeID

	// Find pipes with high individual risk that could be seeds for projects
	for id, node := range net.Nodes {
		if node.Score >= 0.5 { // High risk threshold
			highRiskSeeds = append(highRiskSeeds, id)
		}
	}

	fmt.Printf("Found %d high-risk seed pipes\n", len(highRiskSeeds))

	// For demonstration, let's analyze neighborhoods around the top 10 highest risk pipes
	if len(highRiskSeeds) > 10 {
		highRiskSeeds = highRiskSeeds[:10]
	}

	return highRiskSeeds
}

// createProjectsFromHighRiskAreas creates project clusters around high-risk areas
func createProjectsFromHighRiskAreas(net *graph.Network, seeds []PipeID, cfg planner.UserCfg) []planner.Cluster {
	var projects []planner.Cluster
	usedPipes := make(map[PipeID]bool)

	for i, seed := range seeds {
		if usedPipes[seed] {
			continue // Skip if already included in another project
		}

		// Use neighbourhood function to find pipes within project budget
		neighbors := neighbourhood.Neighbourhood(net, seed, cfg.MaxLength)

		// Create a cluster from the neighborhood
		var clusterNodes []graph.ID
		var totalRisk float64

		for pipeID := range neighbors {
			if !usedPipes[pipeID] {
				clusterNodes = append(clusterNodes, graph.ID(pipeID))
				totalRisk += net.Nodes[graph.ID(pipeID)].Score
				usedPipes[pipeID] = true
			}
		}

		if len(clusterNodes) > 0 {
			projects = append(projects, planner.Cluster{
				ID:    i + 1,
				Nodes: clusterNodes,
				Score: totalRisk, // Total risk in this cluster
			})
		}

		if len(projects) >= cfg.TargetCount {
			break // Reached target number of projects
		}
	}

	return projects
}

// displayProjects shows the recommended projects with their characteristics
func displayProjects(net *graph.Network, projects []planner.Cluster) {
	for _, project := range projects {
		var totalLength float64
		var avgRisk float64

		for _, nodeID := range project.Nodes {
			node := net.Nodes[nodeID]
			totalLength += node.Length
		}

		if len(project.Nodes) > 0 {
			avgRisk = project.Score / float64(len(project.Nodes))
		}

		fmt.Printf("Project %d:\n", project.ID)
		fmt.Printf("  Pipes Count: %d\n", len(project.Nodes))
		fmt.Printf("  Total Length: %.1f m\n", totalLength)
		fmt.Printf("  Total Risk: %.3f\n", project.Score)
		fmt.Printf("  Average Risk: %.3f\n", avgRisk)
		fmt.Printf("  Risk Density: %.3f (risk/km)\n", project.Score/(totalLength/1000))
		fmt.Printf("  Sample Pipes: %v\n", project.Nodes[:min(5, len(project.Nodes))])
		fmt.Println()
	}
}

// calculateRiskReduction estimates the risk reduction from implementing projects
func calculateRiskReduction(net *graph.Network, projects []planner.Cluster) {
	var totalRiskBeforeProjects float64
	var totalRiskInProjects float64

	for _, node := range net.Nodes {
		totalRiskBeforeProjects += node.Score
	}

	for _, project := range projects {
		totalRiskInProjects += project.Score
	}

	// Assume projects reduce risk by 70% for replaced/renovated pipes
	riskReductionFactor := 0.7
	riskReduction := totalRiskInProjects * riskReductionFactor

	percentageReduction := (riskReduction / totalRiskBeforeProjects) * 100

	fmt.Printf("Risk Reduction Analysis:\n")
	fmt.Printf("  Total Network Risk: %.2f\n", totalRiskBeforeProjects)
	fmt.Printf("  Risk in Selected Projects: %.2f\n", totalRiskInProjects)
	fmt.Printf("  Potential Risk Reduction: %.2f (%.1f%% of total network)\n", riskReduction, percentageReduction)
	fmt.Printf("  Remaining Network Risk: %.2f\n", totalRiskBeforeProjects-riskReduction)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
