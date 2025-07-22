func main() {
	net := planner.MockExampleNetwork()

	cfg := planner.UserCfg{}

	clusters := planner.CreateCandidateClusters(net, cfg)
}