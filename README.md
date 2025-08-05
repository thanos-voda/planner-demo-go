# Pipes Network Risk Assessment and Planning System

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

## Usage

### Running the Demos

```bash
# Basic demo
go run ./cmd/demo

# Advanced demo with ROI analysis
go run ./cmd/advanced_demo
```

### Running Tests

```bash
# All tests
go test ./...

# Neighbourhood tests only
go test ./internal/neighbourhood -v

# Integration tests
go test ./internal/neighbourhood -v -run TestNeighbourhoodIntegration
```

## Sample Output

### Network Analysis
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

## Future Enhancements

1. **Advanced Metrics**
   - Age-based risk modeling
   - Material-specific failure rates
   - Environmental impact factors

2. **Optimization Algorithms**
   - Genetic algorithms for project selection
   - Multi-objective optimization
   - Budget constraint handling

3. **Visualization**
   - Network graph visualization
   - Risk heat maps
   - Project impact analysis

4. **Real-world Integration**
   - GIS data import
   - Historical failure data
   - Maintenance scheduling

## Performance

- Handles 2000+ pipe networks efficiently
- O(n log n) neighbourhood finding complexity
- Memory-efficient graph representation
- Comprehensive test coverage with multiple scenarios

## Testing

The system includes extensive test coverage:

- **Unit Tests**: Core neighbourhood algorithm functionality
- **Integration Tests**: Cross-package functionality
- **Edge Cases**: Single nodes, disconnected networks, boundary conditions
- **Mock Networks**: Realistic test scenarios
- **Performance Tests**: Large network handling

All tests pass and provide confidence in the system's reliability and correctness.
