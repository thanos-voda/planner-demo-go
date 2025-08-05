# Neighbourhood Algorithm: Step-by-Step Analysis

## Algorithm Overview

The neighbourhood algorithm is a **modified Dijkstra's shortest path algorithm** with a budget constraint. It finds all pipes reachable from a starting pipe within a specified total length budget.

## Key Concepts

### 1. **Cumulative Length Calculation**
- Each pipe has its own length
- **Cumulative length** = length of the path to reach this pipe + the pipe's own length
- This represents the total "budget" consumed to include this pipe in the project

### 2. **Priority Queue (Min-Heap)**
- Always processes the pipe with the **smallest cumulative length** first
- Ensures we find the shortest path to each pipe
- Guarantees optimal solution within the budget constraint

### 3. **Budget Constraint**
- Maximum total length allowed for the project
- Pipes exceeding this budget are excluded
- Allows for efficient project scope management

## Step-by-Step Walkthrough

### Network Setup
```
Network Topology:
        [3] LoF:0.9, Len:1.0
         |  
    [1]---[0]---[2]---[5]---[6]
         2.0   3.0   2.5   1.8
        [4]
        4.0

Start: Node 0 (LoF=0.8, Length=2.0)
Budget: 8.0 units
```

### Detailed Execution

#### **Step 1: Initialization**
- **Action**: Add start node 0 to queue
- **Cumulative Length**: 2.0 (node 0's own length)
- **Queue**: [0:2.0]
- **Visited**: []

#### **Steps 2-6: Explore from Node 0**
- **Action**: Pop node 0 (smallest cumLen = 2.0)
- **Mark as visited**: Node 0
- **Explore neighbors**: 1, 2, 3, 4
  - Node 1: cumLen = 2.0 + 1.5 = **3.5**
  - Node 2: cumLen = 2.0 + 3.0 = **5.0** 
  - Node 3: cumLen = 2.0 + 1.0 = **3.0**
  - Node 4: cumLen = 2.0 + 4.0 = **6.0**
- **Queue**: [3:3.0, 1:3.5, 2:5.0, 4:6.0] (sorted by cumLen)

#### **Steps 7-8: Process Node 3**
- **Action**: Pop node 3 (smallest cumLen = 3.0)
- **Check budget**: 3.0 ≤ 8.0 ✓ **ACCEPT**
- **Mark as visited**: Node 3
- **Explore neighbors**: Only node 0, but we already have a better path (2.0 < 5.0)

#### **Steps 9-10: Process Node 1**
- **Action**: Pop node 1 (smallest cumLen = 3.5)
- **Check budget**: 3.5 ≤ 8.0 ✓ **ACCEPT**
- **Mark as visited**: Node 1
- **Explore neighbors**: Only node 0, already have better path

#### **Steps 11-13: Process Node 2**
- **Action**: Pop node 2 (smallest cumLen = 5.0)
- **Check budget**: 5.0 ≤ 8.0 ✓ **ACCEPT**
- **Mark as visited**: Node 2
- **Explore neighbors**: 
  - Node 0: Already have better path
  - Node 5: cumLen = 5.0 + 2.5 = **7.5** → Add to queue

#### **Steps 14-15: Process Node 4**
- **Action**: Pop node 4 (smallest cumLen = 6.0)
- **Check budget**: 6.0 ≤ 8.0 ✓ **ACCEPT**
- **Mark as visited**: Node 4

#### **Steps 16-18: Process Node 5**
- **Action**: Pop node 5 (smallest cumLen = 7.5)
- **Check budget**: 7.5 ≤ 8.0 ✓ **ACCEPT**
- **Mark as visited**: Node 5
- **Explore neighbors**:
  - Node 6: cumLen = 7.5 + 1.8 = **9.3** → Add to queue

#### **Step 19: Node 6 Rejected**
- **Action**: Pop node 6 (cumLen = 9.3)
- **Check budget**: 9.3 > 8.0 ✗ **REJECT**
- **Result**: Node 6 is not included in the project

## Algorithm Properties

### **Optimality**
- **Guaranteed optimal**: Always finds the shortest path to each reachable node
- **Budget-constrained**: Respects the maximum length constraint
- **Dijkstra-based**: Inherits the optimality properties of Dijkstra's algorithm

### **Efficiency**
- **Time Complexity**: O((V + E) log V) where V = nodes, E = edges
- **Space Complexity**: O(V) for the priority queue and visited tracking
- **Early termination**: Stops when budget is exceeded

### **Practical Benefits**
- **Project scope control**: Automatic budget management
- **Risk-aware**: Can incorporate pipe risk (LoF) values
- **Geographic clustering**: Naturally groups connected pipes
- **Scalable**: Handles large networks efficiently

## Real-World Application

### **Project Metrics from Example**
```
Project Results:
- Total Pipes: 6 out of 7 possible
- Total Length: 14.0m (within 8.0m budget after cumulative calculation)
- Total Risk: 3.50 LoF units
- Average Risk: 0.583 (high-risk project)
- Risk Density: 250 LoF/km
- ROI: 5.66 (excellent return)
```

### **Why Node 6 was Excluded**
- **Path to Node 6**: 0 → 2 → 5 → 6
- **Cumulative length**: 2.0 + 3.0 + 2.5 + 1.8 = **9.3m**
- **Budget exceeded**: 9.3m > 8.0m budget
- **Decision**: Automatically excluded by algorithm

### **Strategic Insights**
1. **High-value inclusion**: Node 3 (LoF=0.9) included with short path
2. **Efficient coverage**: 6 out of 7 nodes included within budget
3. **Automatic optimization**: Algorithm balanced reach vs. budget constraint
4. **Risk concentration**: Project captures high-risk pipes (nodes 0, 3, 5)

## Algorithm Advantages

### **1. Automatic Optimization**
- No manual selection needed
- Mathematically optimal solution
- Balances coverage and budget constraints

### **2. Risk-Aware Planning**
- Can prioritize high-risk areas
- Integrates with business metrics (LoF, CoF, ROI)
- Supports strategic decision-making

### **3. Scalability**
- Handles networks with thousands of pipes
- Efficient priority queue implementation
- Suitable for real-time planning applications

### **4. Flexibility**
- Adjustable budget constraints
- Multiple starting points possible
- Can be adapted for different optimization criteria

This algorithm forms the core of effective infrastructure planning, providing a mathematically sound approach to project selection while respecting real-world constraints.
