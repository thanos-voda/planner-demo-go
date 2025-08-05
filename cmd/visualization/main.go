package main

import (
	"fmt"
	"sort"

	"github.com/thanos-fil/planner-demo-go/internal/graph"
	"github.com/thanos-fil/planner-demo-go/internal/neighbourhood"
)

// Step represents one step in the algorithm execution
type Step struct {
	StepNum      int
	Action       string
	CurrentNode  graph.ID
	CumLen       float64
	QueueState   []QueueItem
	VisitedNodes map[graph.ID]float64
	InQueue      map[graph.ID]float64
}

// QueueItem represents an item in the priority queue for visualization
type QueueItem struct {
	ID     graph.ID
	CumLen float64
}

func main() {
	fmt.Println("=== Neighbourhood Algorithm Step-by-Step Visualization ===")
	fmt.Println()

	// Create a small, traceable network for demonstration
	net := createSmallDemoNetwork()

	// Display the network structure
	fmt.Println("=== Network Structure ===")
	displayNetworkStructure(net)

	// Choose parameters for the demonstration
	startNode := graph.ID(0)
	maxLength := 8.0

	fmt.Printf("\n=== Algorithm Parameters ===\n")
	fmt.Printf("Start Node: %d\n", startNode)
	fmt.Printf("Max Length Budget: %.1f\n", maxLength)
	fmt.Printf("Start Node Length: %.1f\n", net.Nodes[startNode].Length)

	// Run the algorithm with step-by-step tracking
	fmt.Printf("\n=== Step-by-Step Algorithm Execution ===\n")
	result, steps := runNeighbourhoodWithSteps(net, startNode, maxLength)

	// Display each step
	displaySteps(steps)

	// Show final result
	fmt.Printf("\n=== Final Result ===\n")
	fmt.Printf("Reachable nodes within budget %.1f:\n", maxLength)
	for nodeID := range result {
		node := net.Nodes[graph.ID(nodeID)]
		fmt.Printf("  Node %d: LoF=%.2f, Length=%.1f\n", nodeID, node.Score, node.Length)
	}

	// Create ASCII graph visualization
	fmt.Printf("\n=== Network Graph Visualization ===\n")
	createASCIIGraph(net, result, startNode, maxLength)

	// Analyze the project metrics
	fmt.Printf("\n=== Project Analysis ===\n")
	analyzeProject(net, result)
}

// createSmallDemoNetwork creates a small network perfect for step-by-step demonstration
func createSmallDemoNetwork() *graph.Network {
	net := graph.New()

	// Create nodes with specific, easy-to-follow values
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

	// Create a specific topology that's easy to visualize
	// Network looks like:
	//     3
	//     |
	// 1---0---2---5
	//     |       |
	//     4       6
	//
	net.AddUndirectedEdge(0, 1) // 0 connects to 1
	net.AddUndirectedEdge(0, 2) // 0 connects to 2
	net.AddUndirectedEdge(0, 3) // 0 connects to 3
	net.AddUndirectedEdge(0, 4) // 0 connects to 4
	net.AddUndirectedEdge(2, 5) // 2 connects to 5
	net.AddUndirectedEdge(5, 6) // 5 connects to 6

	return net
}

// displayNetworkStructure shows the network in a table format
func displayNetworkStructure(net *graph.Network) {
	fmt.Println("Node | LoF  | Length | Neighbors")
	fmt.Println("-----|------|--------|----------")

	// Sort nodes by ID for consistent display
	var nodeIDs []int
	for id := range net.Nodes {
		nodeIDs = append(nodeIDs, int(id))
	}
	sort.Ints(nodeIDs)

	for _, id := range nodeIDs {
		node := net.Nodes[graph.ID(id)]
		neighbors := net.Edges[graph.ID(id)]

		fmt.Printf("  %d  | %.1f | %5.1f  | %v\n",
			id, node.Score, node.Length, neighbors)
	}
}

