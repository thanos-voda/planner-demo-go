package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/thanos-fil/planner-demo-go/internal/graph"
	"github.com/thanos-fil/planner-demo-go/internal/neighbourhood"
	"github.com/thanos-fil/planner-demo-go/internal/planner"
)

// PipeID is an alias for graph.ID for consistency
type PipeID = graph.ID

// ProjectMetrics contains various metrics for evaluating projects
type ProjectMetrics struct {
	planner.Cluster
	TotalLength    float64
	AvgRisk        float64
	RiskDensity    float64 // Risk per km
	LoFLengthRatio float64 // LoF divided by Length
	EstimatedCost  float64 // Cost to implement project
	BRE            float64 // Business Risk Exposure (LoF × CoF)
	ROI            float64 // Return on Investment (BRE/Cost)
	Priority       float64 // Overall priority score
}

// CostOfFailure represents the cost of failure for different pipe scenarios
type CostOfFailure struct {
	MinorLeak   float64 // €1,000 - €5,000
	MajorLeak   float64 // €5,000 - €25,000
	Burst       float64 // €25,000 - €100,000
	ServiceDist float64 // €10,000 - €50,000 for service disruption
}

func main() {
	fmt.Println("=== Advanced Pipes Network Risk Assessment and ROI Planning ===")
	fmt.Println()

	// Create a realistic large network
	networkSize := 2000 // Increase to 2000 pipes for more realistic scenario
	network := createLargePipeNetwork(networkSize)
	fmt.Printf("Created network with %d pipes\n", len(network.Nodes))

	// Configuration for planning
	cfg := planner.UserCfg{
		TargetCount:         15,    // Want 15 clusters/projects
		MaxLength:           200.0, // Maximum length per project (200m)
		OvershootFactor:     1.2,   // Allow 20% overshoot
		LongestPathFraction: 0.8,   // Use 80% of max length for longest paths
	}

	// Cost parameters
	costOfFailure := CostOfFailure{
		MinorLeak:   3000,  // €3,000 average
		MajorLeak:   15000, // €15,000 average
		Burst:       62500, // €62,500 average
		ServiceDist: 30000, // €30,000 average
	}

	// Analyze the network risk
	fmt.Println("\n=== Network Risk Analysis ===")
	analyzeNetworkRisk(network)

	// Find high-risk neighborhoods using different strategies
	fmt.Println("\n=== High-Risk Identification Strategies ===")
	strategies := map[string][]PipeID{
		"High LoF":          findHighLoFPipes(network),
		"High Risk Density": findHighRiskDensityAreas(network, cfg.MaxLength),
		"High LoF/Length":   findHighLoFLengthRatioPipes(network),
	}

	// Create projects using different strategies
	fmt.Println("\n=== Project Planning with Different Strategies ===")
	allProjects := make(map[string][]ProjectMetrics)

	for strategy, seeds := range strategies {
		fmt.Printf("\n--- Strategy: %s ---\n", strategy)
		projects := createAdvancedProjects(network, seeds, cfg, costOfFailure)

		// Sort projects by priority (ROI)
		sort.Slice(projects, func(i, j int) bool {
			return projects[i].ROI > projects[j].ROI
		})

		// Keep top projects within budget
		if len(projects) > cfg.TargetCount {
			projects = projects[:cfg.TargetCount]
		}

		allProjects[strategy] = projects
		displayAdvancedProjects(projects, strategy)
	}

	// Compare strategies
	fmt.Println("\n=== Strategy Comparison ===")
	compareStrategies(allProjects)

	// Recommend the best strategy
	fmt.Println("\n=== Recommendation ===")
	recommendBestStrategy(allProjects)
}

