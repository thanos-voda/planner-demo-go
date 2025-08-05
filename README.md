# Pipes Network Risk Assessment and Planning System

## Quick Start

```bash
# Test the neighbourhood algorithm
go test ./internal/neighbourhood -v

# Run step-by-step visualization
go run cmd/visualization/main.go

# Run budget analysis demo
go run cmd/budget_analysis/main.go

# Run advanced network demo
go run cmd/advanced_demo/main.go
```

## Overview

This Go application provides a comprehensive risk assessment and project planning system for large-scale pipe networks. It implements sophisticated algorithms to identify high-risk areas and optimize infrastructure investment decisions based on various metrics including Likelihood of Failure (LoF), Return on Investment (ROI), and Business Risk Exposure (BRE).

## Architecture

### Core Components

1. **Graph Package** (`internal/graph/`)
   - `Network`: Represents the pipe network as a graph
   - `Node`: Represents individual pipes with ID, risk score (LoF), and length
   - `AddUndirectedEdge`: Creates bidirectional connections between pipes

2. **Neighbourhood Package** (`internal/neighbourhood/`)
   - Implements Dijkstra-like algorithm for finding connected pipes within a budget
   - Uses priority queue for efficient pathfinding
   - Supports both private (`neighbourhood`) and public (`Neighbourhood`) APIs

3. **Planner Package** (`internal/planner/`)
   - `Cluster`: Represents a project containing multiple pipes
   - `UserCfg`: Configuration for planning parameters
   - Project optimization and selection logic

## Key Features

### 1. Network Generation
- Creates realistic pipe networks with 1000-2000+ pipes
- Simulates realistic LoF values based on age, material, and environmental factors
- Grid-based topology with random cross-connections
- Pipe lengths between 5-50 meters (realistic urban segments)

### 2. Risk Analysis
- **Likelihood of Failure (LoF)**: Primary risk metric (0.0 - 1.0)
- **Risk Distribution**: Categorizes pipes into low/medium/high/critical risk
- **Network Overview**: Total length, average risk, pipe count statistics

### 3. Project Planning Strategies

#### Strategy 1: High LoF
- Focuses on pipes with highest likelihood of failure
- Prioritizes immediate risk reduction
- Good for emergency interventions

#### Strategy 2: High Risk Density
- Identifies areas with concentrated high-risk pipes
- Optimizes for geographic clustering
- Efficient for large-scale renovations

#### Strategy 3: High LoF/Length Ratio
- **Recommended Strategy** - Best ROI performance
- Balances risk reduction with implementation cost
- Prioritizes short, high-risk pipes

### 4. Financial Metrics

#### Cost of Failure (CoF)
- Minor Leak: €3,000 average
- Major Leak: €15,000 average
- Burst: €62,500 average
- Service Disruption: €30,000 average

#### Project Costs
- €500 per meter of pipe replacement
- €10,000 fixed costs per project

#### ROI Calculation
- **BRE (Business Risk Exposure)**: LoF × CoF
- **ROI**: BRE / Project Cost
- **Priority Score**: Weighted combination of ROI, risk density, and efficiency metrics

## Demo Applications

### Basic Demo (`cmd/demo/`)
- Simple risk assessment and project creation
- 1000 pipe network
- Basic clustering and risk reduction analysis

### Advanced Demo (`cmd/advanced_demo/`)
- Comprehensive ROI and financial analysis
- 2000 pipe network
- Multiple strategy comparison
- Detailed metrics and recommendations

### Visualization Demo (`cmd/visualization/`)
- Step-by-step algorithm execution
- Shows every decision the algorithm makes
- ASCII graph visualization
- Perfect for understanding algorithm mechanics

### Budget Analysis Demo (`cmd/budget_analysis/`)
- Tests different budget scenarios
- Shows impact of starting node selection
- Demonstrates ROI optimization
- Compares algorithm behavior under constraints

## Usage

### Running the Demos

```bash
# Basic demo (1000 pipes, simple analysis)
go run ./cmd/demo

# Advanced demo with ROI analysis (2000 pipes, financial metrics)
go run ./cmd/advanced_demo

# Step-by-step algorithm visualization (small network, detailed steps)
go run ./cmd/visualization

# Budget impact analysis (shows different budget scenarios)
go run ./cmd/budget_analysis
```

### Running Tests

```bash
# All tests
go test ./...

# Neighbourhood tests only (recommended for development)
go test ./internal/neighbourhood -v

# Specific test suites
go test ./internal/neighbourhood -v -run TestNeighbourhood
go test ./internal/neighbourhood -v -run TestNeighbourhoodIntegration
go test ./internal/neighbourhood -v -run TestNeighbourhoodEdgeCases

# Build verification (ensure everything compiles)
go build ./...
```

### Visualization and Analysis Workflow

For understanding the algorithm:

1. **Start with step-by-step visualization**:
   ```bash
   go run ./cmd/visualization
   ```
   This shows exactly how the algorithm makes decisions.

2. **Explore budget impacts**:
   ```bash
   go run ./cmd/budget_analysis
   ```
   See how different budgets affect project selection.

