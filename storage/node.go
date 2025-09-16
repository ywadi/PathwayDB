package storage

import (
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/utils"
)

// CreateNode creates a new node in the specified graph
func (e *BadgerEngine) CreateNode(graphID models.GraphID, node *models.Node) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}

	return e.db.Update(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}
		return tx.CreateNode(graphID, node)
	})
}

// GetNode retrieves a node by ID from the specified graph
func (e *BadgerEngine) GetNode(graphID models.GraphID, nodeID models.NodeID) (*models.Node, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	var node *models.Node
	err := e.db.View(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}
		var err error
		node, err = tx.GetNode(graphID, nodeID)
		return err
	})

	return node, err
}

// UpdateNode updates an existing node
func (e *BadgerEngine) UpdateNode(graphID models.GraphID, node *models.Node) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}

	return e.db.Update(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}
		return tx.UpdateNode(graphID, node)
	})
}

// DeleteNode deletes a node and all its associated edges
func (e *BadgerEngine) DeleteNode(graphID models.GraphID, nodeID models.NodeID) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}

	return e.db.Update(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}
		return tx.DeleteNode(graphID, nodeID)
	})
}

// ListNodes returns all nodes in the specified graph
func (e *BadgerEngine) ListNodes(graphID models.GraphID) ([]*models.Node, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	var nodes []*models.Node
	prefix := utils.CreateNodeIteratorPrefix(graphID)

	err := e.iterateWithPrefix(prefix, func(key []byte, value []byte) error {
		node := &models.Node{}
		err := node.FromJSON(value)
		if err != nil {
			return fmt.Errorf("failed to deserialize node: %w", err)
		}
		nodes = append(nodes, node)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	return nodes, nil
}

// ListNodesByType returns all nodes of a specific type in the specified graph
func (e *BadgerEngine) ListNodesByType(graphID models.GraphID, nodeType models.NodeType) ([]*models.Node, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	var nodes []*models.Node
	prefix := utils.CreateTypeIteratorPrefix(graphID, "n", string(nodeType))

	err := e.iterateWithPrefix(prefix, func(key []byte, value []byte) error {
		// The value in type index is just the node ID, we need to fetch the actual node
		keyStr := string(key)
		parts := strings.Split(keyStr, ":")
		if len(parts) >= 4 {
			nodeID := models.NodeID(parts[len(parts)-1])
			node, err := e.GetNode(graphID, nodeID)
			if err != nil {
				return err
			}
			nodes = append(nodes, node)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list nodes by type: %w", err)
	}

	return nodes, nil
}

// FindNodesByAttribute finds nodes that have a specific attribute value
func (e *BadgerEngine) FindNodesByAttribute(graphID models.GraphID, attrKey string, attrValue interface{}) ([]*models.Node, error) {
	if e.db == nil {
		return nil, fmt.Errorf("database not opened")
	}

	// For now, we'll do a full scan of nodes and filter by attribute
	// In a production system, you'd want to maintain attribute indexes
	allNodes, err := e.ListNodes(graphID)
	if err != nil {
		return nil, err
	}

	var matchingNodes []*models.Node
	for _, node := range allNodes {
		if value, exists := node.GetAttribute(attrKey); exists {
			if value == attrValue {
				matchingNodes = append(matchingNodes, node)
			}
		}
	}

	return matchingNodes, nil
}

// Transaction methods for nodes

// CreateNode creates a node within a transaction
func (t *BadgerTransaction) CreateNode(graphID models.GraphID, node *models.Node) error {
	// Store the node
	nodeKey := utils.EncodeNodeKey(graphID, node.ID)
	nodeValue, err := node.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize node: %w", err)
	}

	err = t.set(nodeKey, nodeValue)
	if err != nil {
		return fmt.Errorf("failed to store node: %w", err)
	}

	// Create type index
	typeIndexKey := utils.EncodeNodeTypeIndexKey(graphID, node.Type, node.ID)
	err = t.set(typeIndexKey, []byte(node.ID))
	if err != nil {
		return fmt.Errorf("failed to create type index: %w", err)
	}

	// Add to expiry index if TTL is set
	if node.ExpiresAt != nil {
		key := utils.EncodeExpiryIndexKey(graphID, node.ID, *node.ExpiresAt)
		err = t.set(key, []byte(node.ID))
		if err != nil {
			return fmt.Errorf("failed to create expiry index: %w", err)
		}
	}

	return nil
}

// GetNode retrieves a node within a transaction
func (t *BadgerTransaction) GetNode(graphID models.GraphID, nodeID models.NodeID) (*models.Node, error) {
	nodeKey := utils.EncodeNodeKey(graphID, nodeID)
	nodeValue, err := t.get(nodeKey)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, fmt.Errorf("node not found: %s", nodeID)
		}
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	node := &models.Node{}
	err = node.FromJSON(nodeValue)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize node: %w", err)
	}

	return node, nil
}

// UpdateNode updates a node within a transaction
func (t *BadgerTransaction) UpdateNode(graphID models.GraphID, node *models.Node) error {
	// Get the existing node to check if type changed
	existingNode, err := t.GetNode(graphID, node.ID)
	if err != nil {
		return fmt.Errorf("node does not exist: %w", err)
	}

	// If type changed, update the type index
	if existingNode.Type != node.Type {
		// Remove old type index
		oldTypeIndexKey := utils.EncodeNodeTypeIndexKey(graphID, existingNode.Type, node.ID)
		err = t.delete(oldTypeIndexKey)
		if err != nil {
			return fmt.Errorf("failed to remove old type index: %w", err)
		}

		// Add new type index
		newTypeIndexKey := utils.EncodeNodeTypeIndexKey(graphID, node.Type, node.ID)
		err = t.set(newTypeIndexKey, []byte(node.ID))
		if err != nil {
			return fmt.Errorf("failed to create new type index: %w", err)
		}
	}

	// Handle expiry index update
	if (existingNode.ExpiresAt != nil && node.ExpiresAt == nil) || (existingNode.ExpiresAt != nil && node.ExpiresAt != nil && !existingNode.ExpiresAt.Equal(*node.ExpiresAt)) {
		// TTL was removed or changed, so remove old index
		oldExpiryKey := utils.EncodeExpiryIndexKey(graphID, existingNode.ID, *existingNode.ExpiresAt)
		if err := t.delete(oldExpiryKey); err != nil {
			return fmt.Errorf("failed to remove old expiry index: %w", err)
		}
	}
	if (node.ExpiresAt != nil && existingNode.ExpiresAt == nil) || (node.ExpiresAt != nil && existingNode.ExpiresAt != nil && !node.ExpiresAt.Equal(*existingNode.ExpiresAt)) {
		// TTL was added or changed, so add new index
		newExpiryKey := utils.EncodeExpiryIndexKey(graphID, node.ID, *node.ExpiresAt)
		if err := t.set(newExpiryKey, []byte(node.ID)); err != nil {
			return fmt.Errorf("failed to create new expiry index: %w", err)
		}
	}

	// Update the node
	nodeKey := utils.EncodeNodeKey(graphID, node.ID)
	nodeValue, err := node.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize node: %w", err)
	}

	return t.set(nodeKey, nodeValue)
}

