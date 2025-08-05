package neighbourhood

import (
	"testing"

	"github.com/thanos-fil/planner-demo-go/internal/graph"
)

func TestNeighbourhood(t *testing.T) {
	// Create a simple test network
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

	tests := []struct {
		name     string
		root     graph.ID
		maxLen   float64
		expected []graph.ID
	}{
		{
			name:     "root only - very small budget",
			root:     0,
			maxLen:   0.5,
			expected: []graph.ID{}, // budget too small even for root
		},
		{
			name:     "root only - exact budget",
			root:     0,
			maxLen:   1.0,
			expected: []graph.ID{0},
		},
		{
			name:     "root and one neighbor",
			root:     0,
			maxLen:   2.0,
			expected: []graph.ID{0, 1},
		},
		{
			name:     "root and two neighbors",
			root:     0,
			maxLen:   4.0,
			expected: []graph.ID{0, 1, 2},
		},
		{
			name:     "all reachable nodes",
			root:     0,
			maxLen:   10.0,
			expected: []graph.ID{0, 1, 2, 3, 4},
		},
		{
			name:     "from middle node",
			root:     2,
			maxLen:   3.0,
			expected: []graph.ID{2, 1}, // Node 3 would need cumLen = 2.0 + 1.5 = 3.5, which exceeds 3.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := neighbourhood(net, tt.root, tt.maxLen)

			// Check that we got the expected number of nodes
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d nodes, got %d", len(tt.expected), len(result))
			}

			// Check that all expected nodes are present
			for _, expectedID := range tt.expected {
				if _, found := result[PipeID(expectedID)]; !found {
					t.Errorf("expected node %d to be in result, but it wasn't", expectedID)
				}
			}

			// Check that no unexpected nodes are present
			for nodeID := range result {
				found := false
				for _, expectedID := range tt.expected {
					if nodeID == PipeID(expectedID) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("unexpected node %d in result", nodeID)
				}
			}
		})
	}
}

func TestNeighbourhoodWithMockNetwork(t *testing.T) {
	// Test with the mock network from the planner package
	// This requires importing the planner package
	// We'll create a simple version here instead

	net := graph.New()

	// Add some nodes from the mock (simplified)
	net.AddNode(&graph.Node{ID: 1, Score: 0.5, Length: 1.0})
	net.AddNode(&graph.Node{ID: 5, Score: 0.4, Length: 1.0})
	net.AddNode(&graph.Node{ID: 9, Score: 0.6, Length: 1.0})
	net.AddNode(&graph.Node{ID: 11, Score: 0.3, Length: 1.0})
	net.AddNode(&graph.Node{ID: 12, Score: 0.3, Length: 1.5})
	net.AddNode(&graph.Node{ID: 16, Score: 0.8, Length: 1.0})

	// Add edges from the mock
	net.AddUndirectedEdge(1, 5)
	net.AddUndirectedEdge(5, 9)
	net.AddUndirectedEdge(11, 12)
	net.AddUndirectedEdge(12, 16)

	t.Run("neighbourhood from node 1", func(t *testing.T) {
		result := neighbourhood(net, 1, 2.5)

		// Should include nodes 1, 5, and 9 (total length: 1 + 1 + 1 = 3, but cumulative distances matter)
		expectedNodes := []graph.ID{1, 5}

		if len(result) != len(expectedNodes) {
			t.Errorf("expected %d nodes, got %d", len(expectedNodes), len(result))
		}

		for _, expectedID := range expectedNodes {
			if _, found := result[PipeID(expectedID)]; !found {
				t.Errorf("expected node %d to be in result", expectedID)
			}
		}
	})

	t.Run("neighbourhood from node 12", func(t *testing.T) {
		result := neighbourhood(net, 12, 2.5)

		// From node 12 (length 1.5):
		// - Node 12: cumLen = 1.5
		// - Node 11: cumLen = 1.5 + 1.0 = 2.5 (reachable)
		// - Node 16: cumLen = 1.5 + 1.0 = 2.5 (reachable)
		expectedNodes := []graph.ID{12, 11, 16}

		if len(result) != len(expectedNodes) {
			// Debug: print what we actually got
			t.Logf("Got nodes: %v", result)
			t.Errorf("expected %d nodes, got %d", len(expectedNodes), len(result))
		}

		for _, expectedID := range expectedNodes {
			if _, found := result[PipeID(expectedID)]; !found {
				t.Errorf("expected node %d to be in result", expectedID)
			}
		}
	})
}

func TestNeighbourhoodEdgeCases(t *testing.T) {
	net := graph.New()

	t.Run("single node network", func(t *testing.T) {
		net.AddNode(&graph.Node{ID: 0, Score: 0.0, Length: 1.0})

		result := neighbourhood(net, 0, 1.0)

		if len(result) != 1 {
			t.Errorf("expected 1 node, got %d", len(result))
		}

		if _, found := result[0]; !found {
			t.Error("expected node 0 to be in result")
		}
	})

	t.Run("disconnected network", func(t *testing.T) {
		net := graph.New()
		net.AddNode(&graph.Node{ID: 0, Score: 0.0, Length: 1.0})
		net.AddNode(&graph.Node{ID: 1, Score: 0.0, Length: 1.0})
		// No edges - nodes are disconnected

		result := neighbourhood(net, 0, 10.0)

		if len(result) != 1 {
			t.Errorf("expected 1 node, got %d", len(result))
		}

		if _, found := result[0]; !found {
			t.Error("expected node 0 to be in result")
		}
	})
}