3. **Run realistic scenarios**:
   ```bash
   go run ./cmd/advanced_demo
   ```
   Experience full-scale network analysis with financial metrics.

4. **Verify with tests**:
   ```bash
   go test ./internal/neighbourhood -v
   ```
   Ensure algorithm correctness across various scenarios.

## Sample Output

### Step-by-Step Visualization Output
```
=== Neighbourhood Algorithm Step-by-Step Visualization ===

Network Structure:
Node | LoF  | Length | Neighbors
-----|------|--------|----------
  0  | 0.8 |   2.0  | [1 2 3 4]
  1  | 0.3 |   1.5  | [0]
  2  | 0.6 |   3.0  | [0 5]

Algorithm Execution:
Step  1: Initialize: Add root node 0 with cumLen=2.0 to queue
         Queue: [0:2.0]
Step  2: Pop and visit node 0 (cumLen=2.0)
         Queue: [empty]
         Visited: 0(2.0)
Step  3:   → Explore neighbor 1: newCumLen=3.5 (first time) → Add to queue

ASCII Graph Visualization:
        [3] LoF:0.9, Len:1.0
         |  ✓
    [1]--->>>[0]<<<---[2]---[5]
   1.5     2.0     3.0     2.5
        [4]               (6)
        4.0               1.8

Legend: [X] = Visited, ( ) = Not visited, >>> = Start node
```

### Budget Analysis Output
```
=== Budget Scenario Analysis ===
Budget: 3.0 → Nodes: [0,3] → ROI: 4.07 → High-risk: 2/3
Budget: 5.0 → Nodes: [0,1,2,3] → ROI: 5.20 → High-risk: 2/3  
Budget: 8.0 → Nodes: [0,1,2,3,4,5] → ROI: 5.66 → High-risk: 3/3
Budget: 12.0 → Nodes: [0,1,2,3,4,5,6] → ROI: 5.99 → High-risk: 3/3

Recommended: Budget 8.0 provides optimal ROI (5.66) with full high-risk coverage
```

### Network Analysis (Advanced Demo)
```
Network Overview:
  Total Pipes: 2000
  Total Length: 54.16 km
  Average Risk (LoF): 0.400
  Average Length: 27.1 m

Risk Distribution:
  Critical Risk (≥0.7): 0 pipes (0.0%)
  High Risk (0.5-0.7): 421 pipes (21.1%)
  Medium Risk (0.3-0.5): 1166 pipes (58.3%)
  Low Risk (<0.3): 413 pipes (20.6%)
```

### Project Recommendations
```
Top Projects for High LoF/Length Strategy:
ID | Pipes | Length(m) | AvgRisk | RiskDens | Cost(€) | BRE(€) | ROI | Priority
 1 |   445 |   10201 |   0.405 |     17.7 | 5110682 | 4984456 | 0.98 |     7.55
 2 |   422 |   10408 |   0.399 |     16.2 | 5213845 | 4645685 | 0.89 |     6.90
```

### ROI Analysis
```
Expected Outcomes:
  Total Investment: €19,832,861
  Total Risk Exposure Addressed: €17,223,724
  Expected Payback Period: 11.5 years
  Number of Projects: 5
```

### Test Output
```
=== RUN   TestNeighbourhood
=== RUN   TestNeighbourhood/root_only_-_exact_budget
=== RUN   TestNeighbourhood/root_and_one_neighbor
=== RUN   TestNeighbourhood/all_reachable_nodes
--- PASS: TestNeighbourhood (0.00s)
=== RUN   TestNeighbourhoodIntegration
--- PASS: TestNeighbourhoodIntegration (0.00s)
PASS
ok      github.com/thanos-fil/planner-demo-go/internal/neighbourhood    0.168s
```

## Key Algorithms

### Neighbourhood Finding (Dijkstra-like)
1. Start from seed pipe with its length as initial cost
2. Use priority queue to process pipes by cumulative distance
3. Add neighbors if they improve the shortest path
4. Stop when budget (maxLength) is exceeded
5. Return all reachable pipes within budget

### Project Optimization
1. **Seed Selection**: Identify high-priority starting points
2. **Clustering**: Use neighbourhood algorithm to group connected pipes
3. **Metric Calculation**: Compute ROI, risk density, and priority scores
4. **Ranking**: Sort projects by combined priority metrics
5. **Selection**: Choose top projects within budget constraints

## Troubleshooting

### Common Issues

**Import errors or module issues:**
```bash
# Ensure you're in the project root directory
cd /path/to/planner-demo-go

# Clean and rebuild
go clean -cache
go mod tidy
```

**Tests failing:**
```bash
# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./internal/neighbourhood -v
```

**Large output in demos:**
```bash
# Pipe output to less for easier viewing
go run cmd/visualization/main.go | less

# Save output to file
go run cmd/advanced_demo/main.go > network_analysis.txt
```

### Performance Notes

- The neighbourhood algorithm is O(V log V + E) where V is nodes and E is edges
- Memory usage scales linearly with network size
- For networks >5000 pipes, consider implementing result caching
