# PathwayDB Redis Command Reference

This document provides a comprehensive reference for all custom Redis commands supported by **PathwayDB**.

All PathwayDB commands are **namespaced** to avoid conflicts with standard Redis commands.  
The available namespaces are: `GRAPH`, `NODE`, `EDGE`, and `ANALYSIS`.

---

## `GRAPH` Commands

Commands for managing graphs.

### `GRAPH.CREATE`

Creates a new graph.

- **Syntax**:
```redis
GRAPH.CREATE <name> [description]
```

- **Example Input**:
```redis
> GRAPH.CREATE my-graph "My first graph"
```

- **Example Output**:
```redis
OK
```

### `GRAPH.DELETE`

Deletes a graph and all of its associated nodes, edges, and indexes.

- **Syntax**:
```redis
GRAPH.DELETE <name>
```

- **Example Input**:
```redis
> GRAPH.DELETE my-graph
```

- **Example Output**:
```redis
OK
```

### `GRAPH.LIST`

Lists all graphs in the database.

- **Syntax**:
```redis
GRAPH.LIST
```

- **Example Input**:
```redis
> GRAPH.LIST
```

- **Example Output**:
```redis
1) "my-graph"
2) "My first graph"
3) "another-graph"
4) ""
```

### `GRAPH.GET`

Retrieves information about a specific graph.

- **Syntax**:
```redis
GRAPH.GET <name>
```

- **Example Input**:
```redis
> GRAPH.GET my-graph
```

- **Example Output**:
```redis
1) "my-graph"
2) "my-graph"
3) "My first graph"
4) "15"  # Node count
5) "30"  # Edge count
```

### `GRAPH.EXISTS`

Checks if a graph with the given name exists.

- **Syntax**:
```redis
GRAPH.EXISTS <name>
```

- **Example Input**:
```redis
> GRAPH.EXISTS my-graph
```

- **Example Output**:
```redis
(integer) 1
```

---

## `NODE` Commands

Commands for managing nodes within a graph.

### `NODE.CREATE`

Creates or fully replaces (upserts) a node in a graph.

- **Syntax**:
```redis
NODE.CREATE <graph> <id> <type> [attributes_json] [TTL <seconds>]
```

- **Example Input**:
```redis
> NODE.CREATE my-graph service-a service '{"version":"1.0", "region":"us-east-1"}' TTL 3600
```

- **Example Output**:
```redis
OK
```

### `NODE.GET`

Retrieves the details of a specific node.

- **Syntax**:
```redis
NODE.GET <graph> <id>
```

- **Example Input**:
```redis
> NODE.GET my-graph service-a
```

- **Example Output**:
```redis
1) "service-a"
2) "service"
3) "{"region":"us-east-1","version":"1.0"}"
4) "2025-09-16T08:46:12Z"
```

### `NODE.UPDATE`

Updates an existing node's type, attributes, and/or TTL. At least one update parameter must be provided.

- **Syntax**:
```redis
NODE.UPDATE <graph> <id> [TYPE <new_type>] [ATTRIBUTES <attributes_json>] [TTL <seconds>]
```

- **Parameters**:
  - `TYPE <new_type>`: (Optional) Updates the node's type
  - `ATTRIBUTES <attributes_json>`: (Optional) Updates the node's attributes with JSON
  - `TTL <seconds>`: (Optional) Sets expiration time in seconds (0 removes expiration)

- **Example Inputs**:
```redis
> NODE.UPDATE my-graph service-a TYPE microservice
> NODE.UPDATE my-graph service-a ATTRIBUTES '{"version":"1.1"}'
> NODE.UPDATE my-graph service-a TYPE microservice ATTRIBUTES '{"version":"2.0"}' TTL 3600
> NODE.UPDATE my-graph service-a TTL 0
```

- **Legacy Syntax** (still supported):
```redis
NODE.UPDATE <graph> <id> <attributes_json> [TTL <seconds>]
```

- **Example Output**:
```redis
OK
```

### `NODE.DELETE`

Deletes a node and all of its incoming and outgoing edges.

- **Syntax**:
```redis
NODE.DELETE <graph> <id>
```