// createLargePipeNetwork creates a realistic pipe network with varying LoF values
func createLargePipeNetwork(numPipes int) *graph.Network {
	net := graph.New()
	rand.Seed(time.Now().UnixNano())

	// Create pipes with realistic characteristics
	for i := 0; i < numPipes; i++ {
		// Simulate realistic pipe properties
		length := 5.0 + rand.Float64()*45.0 // 5-50m pipes (more realistic urban segments)

		// LoF (Likelihood of Failure) with more realistic distribution
		var lof float64
		age := rand.Float64()         // Simulate pipe age factor
		material := rand.Float64()    // Simulate material quality
		environment := rand.Float64() // Simulate environmental factors

		// Combine factors for more realistic LoF
		baseLof := 0.1 + age*0.3 + (1-material)*0.2 + environment*0.1
		lof = math.Min(0.9, baseLof) // Cap at 0.9

		net.AddNode(&graph.Node{
			ID:     graph.ID(i),
			Score:  lof,
			Length: length,
		})
	}

	// Create realistic network topology
	gridSize := int(math.Sqrt(float64(numPipes)))

	// Add grid connections with more realistic topology
	for i := 0; i < numPipes; i++ {
		row := i / gridSize
		col := i % gridSize

		// Connect to neighbors with some probability
		neighbors := []int{
			row*gridSize + col + 1, // right
			(row+1)*gridSize + col, // bottom
			row*gridSize + col - 1, // left
			(row-1)*gridSize + col, // top
		}

		for _, neighbor := range neighbors {
			if neighbor >= 0 && neighbor < numPipes && neighbor != i {
				// Add connection with 70% probability for adjacent pipes
				if rand.Float64() < 0.7 {
					net.AddUndirectedEdge(graph.ID(i), graph.ID(neighbor))
				}
			}
		}

		// Add some random cross-connections (10% chance)
		if rand.Float64() < 0.1 {
			randomNeighbor := rand.Intn(numPipes)
			if randomNeighbor != i {
				net.AddUndirectedEdge(graph.ID(i), graph.ID(randomNeighbor))
			}
		}
	}

	return net
}

// analyzeNetworkRisk provides detailed network analysis
func analyzeNetworkRisk(net *graph.Network) {
	var totalRisk, totalLength float64
	var riskCategories [4]int   // [low, medium, high, critical]
	var lengthCategories [3]int // [short, medium, long]

	for _, node := range net.Nodes {
		totalRisk += node.Score
		totalLength += node.Length

		// Risk categories
		switch {
		case node.Score >= 0.7:
			riskCategories[3]++ // critical
		case node.Score >= 0.5:
			riskCategories[2]++ // high
		case node.Score >= 0.3:
			riskCategories[1]++ // medium
		default:
			riskCategories[0]++ // low
		}

		// Length categories
		switch {
		case node.Length >= 30:
			lengthCategories[2]++ // long
		case node.Length >= 15:
			lengthCategories[1]++ // medium
		default:
			lengthCategories[0]++ // short
		}
	}

	avgRisk := totalRisk / float64(len(net.Nodes))
	avgLength := totalLength / float64(len(net.Nodes))

	fmt.Printf("Network Overview:\n")
	fmt.Printf("  Total Pipes: %d\n", len(net.Nodes))
	fmt.Printf("  Total Length: %.2f km\n", totalLength/1000)
	fmt.Printf("  Average Risk (LoF): %.3f\n", avgRisk)
	fmt.Printf("  Average Length: %.1f m\n", avgLength)
	fmt.Printf("\nRisk Distribution:\n")
	fmt.Printf("  Critical Risk (≥0.7): %d pipes (%.1f%%)\n", riskCategories[3], float64(riskCategories[3])/float64(len(net.Nodes))*100)
	fmt.Printf("  High Risk (0.5-0.7): %d pipes (%.1f%%)\n", riskCategories[2], float64(riskCategories[2])/float64(len(net.Nodes))*100)
	fmt.Printf("  Medium Risk (0.3-0.5): %d pipes (%.1f%%)\n", riskCategories[1], float64(riskCategories[1])/float64(len(net.Nodes))*100)
	fmt.Printf("  Low Risk (<0.3): %d pipes (%.1f%%)\n", riskCategories[0], float64(riskCategories[0])/float64(len(net.Nodes))*100)
	fmt.Printf("\nLength Distribution:\n")
	fmt.Printf("  Long Pipes (≥30m): %d pipes (%.1f%%)\n", lengthCategories[2], float64(lengthCategories[2])/float64(len(net.Nodes))*100)
	fmt.Printf("  Medium Pipes (15-30m): %d pipes (%.1f%%)\n", lengthCategories[1], float64(lengthCategories[1])/float64(len(net.Nodes))*100)
	fmt.Printf("  Short Pipes (<15m): %d pipes (%.1f%%)\n", lengthCategories[0], float64(lengthCategories[0])/float64(len(net.Nodes))*100)
}

