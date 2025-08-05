package neighbourhood

import (
	"fmt"
	"testing"

	"github.com/thanos-fil/planner-demo-go/internal/graph"
)

func TestNeighbourhoodDebug(t *testing.T) {
	// Create a simple test network for debugging
	net := graph.New()

	// Add nodes with different lengths
	nodes := []*graph.Node{
		{ID: 0, Score: 0.0, Length: 1.0},
		{ID: 1, Score: 0.5, Length: 1.0},
		{ID: 2, Score: 0.0, Length: 2.0},
		{ID: 3, Score: 0.0, Length: 1.5},
		{ID: 4, Score: 0.0, Length: 1.0},
	}

	for _, node := range nodes {
		net.AddNode(node)
	}

	// Create a simple chain: 0 - 1 - 2 - 3 - 4
	net.AddUndirectedEdge(0, 1)
	net.AddUndirectedEdge(1, 2)
	net.AddUndirectedEdge(2, 3)
	net.AddUndirectedEdge(3, 4)

	t.Run("debug from node 2 with maxLen 3.0", func(t *testing.T) {
		result := neighbourhood(net, 2, 3.0)

		fmt.Printf("Network structure:\n")
		for id, node := range net.Nodes {
			fmt.Printf("Node %d: Length=%.1f, Neighbors=%v\n", id, node.Length, net.Edges[id])
		}

		fmt.Printf("\nResult from node 2 with maxLen=3.0:\n")
		for nodeID := range result {
			fmt.Printf("Node %d is reachable\n", nodeID)
		}

		// Let's manually calculate the expected distances:
		// From node 2:
		// - Node 2: cumLen = 2.0 (its own length)
		// - Node 1: cumLen = 2.0 + 1.0 = 3.0 (node 2's length + node 1's length)
		// - Node 3: cumLen = 2.0 + 1.5 = 3.5 (node 2's length + node 3's length)
		// - Node 0: cumLen = 3.0 + 1.0 = 4.0 (via node 1)
		// - Node 4: cumLen = 3.5 + 1.0 = 4.5 (via node 3)

		// So with maxLen=3.0, we should get nodes 2 and 1
		expectedCount := 2
		if len(result) != expectedCount {
			t.Errorf("expected %d nodes, got %d", expectedCount, len(result))
		}
	})

	t.Run("debug mock network node 12", func(t *testing.T) {
		net2 := graph.New()

		net2.AddNode(&graph.Node{ID: 11, Score: 0.3, Length: 1.0})
		net2.AddNode(&graph.Node{ID: 12, Score: 0.3, Length: 1.5})
		net2.AddNode(&graph.Node{ID: 16, Score: 0.8, Length: 1.0})

		net2.AddUndirectedEdge(11, 12)
		net2.AddUndirectedEdge(12, 16)

		result := neighbourhood(net2, 12, 2.0)

		fmt.Printf("\nMock network structure:\n")
		for id, node := range net2.Nodes {
			fmt.Printf("Node %d: Length=%.1f, Neighbors=%v\n", id, node.Length, net2.Edges[id])
		}

		fmt.Printf("\nResult from node 12 with maxLen=2.0:\n")
		for nodeID := range result {
			fmt.Printf("Node %d is reachable\n", nodeID)
		}

		// From node 12:
		// - Node 12: cumLen = 1.5 (its own length)
		// - Node 11: cumLen = 1.5 + 1.0 = 2.5 (too far)
		// - Node 16: cumLen = 1.5 + 1.0 = 2.5 (too far)

		// So with maxLen=2.0, we should only get node 12
		expectedCount := 1
		if len(result) != expectedCount {
			t.Errorf("expected %d nodes, got %d", expectedCount, len(result))
		}
	})
}
