package storage

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/utils"
)

// CreateEdge creates a new edge in the specified graph
func (e *BadgerEngine) CreateEdge(graphID models.GraphID, edge *models.Edge) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}

	return e.db.Update(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}
		return tx.CreateEdge(graphID, edge)
	})
}

// GetEdge retrieves an edge by ID from the specified graph
func (e *BadgerEngine) GetEdge(graphID models.GraphID, edgeID models.EdgeID) (*models.Edge, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	var edge *models.Edge
	err := e.db.View(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}
		var err error
		edge, err = tx.GetEdge(graphID, edgeID)
		return err
	})

	return edge, err
}

// UpdateEdge updates an existing edge
func (e *BadgerEngine) UpdateEdge(graphID models.GraphID, edge *models.Edge) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}

	return e.db.Update(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}
		return tx.UpdateEdge(graphID, edge)
	})
}

// DeleteEdge deletes an edge
func (e *BadgerEngine) DeleteEdge(graphID models.GraphID, edgeID models.EdgeID) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}

	return e.db.Update(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}
		return tx.DeleteEdge(graphID, edgeID)
	})
}

// ListEdges returns all edges in the specified graph
func (e *BadgerEngine) ListEdges(graphID models.GraphID) ([]*models.Edge, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	var edges []*models.Edge
	prefix := utils.CreateEdgeIteratorPrefix(graphID)

	err := e.iterateWithPrefix(prefix, func(key []byte, value []byte) error {
		edge := &models.Edge{}
		err := edge.FromJSON(value)
		if err != nil {
			return fmt.Errorf("failed to deserialize edge: %w", err)
		}
		edges = append(edges, edge)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %w", err)
	}

	return edges, nil
}

// ListEdgesByType returns all edges of a specific type in the specified graph
func (e *BadgerEngine) ListEdgesByType(graphID models.GraphID, edgeType models.EdgeType) ([]*models.Edge, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	var edges []*models.Edge
	prefix := utils.CreateTypeIteratorPrefix(graphID, "e", string(edgeType))

	err := e.iterateWithPrefix(prefix, func(key []byte, value []byte) error {
		// The value in the type index is the edge ID. We extract it from the value, not the key.
		edgeID := models.EdgeID(value)
		edge, err := e.GetEdge(graphID, edgeID)
		if err != nil {
			// It's possible the edge was deleted but the index remains, so we can skip.
			return nil
		}
		edges = append(edges, edge)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list edges by type: %w", err)
	}

	return edges, nil
}

// GetOutgoingEdges returns all edges going out from a specific node
func (e *BadgerEngine) GetOutgoingEdges(graphID models.GraphID, nodeID models.NodeID) ([]*models.Edge, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	var edges []*models.Edge
	prefix := []byte(fmt.Sprintf("%sout:%s:%s:", utils.NodeIndexPrefix, graphID, nodeID))

	err := e.iterateWithPrefix(prefix, func(key []byte, value []byte) error {
		// The value in the index is the edge ID.
		edgeID := models.EdgeID(value)
		edge, err := e.GetEdge(graphID, edgeID)
		if err != nil {
			// It's possible the edge was deleted but the index remains, so we can skip.
			return nil
		}
		edges = append(edges, edge)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get outgoing edges: %w", err)
	}

	return edges, nil
}

// GetIncomingEdges returns all edges coming into a specific node
func (e *BadgerEngine) GetIncomingEdges(graphID models.GraphID, nodeID models.NodeID) ([]*models.Edge, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	var edges []*models.Edge
	prefix := []byte(fmt.Sprintf("%sin:%s:%s:", utils.NodeIndexPrefix, graphID, nodeID))

	err := e.iterateWithPrefix(prefix, func(key []byte, value []byte) error {
		// The value in the index is the edge ID.
		edgeID := models.EdgeID(value)
		edge, err := e.GetEdge(graphID, edgeID)
		if err != nil {
			// It's possible the edge was deleted but the index remains, so we can skip.
			return nil
		}
		edges = append(edges, edge)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get incoming edges: %w", err)
	}

	return edges, nil
}

// GetConnectedNodes returns all nodes connected to a specific node (both incoming and outgoing)
func (e *BadgerEngine) GetConnectedNodes(graphID models.GraphID, nodeID models.NodeID) ([]*models.Node, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	nodeMap := make(map[models.NodeID]*models.Node)

	// Get outgoing edges and their target nodes
	outgoingEdges, err := e.GetOutgoingEdges(graphID, nodeID)
	if err != nil {
		return nil, err
	}

	for _, edge := range outgoingEdges {
		if _, exists := nodeMap[edge.ToNodeID]; !exists {
			node, err := e.GetNode(graphID, edge.ToNodeID)
			if err != nil {
				continue // Skip if node doesn't exist
			}
			nodeMap[edge.ToNodeID] = node
		}
	}

	// Get incoming edges and their source nodes
	incomingEdges, err := e.GetIncomingEdges(graphID, nodeID)
	if err != nil {
		return nil, err
	}

	for _, edge := range incomingEdges {
		if _, exists := nodeMap[edge.FromNodeID]; !exists {
			node, err := e.GetNode(graphID, edge.FromNodeID)
			if err != nil {
				continue // Skip if node doesn't exist
			}
			nodeMap[edge.FromNodeID] = node
		}
	}

	// Convert map to slice
	var nodes []*models.Node
	for _, node := range nodeMap {
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// FindEdgesByAttribute finds edges that have a specific attribute value
func (e *BadgerEngine) FindEdgesByAttribute(graphID models.GraphID, attrKey string, attrValue interface{}) ([]*models.Edge, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	// For now, we'll do a full scan of edges and filter by attribute
	// In a production system, you'd want to maintain attribute indexes
	allEdges, err := e.ListEdges(graphID)
	if err != nil {
		return nil, err
	}

	var matchingEdges []*models.Edge
	for _, edge := range allEdges {
		if value, exists := edge.GetAttribute(attrKey); exists {
			if value == attrValue {
				matchingEdges = append(matchingEdges, edge)
			}
		}
	}

	return matchingEdges, nil
}

// Transaction methods for edges

// CreateEdge creates an edge within a transaction
func (t *BadgerTransaction) CreateEdge(graphID models.GraphID, edge *models.Edge) error {
	// Verify that both nodes exist
	_, err := t.GetNode(graphID, edge.FromNodeID)
	if err != nil {
		return fmt.Errorf("source node does not exist: %w", err)
	}

	_, err = t.GetNode(graphID, edge.ToNodeID)
	if err != nil {
		return fmt.Errorf("target node does not exist: %w", err)
	}

	// Store the edge
	edgeKey := utils.EncodeEdgeKey(graphID, edge.ID)
	edgeValue, err := edge.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize edge: %w", err)
	}

	if edge.ExpiresAt != nil {
		ttl := time.Until(*edge.ExpiresAt)
		if ttl > 0 {
			err = t.setWithTTL(edgeKey, edgeValue, ttl)
		} else {
			// If TTL is already expired, don't even add it.
			return nil
		}
	} else {
		err = t.set(edgeKey, edgeValue)
	}
	if err != nil {
		return fmt.Errorf("failed to store edge: %w", err)
	}

	// Create type index
	typeIndexKey := utils.EncodeEdgeTypeIndexKey(graphID, edge.Type, edge.ID)
	err = t.set(typeIndexKey, []byte(edge.ID))
	if err != nil {
		return fmt.Errorf("failed to create type index: %w", err)
	}

	// Create outgoing edge index
	outIndexKey := utils.EncodeNodeOutEdgeIndexKey(graphID, edge.FromNodeID, edge.ID)
	err = t.set(outIndexKey, []byte(edge.ID))
	if err != nil {
		return fmt.Errorf("failed to create outgoing edge index: %w", err)
	}

	// Create incoming edge index
	inIndexKey := utils.EncodeNodeInEdgeIndexKey(graphID, edge.ToNodeID, edge.ID)
	err = t.set(inIndexKey, []byte(edge.ID))
	if err != nil {
		return fmt.Errorf("failed to create incoming edge index: %w", err)
	}

	return nil
}

// GetEdge retrieves an edge within a transaction
func (t *BadgerTransaction) GetEdge(graphID models.GraphID, edgeID models.EdgeID) (*models.Edge, error) {
	edgeKey := utils.EncodeEdgeKey(graphID, edgeID)
	edgeValue, err := t.get(edgeKey)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("edge not found: %s", edgeID)
		}
		return nil, fmt.Errorf("failed to get edge: %w", err)
	}

	edge := &models.Edge{}
	err = edge.FromJSON(edgeValue)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize edge: %w", err)
	}

	return edge, nil
}

// UpdateEdge updates an edge within a transaction
func (t *BadgerTransaction) UpdateEdge(graphID models.GraphID, edge *models.Edge) error {
	// Get the existing edge to check if connections or type changed
	existingEdge, err := t.GetEdge(graphID, edge.ID)
	if err != nil {
		return fmt.Errorf("edge does not exist: %w", err)
	}

	// If type changed, update the type index
	if existingEdge.Type != edge.Type {
		// Remove old type index
		oldTypeIndexKey := utils.EncodeEdgeTypeIndexKey(graphID, existingEdge.Type, edge.ID)
		err = t.delete(oldTypeIndexKey)
		if err != nil {
			return fmt.Errorf("failed to remove old type index: %w", err)
		}

		// Add new type index
		newTypeIndexKey := utils.EncodeEdgeTypeIndexKey(graphID, edge.Type, edge.ID)
		err = t.set(newTypeIndexKey, []byte(edge.ID))
		if err != nil {
			return fmt.Errorf("failed to create new type index: %w", err)
		}
	}

	// If connections changed, update the node indexes
	if existingEdge.FromNodeID != edge.FromNodeID || existingEdge.ToNodeID != edge.ToNodeID {
		// Remove old indexes
		oldOutIndexKey := utils.EncodeNodeOutEdgeIndexKey(graphID, existingEdge.FromNodeID, edge.ID)
		err = t.delete(oldOutIndexKey)
		if err != nil {
			return fmt.Errorf("failed to remove old outgoing edge index: %w", err)
		}

		oldInIndexKey := utils.EncodeNodeInEdgeIndexKey(graphID, existingEdge.ToNodeID, edge.ID)
		err = t.delete(oldInIndexKey)
		if err != nil {
			return fmt.Errorf("failed to remove old incoming edge index: %w", err)
		}

		// Verify new nodes exist
		_, err = t.GetNode(graphID, edge.FromNodeID)
		if err != nil {
			return fmt.Errorf("new source node does not exist: %w", err)
		}

		_, err = t.GetNode(graphID, edge.ToNodeID)
		if err != nil {
			return fmt.Errorf("new target node does not exist: %w", err)
		}

		// Add new indexes
		newOutIndexKey := utils.EncodeNodeOutEdgeIndexKey(graphID, edge.FromNodeID, edge.ID)
		err = t.set(newOutIndexKey, []byte(edge.ID))
		if err != nil {
			return fmt.Errorf("failed to create new outgoing edge index: %w", err)
		}

		newInIndexKey := utils.EncodeNodeInEdgeIndexKey(graphID, edge.ToNodeID, edge.ID)
		err = t.set(newInIndexKey, []byte(edge.ID))
		if err != nil {
			return fmt.Errorf("failed to create new incoming edge index: %w", err)
		}
	}

	// Update the edge
	edgeKey := utils.EncodeEdgeKey(graphID, edge.ID)
	edgeValue, err := edge.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize edge: %w", err)
	}

	if edge.ExpiresAt != nil {
		ttl := time.Until(*edge.ExpiresAt)
		if ttl > 0 {
			return t.setWithTTL(edgeKey, edgeValue, ttl)
		} else {
			// If TTL is expired, this update effectively becomes a delete.
			return t.DeleteEdge(graphID, edge.ID)
		}
	}

	return t.set(edgeKey, edgeValue)
}

// DeleteEdge deletes an edge within a transaction
func (t *BadgerTransaction) DeleteEdge(graphID models.GraphID, edgeID models.EdgeID) error {
	// Get the edge first to access its properties
	edge, err := t.GetEdge(graphID, edgeID)
	if err != nil {
		return fmt.Errorf("edge does not exist: %w", err)
	}

	// Delete the edge
	edgeKey := utils.EncodeEdgeKey(graphID, edgeID)
	err = t.delete(edgeKey)
	if err != nil {
		return fmt.Errorf("failed to delete edge: %w", err)
	}

	// Delete type index
	typeIndexKey := utils.EncodeEdgeTypeIndexKey(graphID, edge.Type, edgeID)
	err = t.delete(typeIndexKey)
	if err != nil {
		return fmt.Errorf("failed to delete type index: %w", err)
	}

	// Delete outgoing edge index
	outIndexKey := utils.EncodeNodeOutEdgeIndexKey(graphID, edge.FromNodeID, edgeID)
	err = t.delete(outIndexKey)
	if err != nil {
		return fmt.Errorf("failed to delete outgoing edge index: %w", err)
	}

	// Delete incoming edge index
	inIndexKey := utils.EncodeNodeInEdgeIndexKey(graphID, edge.ToNodeID, edgeID)
	err = t.delete(inIndexKey)
	if err != nil {
		return fmt.Errorf("failed to delete incoming edge index: %w", err)
	}

	return nil
}