// runNeighbourhoodWithSteps runs the algorithm and captures each step
func runNeighbourhoodWithSteps(net *graph.Network, root graph.ID, maxLen float64) (map[neighbourhood.PipeID]struct{}, []Step) {
	var steps []Step
	stepNum := 0

	// Initialize tracking variables
	inQueue := map[neighbourhood.PipeID]float64{root: net.Nodes[root].Length}
	visited := map[neighbourhood.PipeID]struct{}{}
	visitedWithDistance := map[graph.ID]float64{}

	// We'll simulate the priority queue for demonstration
	type pqItem struct {
		id     neighbourhood.PipeID
		cumLen float64
	}

	var pq []pqItem
	pq = append(pq, pqItem{id: root, cumLen: inQueue[root]})

	// Initial state
	stepNum++
	steps = append(steps, Step{
		StepNum:      stepNum,
		Action:       fmt.Sprintf("Initialize: Add root node %d with cumLen=%.1f to queue", root, inQueue[root]),
		CurrentNode:  root,
		CumLen:       inQueue[root],
		QueueState:   []QueueItem{{ID: root, CumLen: inQueue[root]}},
		VisitedNodes: make(map[graph.ID]float64),
		InQueue:      copyPipeIDMap(inQueue),
	})

	for len(pq) > 0 {
		// Sort pq to simulate priority queue (smallest cumLen first)
		sort.Slice(pq, func(i, j int) bool {
			return pq[i].cumLen < pq[j].cumLen
		})

		// Pop the item with smallest cumulative length
		current := pq[0]
		pq = pq[1:]

		stepNum++
		var queueState []QueueItem
		for _, item := range pq {
			queueState = append(queueState, QueueItem{ID: item.id, CumLen: item.cumLen})
		}

		if current.cumLen > maxLen {
			steps = append(steps, Step{
				StepNum:      stepNum,
				Action:       fmt.Sprintf("Pop node %d (cumLen=%.1f) - EXCEEDS BUDGET, skip", current.id, current.cumLen),
				CurrentNode:  current.id,
				CumLen:       current.cumLen,
				QueueState:   queueState,
				VisitedNodes: copyGraphIDMap(visitedWithDistance),
				InQueue:      copyPipeIDMap(inQueue),
			})
			continue
		}

		// Visit the node
		visited[current.id] = struct{}{}
		visitedWithDistance[graph.ID(current.id)] = current.cumLen

		steps = append(steps, Step{
			StepNum:      stepNum,
			Action:       fmt.Sprintf("Pop and visit node %d (cumLen=%.1f)", current.id, current.cumLen),
			CurrentNode:  current.id,
			CumLen:       current.cumLen,
			QueueState:   queueState,
			VisitedNodes: copyGraphIDMap(visitedWithDistance),
			InQueue:      copyPipeIDMap(inQueue),
		})

		// Explore neighbors
		for _, nbr := range net.Edges[current.id] {
			nextLen := current.cumLen + net.Nodes[nbr].Length

			stepNum++
			if prev, ok := inQueue[neighbourhood.PipeID(nbr)]; !ok || nextLen < prev {
				inQueue[neighbourhood.PipeID(nbr)] = nextLen
				pq = append(pq, pqItem{id: neighbourhood.PipeID(nbr), cumLen: nextLen})

				// Update queue state for display
				queueState = nil
				for _, item := range pq {
					queueState = append(queueState, QueueItem{ID: item.id, CumLen: item.cumLen})
				}

				action := fmt.Sprintf("  → Explore neighbor %d: newCumLen=%.1f", nbr, nextLen)
				if ok {
					action += fmt.Sprintf(" (improved from %.1f)", prev)
				} else {
					action += " (first time)"
				}
				action += " → Add to queue"

				steps = append(steps, Step{
					StepNum:      stepNum,
					Action:       action,
					CurrentNode:  current.id,
					CumLen:       current.cumLen,
					QueueState:   queueState,
					VisitedNodes: copyGraphIDMap(visitedWithDistance),
					InQueue:      copyPipeIDMap(inQueue),
				})
			} else {
				steps = append(steps, Step{
					StepNum:      stepNum,
					Action:       fmt.Sprintf("  → Explore neighbor %d: newCumLen=%.1f ≥ existing %.1f → Skip", nbr, nextLen, prev),
					CurrentNode:  current.id,
					CumLen:       current.cumLen,
					QueueState:   queueState,
					VisitedNodes: copyGraphIDMap(visitedWithDistance),
					InQueue:      copyPipeIDMap(inQueue),
				})
			}
		}
	}

	return visited, steps
}