- **Example Input**:
```redis
> NODE.DELETE my-graph service-a
```

- **Example Output**:
```redis
OK
```

### `NODE.FILTER`

Finds all nodes in a graph that have a specific attribute key-value pair.

- **Syntax**:
```redis
NODE.FILTER <graph> <attribute_key> <attribute_value>
```

- **Example Input**:
```redis
> NODE.FILTER my-graph region us-east-1
```

- **Example Output**:
```redis
1) "service-a"
2) "service"
3) "{"region":"us-east-1","version":"1.1"}"
4) "service-c"
5) "service"
6) "{"region":"us-east-1","version":"1.0"}"
```

### `NODE.LIST`

Lists all nodes in a specific graph.

- **Syntax**:
```redis
NODE.LIST <graph>
```

- **Example Input**:
```redis
> NODE.LIST my-graph
```

- **Example Output**:
```redis
1) "service-a:service"
2) "service-b:database"
```

### `NODE.EXISTS`

Checks if a node with the given ID exists in a graph.

- **Syntax**:
```redis
NODE.EXISTS <graph> <id>
```

- **Example Input**:
```redis
> NODE.EXISTS my-graph service-a
```

- **Example Output**:
```redis
(integer) 1
```

---

## `EDGE` Commands

Commands for managing edges within a graph.

### `EDGE.CREATE`

Creates or fully replaces (upserts) an edge between two nodes.

- **Syntax**:
```redis
EDGE.CREATE <graph> <id> <from> <to> <type> [attributes_json] [TTL <seconds>]
```

- **Example Input**:
```redis
> EDGE.CREATE my-graph edge-ab service-a service-b depends_on '{"protocol":"http"}'
```

- **Example Output**:
```redis
OK
```

### `EDGE.GET`

Retrieves the details of a specific edge.

- **Syntax**:
```redis
EDGE.GET <graph> <id>
```

- **Example Input**:
```redis
> EDGE.GET my-graph edge-ab
```

- **Example Output**:
```redis
1) "edge-ab"
2) "service-a"
3) "service-b"
4) "depends_on"
5) "{"protocol":"http"}"
6) ""
```

### `EDGE.UPDATE`

Updates the attributes of an existing edge.

- **Syntax**:
```redis
EDGE.UPDATE <graph> <id> <attributes_json> [TTL <seconds>]
```

- **Example Input**:
```redis
> EDGE.UPDATE my-graph edge-ab '{"protocol":"https"}'
```

- **Example Output**:
```redis
OK
```

### `EDGE.DELETE`

Deletes an edge.

- **Syntax**:
```redis
EDGE.DELETE <graph> <id>
```

- **Example Input**:
```redis
> EDGE.DELETE my-graph edge-ab
```

- **Example Output**:
```redis
OK
```

### `EDGE.FILTER`

Finds all edges in a graph that have a specific attribute key-value pair.

- **Syntax**:
```redis
EDGE.FILTER <graph> <attribute_key> <attribute_value>
```

- **Example Input**:
```redis
> EDGE.FILTER my-graph protocol https
```

- **Example Output**:
```redis
1) "edge-ab"
2) "service-a"
3) "service-b"
4) "depends_on"
5) "{"protocol":"https"}"
```

### `EDGE.NEIGHBORS`

Gets all neighboring nodes connected to a specified node.

- **Syntax**:
```redis
EDGE.NEIGHBORS <graph> <node> [DIRECTION in|out|both] [FORMAT simple|detailed]
```

- **Parameters**:
  - `DIRECTION`: Filter by edge direction relative to the specified node
    - `in`: Only incoming edges (neighbors that connect TO this node)
    - `out`: Only outgoing edges (neighbors that connect FROM this node)  
    - `both`: Both directions (default)
  - `FORMAT`: Output format
    - `simple`: Returns `neighbor_id:neighbor_type`
    - `detailed`: Returns `neighbor_id:neighbor_type<arrow>edge_id:edge_type` where `<arrow>` is `<-` for incoming edges or `->` for outgoing edges

- **Example Input (detailed)**:
```redis
> EDGE.NEIGHBORS my-graph service-a
```

