package storage

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/utils"
)

// CreateGraph creates a new graph
func (e *BadgerEngine) CreateGraph(graph *models.Graph) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}

	key := utils.EncodeGraphKey(graph.ID)
	value, err := graph.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize graph: %w", err)
	}

	return e.set(key, value)
}

// GetGraph retrieves a graph by ID
func (e *BadgerEngine) GetGraph(graphID models.GraphID) (*models.Graph, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	key := utils.EncodeGraphKey(graphID)
	value, err := e.get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("graph not found: %s", graphID)
		}
		return nil, fmt.Errorf("failed to get graph: %w", err)
	}

	graph := &models.Graph{}
	err = graph.FromJSON(value)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize graph: %w", err)
	}

	return graph, nil
}

// UpdateGraph updates an existing graph
func (e *BadgerEngine) UpdateGraph(graph *models.Graph) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}

	// Check if graph exists
	_, err := e.GetGraph(graph.ID)
	if err != nil {
		return fmt.Errorf("graph does not exist: %w", err)
	}

	key := utils.EncodeGraphKey(graph.ID)
	value, err := graph.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize graph: %w", err)
	}

	return e.set(key, value)
}

// DeleteGraph deletes a graph and all its nodes and edges
func (e *BadgerEngine) DeleteGraph(graphID models.GraphID) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}

	return e.db.Update(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}

		// 1. Delete all nodes, which will also trigger cascading deletion of connected edges.
		nodePrefix := utils.CreateNodeIteratorPrefix(graphID)
		nodeIter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer nodeIter.Close()

		for nodeIter.Seek(nodePrefix); nodeIter.ValidForPrefix(nodePrefix); nodeIter.Next() {
			item := nodeIter.Item()
			_, nodeID := utils.DecodeNodeKey(item.Key())
			if nodeID != "" {
				err := tx.DeleteNode(graphID, nodeID)
				if err != nil {
					// We continue even if a single node deletion fails, to attempt to clean up as much as possible.
					// Consider logging this error.
				}
			}
		}

		// 2. Delete any remaining orphaned edges that might not have been connected to the nodes iterated above.
		edgePrefix := utils.CreateEdgeIteratorPrefix(graphID)
		edgeIter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer edgeIter.Close()

		for edgeIter.Seek(edgePrefix); edgeIter.ValidForPrefix(edgePrefix); edgeIter.Next() {
			item := edgeIter.Item()
			_, edgeID := utils.DecodeEdgeKey(item.Key())
			if edgeID != "" {
				err := tx.DeleteEdge(graphID, edgeID)
				if err != nil {
					// Continue on error.
				}
			}
		}

		// 3. Delete the graph record itself.
		graphKey := utils.EncodeGraphKey(graphID)
		return txn.Delete(graphKey)
	})
}

// ListGraphs returns all graphs in the database
// CountNodes returns the total number of nodes in a graph
func (e *BadgerEngine) CountNodes(graphID models.GraphID) (int, error) {
	if e.db == nil {
		return 0, fmt.Errorf("database not opened")
	}

	count := 0
	prefix := utils.CreateNodeIteratorPrefix(graphID)

	err := e.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // We only need to count keys
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			count++
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to count nodes: %w", err)
	}

	return count, nil
}

// CountEdges returns the total number of edges in a graph
func (e *BadgerEngine) CountEdges(graphID models.GraphID) (int, error) {
	if e.db == nil {
		return 0, fmt.Errorf("database not opened")
	}

	count := 0
	prefix := utils.CreateEdgeIteratorPrefix(graphID)

	err := e.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // We only need to count keys
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			count++
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to count edges: %w", err)
	}

	return count, nil
}

func (e *BadgerEngine) ListGraphs() ([]*models.Graph, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	var graphs []*models.Graph
	prefix := []byte(utils.GraphPrefix)

	err := e.iterateWithPrefix(prefix, func(key []byte, value []byte) error {
		graph := &models.Graph{}
		err := graph.FromJSON(value)
		if err != nil {
			return fmt.Errorf("failed to deserialize graph: %w", err)
		}
		graphs = append(graphs, graph)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list graphs: %w", err)
	}

	return graphs, nil
}

// deleteWithPrefix deletes all keys with a given prefix within a transaction
func (e *BadgerEngine) deleteWithPrefix(txn *badger.Txn, prefix []byte) error {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false // We only need keys for deletion
	it := txn.NewIterator(opts)
	defer it.Close()

	var keysToDelete [][]byte
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		key := it.Item().KeyCopy(nil)
		keysToDelete = append(keysToDelete, key)
	}

	for _, key := range keysToDelete {
		err := txn.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
