package neighbourhood

import (
	"testing"

	"github.com/thanos-fil/planner-demo-go/internal/planner"
)

func TestNeighbourhoodIntegration(t *testing.T) {
	// Test the neighbourhood function with the actual mock network
	net := planner.MockExampleNetwork()

	t.Run("neighbourhood from node 1 in mock network", func(t *testing.T) {
		result := neighbourhood(net, 1, 3.0)

		// Should include nodes reachable within distance 3.0 from node 1
		// Based on the mock network: 1 - 5 - 9
		// Node 1: cumLen = 1.0
		// Node 5: cumLen = 1.0 + 1.0 = 2.0
		// Node 9: cumLen = 2.0 + 1.0 = 3.0

		expectedMinNodes := 3 // at least nodes 1, 5, and 9
		if len(result) < expectedMinNodes {
			t.Errorf("expected at least %d nodes, got %d", expectedMinNodes, len(result))
		}

		// Check that key nodes are present
		expectedNodes := []PipeID{1, 5, 9}
		for _, nodeID := range expectedNodes {
			if _, found := result[nodeID]; !found {
				t.Errorf("expected node %d to be in result", nodeID)
			}
		}
	})

	t.Run("neighbourhood from node 16 in mock network", func(t *testing.T) {
		result := neighbourhood(net, 16, 4.0)

		// Should include nodes reachable within distance 4.0 from node 16
		// Based on the mock network: 16 - 12 - 11
		// Node 16: cumLen = 1.0
		// Node 12: cumLen = 1.0 + 1.5 = 2.5
		// Node 11: cumLen = 2.5 + 1.0 = 3.5

		expectedMinNodes := 3 // at least nodes 16, 12, and 11
		if len(result) < expectedMinNodes {
			t.Errorf("expected at least %d nodes, got %d", expectedMinNodes, len(result))
		}

		// Check that key nodes are present
		expectedNodes := []PipeID{16, 12, 11}
		for _, nodeID := range expectedNodes {
			if _, found := result[nodeID]; !found {
				t.Errorf("expected node %d to be in result", nodeID)
			}
		}
	})
}