// findHighLoFPipes finds pipes with highest likelihood of failure
func findHighLoFPipes(net *graph.Network) []PipeID {
	type pipeRisk struct {
		id   PipeID
		risk float64
	}

	var pipes []pipeRisk
	for id, node := range net.Nodes {
		pipes = append(pipes, pipeRisk{id: id, risk: node.Score})
	}

	// Sort by risk (highest first)
	sort.Slice(pipes, func(i, j int) bool {
		return pipes[i].risk > pipes[j].risk
	})

	// Return top 20% highest risk pipes
	numSeeds := len(pipes) / 5
	if numSeeds > 50 {
		numSeeds = 50 // Limit to 50 seeds
	}

	var seeds []PipeID
	for i := 0; i < numSeeds; i++ {
		seeds = append(seeds, pipes[i].id)
	}

	fmt.Printf("Found %d high LoF seed pipes (top %.1f%% by risk)\n", len(seeds), float64(numSeeds)/float64(len(pipes))*100)
	return seeds
}

// findHighRiskDensityAreas finds areas with high concentration of risky pipes
func findHighRiskDensityAreas(net *graph.Network, maxLength float64) []PipeID {
	type areaRisk struct {
		id          PipeID
		riskDensity float64
	}

	var areas []areaRisk

	// For each pipe, calculate risk density in its neighborhood
	for id := range net.Nodes {
		neighbors := neighbourhood.Neighbourhood(net, id, maxLength/2) // Use smaller radius for density calculation

		var totalRisk, totalLength float64
		for neighborID := range neighbors {
			node := net.Nodes[graph.ID(neighborID)]
			totalRisk += node.Score
			totalLength += node.Length
		}

		if totalLength > 0 {
			riskDensity := totalRisk / (totalLength / 1000) // Risk per km
			areas = append(areas, areaRisk{id: id, riskDensity: riskDensity})
		}
	}

	// Sort by risk density
	sort.Slice(areas, func(i, j int) bool {
		return areas[i].riskDensity > areas[j].riskDensity
	})

	// Return top 15% by risk density
	numSeeds := len(areas) * 15 / 100
	if numSeeds > 40 {
		numSeeds = 40
	}

	var seeds []PipeID
	for i := 0; i < numSeeds; i++ {
		seeds = append(seeds, areas[i].id)
	}

	fmt.Printf("Found %d high risk density seed areas\n", len(seeds))
	return seeds
}

// findHighLoFLengthRatioPipes finds pipes with the best LoF/Length ratio
func findHighLoFLengthRatioPipes(net *graph.Network) []PipeID {
	type pipeEfficiency struct {
		id    PipeID
		ratio float64
	}

	var pipes []pipeEfficiency
	for id, node := range net.Nodes {
		ratio := node.Score / (node.Length / 1000) // LoF per km
		pipes = append(pipes, pipeEfficiency{id: id, ratio: ratio})
	}

	// Sort by ratio (highest first)
	sort.Slice(pipes, func(i, j int) bool {
		return pipes[i].ratio > pipes[j].ratio
	})

	// Return top 20%
	numSeeds := len(pipes) / 5
	if numSeeds > 45 {
		numSeeds = 45
	}

	var seeds []PipeID
	for i := 0; i < numSeeds; i++ {
		seeds = append(seeds, pipes[i].id)
	}

	fmt.Printf("Found %d high LoF/Length ratio seed pipes\n", len(seeds))
	return seeds
}

// createAdvancedProjects creates projects with detailed metrics
func createAdvancedProjects(net *graph.Network, seeds []PipeID, cfg planner.UserCfg, cof CostOfFailure) []ProjectMetrics {
	var projects []ProjectMetrics
	usedPipes := make(map[PipeID]bool)

	for i, seed := range seeds {
		if usedPipes[seed] {
			continue
		}

		neighbors := neighbourhood.Neighbourhood(net, seed, cfg.MaxLength)

		var clusterNodes []graph.ID
		var totalRisk, totalLength float64

		for pipeID := range neighbors {
			if !usedPipes[pipeID] {
				node := net.Nodes[graph.ID(pipeID)]
				clusterNodes = append(clusterNodes, graph.ID(pipeID))
				totalRisk += node.Score
				totalLength += node.Length
				usedPipes[pipeID] = true
			}
		}

		if len(clusterNodes) > 0 {
			// Calculate advanced metrics
			avgRisk := totalRisk / float64(len(clusterNodes))
			riskDensity := totalRisk / (totalLength / 1000)
			lofLengthRatio := totalRisk / (totalLength / 1000)

			// Estimate project cost (€500 per meter + fixed costs)
			estimatedCost := totalLength*500 + 10000 // €500/m + €10k fixed

			// Calculate Business Risk Exposure (LoF × Cost of Failure)
			// Assume different failure types based on LoF level
			avgCoF := (cof.MinorLeak + cof.MajorLeak + cof.Burst + cof.ServiceDist) / 4
			bre := totalRisk * avgCoF

			// Calculate ROI (potential savings / investment cost)
			roi := bre / estimatedCost

			// Calculate priority score (weighted combination of factors)
			priority := roi*0.4 + riskDensity*0.3 + avgRisk*0.2 + lofLengthRatio*0.1

			project := ProjectMetrics{
				Cluster: planner.Cluster{
					ID:    i + 1,
					Nodes: clusterNodes,
					Score: totalRisk,
				},
				TotalLength:    totalLength,
				AvgRisk:        avgRisk,
				RiskDensity:    riskDensity,
				LoFLengthRatio: lofLengthRatio,
				EstimatedCost:  estimatedCost,
				BRE:            bre,
				ROI:            roi,
				Priority:       priority,
			}

			projects = append(projects, project)
		}

		if len(projects) >= cfg.TargetCount*2 { // Generate more candidates than needed
			break
		}
	}

	return projects
}

