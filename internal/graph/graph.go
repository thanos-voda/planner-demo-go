package graph

type ID int64

type Node struct {
	ID     ID
	Score  float64
	Length float64
}

type Network struct {
	Nodes map[ID]*Node
	Edges map[ID][]ID
}

func New() *Network {
	return &Network{
		Nodes: make(map[ID]*Node),
		Edges: make(map[ID][]ID),
	}
}

func (n *Network) AddNode(node *Node) {
	n.Nodes[node.ID] = node
}

func (n *Network) AddUndirectedEdge(from, to ID) {
	n.Edges[from] = append(n.Edges[from], to)
	n.Edges[to] = append(n.Edges[to], from)
}
