package planner

import "github.com/thanos-fil/planner-demo-go/internal/graph"

func MockExampleNetwork() *graph.Network {
	net := graph.New()

	// Scores from the picture
	scores := map[int]float64{
		0: 0.0, 1: 0.5, 2: 0.0, 3: 0.0,
		4: 0.0, 5: 0.4, 6: 0.0, 7: 0.0,
		8: 0.0, 9: 0.6, 10: 0.0, 11: 0.3,
		12: 0.3, 13: 0.0, 14: 0.0, 15: 0.0,
		16: 0.8, 17: 0.0, 18: 0.0, 19: 0.0,
		20: 0.0, 21: 0.0, 22: 0.0, 23: 0.0, 24: 0.0,
	}
	length := map[int]float64{
		3: 1.5, 12: 1.5, // long segments
	}
	for i := 0; i <= 24; i++ {
		l := 1.0
		if v, ok := length[i]; ok {
			l = v
		}
		net.AddNode(&graph.Node{
			ID:     graph.ID(i),
			Score:  scores[i],
			Length: l,
		})
	}

	// Edges – only a sparse subset needed for unit tests.  Fill as required.
	net.AddUndirectedEdge(1, 5)
	net.AddUndirectedEdge(5, 9)
	net.AddUndirectedEdge(11, 12)
	net.AddUndirectedEdge(12, 16)
	return net
}
