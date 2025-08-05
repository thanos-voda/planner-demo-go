package neighbourhood

import (
	"container/heap"

	"github.com/thanos-fil/planner-demo-go/internal/graph"
)

// PipeID is an alias for graph.ID
type PipeID = graph.ID

// PriorityQueue implements heap.Interface for pqItem
type PriorityQueue []*pqItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cumLen < pq[j].cumLen
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

type pqItem struct {
	id     PipeID
	cumLen float64
	index  int
}

func neighbourhood(net *graph.Network, root graph.ID, maxLen float64) map[PipeID]struct{} {
	inQueue := map[PipeID]float64{root: net.Nodes[root].Length}
	visited := map[PipeID]struct{}{}

	pq := make(PriorityQueue, 0, 16)
	heap.Push(&pq, &pqItem{id: root, cumLen: inQueue[root]})

	for pq.Len() > 0 {
		it := heap.Pop(&pq).(*pqItem)
		if it.cumLen > maxLen {
			continue // pruned by budget
		}
		visited[it.id] = struct{}{}

		for _, nbr := range net.Edges[it.id] {
			nextLen := it.cumLen + net.Nodes[nbr].Length
			if prev, ok := inQueue[nbr]; !ok || nextLen < prev {
				inQueue[nbr] = nextLen
				heap.Push(&pq, &pqItem{id: nbr, cumLen: nextLen})
			}
		}
	}
	return visited
}

// Neighbourhood is the exported version of the neighbourhood function
func Neighbourhood(net *graph.Network, root graph.ID, maxLen float64) map[PipeID]struct{} {
	return neighbourhood(net, root, maxLen)
}