// Helper functions for copying maps
func copyPipeIDMap(m map[neighbourhood.PipeID]float64) map[neighbourhood.PipeID]float64 {
	result := make(map[neighbourhood.PipeID]float64)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func copyGraphIDMap(m map[graph.ID]float64) map[graph.ID]float64 {
	result := make(map[graph.ID]float64)
	for k, v := range m {
		result[k] = v
	}
	return result
}

// displaySteps shows each step of the algorithm execution
func displaySteps(steps []Step) {
	for _, step := range steps {
		fmt.Printf("Step %2d: %s\n", step.StepNum, step.Action)

		// Show queue state
		if len(step.QueueState) > 0 {
			fmt.Printf("         Queue: ")
			for i, item := range step.QueueState {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("[%d:%.1f]", item.ID, item.CumLen)
			}
			fmt.Println()
		} else {
			fmt.Printf("         Queue: [empty]\n")
		}

		// Show visited nodes
		if len(step.VisitedNodes) > 0 {
			fmt.Printf("         Visited: ")
			var visitedIDs []int
			for id := range step.VisitedNodes {
				visitedIDs = append(visitedIDs, int(id))
			}
			sort.Ints(visitedIDs)
			for i, id := range visitedIDs {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%d(%.1f)", id, step.VisitedNodes[graph.ID(id)])
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

// createASCIIGraph creates a visual representation of the network and result
func createASCIIGraph(net *graph.Network, result map[neighbourhood.PipeID]struct{}, startNode graph.ID, maxLen float64) {
	fmt.Println("Network Topology and Algorithm Result:")
	fmt.Println()
	fmt.Println("Legend:")
	fmt.Println("  [X] = Visited (within budget)")
	fmt.Println("  ( ) = Not visited (outside budget)")
	fmt.Println("  >>> = Start node")
	fmt.Println()

	// Create the ASCII representation
	fmt.Println("        [3] LoF:0.9, Len:1.0")
	if _, visited := result[neighbourhood.PipeID(3)]; visited {
		fmt.Println("         |  ✓")
	} else {
		fmt.Println("         |  ✗")
	}

	fmt.Print("    ")
	if _, visited := result[neighbourhood.PipeID(1)]; visited {
		fmt.Print("[1]")
	} else {
		fmt.Print("(1)")
	}
	fmt.Print("---")

	if startNode == 0 {
		fmt.Print(">>>[0]<<<")
	} else if _, visited := result[neighbourhood.PipeID(0)]; visited {
		fmt.Print("[0]")
	} else {
		fmt.Print("(0)")
	}

	fmt.Print("---")
	if _, visited := result[neighbourhood.PipeID(2)]; visited {
		fmt.Print("[2]")
	} else {
		fmt.Print("(2)")
	}
	fmt.Print("---")
	if _, visited := result[neighbourhood.PipeID(5)]; visited {
		fmt.Print("[5]")
	} else {
		fmt.Print("(5)")
	}
	fmt.Println()

	// Add length info
	fmt.Println("   1.5     2.0     3.0     2.5")

	fmt.Print("        ")
	if _, visited := result[neighbourhood.PipeID(4)]; visited {
		fmt.Print("[4]")
	} else {
		fmt.Print("(4)")
	}
	fmt.Print("               ")
	if _, visited := result[neighbourhood.PipeID(6)]; visited {
		fmt.Print("[6]")
	} else {
		fmt.Print("(6)")
	}
	fmt.Println()

	fmt.Println("        4.0               1.8")
	fmt.Println()

	// Add detailed node information
	fmt.Println("Node Details:")
	var nodeIDs []int
	for id := range net.Nodes {
		nodeIDs = append(nodeIDs, int(id))
	}
	sort.Ints(nodeIDs)

	for _, id := range nodeIDs {
		node := net.Nodes[graph.ID(id)]
		status := "✗ Not visited"
		if _, visited := result[neighbourhood.PipeID(id)]; visited {
			status = "✓ Visited"
		}
		if graph.ID(id) == startNode {
			status += " (START)"
		}
		fmt.Printf("  Node %d: LoF=%.1f, Length=%.1f → %s\n",
			id, node.Score, node.Length, status)
	}
}

// analyzeProject provides analysis of the selected project
func analyzeProject(net *graph.Network, result map[neighbourhood.PipeID]struct{}) {
	var totalRisk, totalLength float64
	var nodeCount int

	fmt.Printf("Project contains %d pipes:\n", len(result))

	var resultNodes []int
	for id := range result {
		resultNodes = append(resultNodes, int(id))
	}
	sort.Ints(resultNodes)

	for _, id := range resultNodes {
		node := net.Nodes[graph.ID(id)]
		totalRisk += node.Score
		totalLength += node.Length
		nodeCount++
		fmt.Printf("  Pipe %d: LoF=%.2f, Length=%.1fm\n", id, node.Score, node.Length)
	}

	fmt.Printf("\nProject Metrics:\n")
	fmt.Printf("  Total Pipes: %d\n", nodeCount)
	fmt.Printf("  Total Length: %.1f m\n", totalLength)
	fmt.Printf("  Total Risk (LoF): %.2f\n", totalRisk)
	fmt.Printf("  Average Risk: %.3f\n", totalRisk/float64(nodeCount))
	fmt.Printf("  Risk Density: %.2f LoF/km\n", totalRisk/(totalLength/1000))

	// Cost analysis
	estimatedCost := totalLength*500 + 10000
	avgCoF := 27500.0 // Average cost of failure
	bre := totalRisk * avgCoF
	roi := bre / estimatedCost

	fmt.Printf("\nFinancial Analysis:\n")
	fmt.Printf("  Estimated Cost: €%.0f (€500/m + €10k fixed)\n", estimatedCost)
	fmt.Printf("  Business Risk Exposure: €%.0f\n", bre)
	fmt.Printf("  ROI: %.2f\n", roi)
}