- **Example Output (detailed)**:
```redis
1) "1"
2) "service-b:service->edge-ab:depends_on"
```

- **Example Input (simple)**:
```redis
> EDGE.NEIGHBORS my-graph service-a FORMAT simple
```

- **Example Output (simple)**:
```redis
1) "service-b:service"
```

### `EDGE.LIST`

Lists all edges in a specific graph.

- **Syntax**:
```redis
EDGE.LIST <graph>
```

- **Example Input**:
```redis
> EDGE.LIST my-graph
```

- **Example Output**:
```redis
1) "edge-ab:depends_on"
2) "edge-bc:depends_on"
```

### `EDGE.EXISTS`

Checks if an edge with the given ID exists in a graph.

- **Syntax**:
```redis
EDGE.EXISTS <graph> <id>
```

- **Example Input**:
```redis
> EDGE.EXISTS my-graph edge-ab
```

- **Example Output**:
```redis
(integer) 1
```

---

## `ANALYSIS` Commands

Commands for performing graph analysis.

### `ANALYSIS.SHORTESTPATH`

Finds the shortest path(s) between two nodes using BFS.

- **Syntax**:
```redis
ANALYSIS.SHORTESTPATH <graph> <from_node> <to_node> [FORMAT simple|detailed]
```

- **Example Input (detailed)**:
```redis
> ANALYSIS.SHORTESTPATH my-graph service-a service-c
```

- **Example Output (detailed)**:
```redis
1) "1"
2) "service-a:service->edge-ab:depends_on->service-b:service->edge-bc:depends_on->service-c:service"
```

- **Example Input (simple)**:
```redis
> ANALYSIS.SHORTESTPATH my-graph service-a service-c FORMAT simple
```

- **Example Output (simple)**:
```redis
1) "service-a:service"
2) "service-b:service"
3) "service-c:service"
```

### `ANALYSIS.CENTRALITY`

Calculates centrality measures for nodes in a graph.

- **Syntax**:
```redis
ANALYSIS.CENTRALITY <graph> <type> [node_id] [DIRECTION in|out|both]
```

- **Example Input**:
```redis
> ANALYSIS.CENTRALITY my-graph degree
```

- **Example Output**:
```redis
1) "service-a"
2) "1"
3) "service-b"
4) "2"
5) "service-c"
6) "1"
```

### `ANALYSIS.CLUSTERING`

Performs clustering analysis on a graph.

- **Syntax**:
```redis
ANALYSIS.CLUSTERING <graph> [algorithm] [parameters_json]
```

- **Example Input (Louvain)**:
```redis
> ANALYSIS.CLUSTERING my-graph louvain
```

- **Example Output (Louvain)**:
```redis
1) 1) "service-a"
   2) "service-b"
2) 1) "service-c"
```

- **Example Input (Connected Components)**:
```redis
> ANALYSIS.CLUSTERING my-graph connected_components
```

- **Example Output (Connected Components)**:
```redis
1) "connected_components"
2) "1"
```

### `ANALYSIS.CYCLES`

Finds all cycles in a graph, with optional filtering.

- **Syntax**:
```redis
ANALYSIS.CYCLES <graph> [NODETYPE type1...] [EDGETYPE type1...] [FORMAT simple|detailed]
```

- **Example Input**:
```redis
> ANALYSIS.CYCLES my-graph
```

- **Example Output**:
```redis
1) "1"
2) "service-a:service->edge-ab:depends_on->service-b:service->edge-ba:depends_on->service-a:service"
```

### `ANALYSIS.TRAVERSE`

Performs a traversal from a starting node.

- **Syntax**:
```redis
ANALYSIS.TRAVERSE <graph> <start_node> [DIRECTION <dir>] [NODETYPES type1...] [EDGETYPES type1...] [FORMAT simple|detailed]
```

- **Example Input**:
```redis
> ANALYSIS.TRAVERSE my-graph service-a
```

- **Example Output**:
```redis
1) "2"
2) "service-a:service->edge-ab:depends_on->service-b:service"
3) "service-a:service->edge-ac:depends_on->service-c:service"
```