// displayAdvancedProjects shows projects with comprehensive metrics
func displayAdvancedProjects(projects []ProjectMetrics, strategy string) {
	fmt.Printf("Top %d Projects for %s Strategy:\n", len(projects), strategy)
	fmt.Println("ID | Pipes | Length(m) | AvgRisk | RiskDens | Cost(€) | BRE(€) | ROI | Priority")
	fmt.Println("----|-------|-----------|---------|----------|---------|---------|-----|--------")

	for i, project := range projects {
		if i >= 10 { // Show top 10
			break
		}
		fmt.Printf("%2d | %5d | %7.0f | %7.3f | %8.1f | %7.0f | %7.0f | %.2f | %8.2f\n",
			project.ID,
			len(project.Nodes),
			project.TotalLength,
			project.AvgRisk,
			project.RiskDensity,
			project.EstimatedCost,
			project.BRE,
			project.ROI,
			project.Priority)
	}
	fmt.Println()
}

// compareStrategies compares the effectiveness of different strategies
func compareStrategies(allProjects map[string][]ProjectMetrics) {
	fmt.Printf("Strategy Comparison Summary:\n")
	fmt.Printf("%-20s | %s | %s | %s | %s\n", "Strategy", "Avg ROI", "Avg Risk", "Total Cost(€)", "Total BRE(€)")
	fmt.Printf("---------------------|---------|----------|-------------|-------------\n")

	for strategy, projects := range allProjects {
		if len(projects) == 0 {
			continue
		}

		var totalROI, totalRisk, totalCost, totalBRE float64
		for _, project := range projects {
			totalROI += project.ROI
			totalRisk += project.AvgRisk
			totalCost += project.EstimatedCost
			totalBRE += project.BRE
		}

		avgROI := totalROI / float64(len(projects))
		avgRisk := totalRisk / float64(len(projects))

		fmt.Printf("%-20s | %7.2f | %8.3f | %11.0f | %11.0f\n",
			strategy, avgROI, avgRisk, totalCost, totalBRE)
	}
}

// recommendBestStrategy provides a final recommendation
func recommendBestStrategy(allProjects map[string][]ProjectMetrics) {
	bestStrategy := ""
	bestScore := 0.0

	for strategy, projects := range allProjects {
		if len(projects) == 0 {
			continue
		}

		var totalROI, totalPriority float64
		for _, project := range projects {
			totalROI += project.ROI
			totalPriority += project.Priority
		}

		// Score based on average ROI and priority
		avgROI := totalROI / float64(len(projects))
		avgPriority := totalPriority / float64(len(projects))
		score := avgROI*0.6 + avgPriority*0.4

		if score > bestScore {
			bestScore = score
			bestStrategy = strategy
		}
	}

	if bestStrategy != "" {
		fmt.Printf("Recommended Strategy: %s\n", bestStrategy)
		fmt.Printf("This strategy offers the best combination of ROI and overall priority.\n")

		projects := allProjects[bestStrategy]
		var totalCost, totalBRE float64
		for _, project := range projects {
			totalCost += project.EstimatedCost
			totalBRE += project.BRE
		}

		fmt.Printf("\nExpected Outcomes:\n")
		fmt.Printf("  Total Investment: €%.0f\n", totalCost)
		fmt.Printf("  Total Risk Exposure Addressed: €%.0f\n", totalBRE)
		fmt.Printf("  Expected Payback Period: %.1f years\n", totalCost/totalBRE*10) // Rough estimate
		fmt.Printf("  Number of Projects: %d\n", len(projects))
	}
}