// DeleteNode deletes a node within a transaction
func (t *BadgerTransaction) DeleteNode(graphID models.GraphID, nodeID models.NodeID) error {
	// Get the node first to access its type
	node, err := t.GetNode(graphID, nodeID)
	if err != nil {
		return fmt.Errorf("node does not exist: %w", err)
	}

	// Delete the node
	nodeKey := utils.EncodeNodeKey(graphID, nodeID)
	err = t.delete(nodeKey)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	// Delete type index
	typeIndexKey := utils.EncodeNodeTypeIndexKey(graphID, node.Type, nodeID)
	err = t.delete(typeIndexKey)
	if err != nil {
		return fmt.Errorf("failed to delete type index: %w", err)
	}

	// Delete from expiry index if TTL was set
	if node.ExpiresAt != nil {
		expiryKey := utils.EncodeExpiryIndexKey(graphID, node.ID, *node.ExpiresAt)
		if err := t.delete(expiryKey); err != nil {
			// This is not a critical failure, as the node is being deleted anyway.
			// A log is sufficient.
			fmt.Printf("warn: failed to remove expiry index for node %s: %v\n", nodeID, err)
		}
	}

	// Delete outgoing edges
	outgoingPrefix := []byte(fmt.Sprintf("%sout:%s:%s:", utils.NodeIndexPrefix, graphID, nodeID))
	outIterOpts := badger.DefaultIteratorOptions
	outIter := t.txn.NewIterator(outIterOpts)
	defer outIter.Close()

	for outIter.Seek(outgoingPrefix); outIter.ValidForPrefix(outgoingPrefix); outIter.Next() {
		item := outIter.Item()
		err := item.Value(func(val []byte) error {
			return t.DeleteEdge(graphID, models.EdgeID(val))
		})
		if err != nil {
			return fmt.Errorf("failed to delete outgoing edge during node deletion: %w", err)
		}
	}

	// Delete incoming edges
	incomingPrefix := []byte(fmt.Sprintf("%sin:%s:%s:", utils.NodeIndexPrefix, graphID, nodeID))
	inIterOpts := badger.DefaultIteratorOptions
	inIter := t.txn.NewIterator(inIterOpts)
	defer inIter.Close()

	for inIter.Seek(incomingPrefix); inIter.ValidForPrefix(incomingPrefix); inIter.Next() {
		item := inIter.Item()
		err := item.Value(func(val []byte) error {
			return t.DeleteEdge(graphID, models.EdgeID(val))
		})
		if err != nil {
			return fmt.Errorf("failed to delete incoming edge during node deletion: %w", err)
		}
	}

	return nil
}
