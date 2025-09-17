#!/bin/bash

# Sample Graph Structure
# Creates a complex graph with multiple nodes and interconnected edges

echo "🏗️  Creating sample graph structure..."

# Create the main graph
redis-cli -p 6379 GRAPH.CREATE myGraph "Sample graph with complex node relationships"

echo "Creating nodes..."
redis-cli -p 6379 NODE.CREATE myGraph node-a node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-b node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-c node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-d node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-e node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-f node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-g node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-h node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-i node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-j node '{}'
redis-cli -p 6379 NODE.CREATE myGraph node-k node '{}'

echo "Creating edges..."
# From node-j (central hub)
redis-cli -p 6379 EDGE.CREATE myGraph edge9 node-j node-f line '{}'
redis-cli -p 6379 EDGE.CREATE myGraph edge10 node-j node-d line '{}'
redis-cli -p 6379 EDGE.CREATE myGraph edge7 node-j node-e line '{}'
redis-cli -p 6379 EDGE.CREATE myGraph edge8 node-j node-k line '{}'
redis-cli -p 6379 EDGE.CREATE myGraph edgeCycle node-j node-g line '{}'

# From node-b
redis-cli -p 6379 EDGE.CREATE myGraph edge1 node-b node-a line '{}'
redis-cli -p 6379 EDGE.CREATE myGraph edge3 node-b node-h line '{}'

# From node-e
redis-cli -p 6379 EDGE.CREATE myGraph edge4 node-e node-i line '{}'

# From node-g
redis-cli -p 6379 EDGE.CREATE myGraph edge6 node-g node-i line '{}'

echo ""
echo "✅ Sample graph created successfully!"
echo "📊 Graph contains:"
echo "   - 11 nodes (node-a through node-k)"
echo "   - 9 edges connecting various nodes"
echo "   - Central hub at node-j with multiple connections"
echo "   - Isolated node-c"
echo ""
echo "🚀 Try these commands:"
echo "   ANALYSIS.TRAVERSE myGraph node-j DIRECTION both"
echo "   EDGE.NEIGHBORS myGraph node-j out"
echo "   NODE.LIST myGraph"