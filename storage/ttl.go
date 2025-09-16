package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/utils"
)

// TTLManager handles the expiration of nodes.
type TTLManager struct {
	engine *BadgerEngine
	stop   chan struct{}
}

// NewTTLManager creates a new TTL manager.
func NewTTLManager(engine *BadgerEngine) *TTLManager {
	return &TTLManager{
		engine: engine,
		stop:   make(chan struct{}),
	}
}

// Start begins the background TTL cleanup process.
func (tm *TTLManager) Start() {
	go tm.run()
}

// Stop halts the background TTL cleanup process.
func (tm *TTLManager) Stop() {
	close(tm.stop)
}

// run is the main loop for the TTL manager.
// Cleanup is an exported method to manually trigger a cleanup for testing.
func (tm *TTLManager) Cleanup() {
	tm.cleanupExpiredNodes()
}

func (tm *TTLManager) run() {
	ticker := time.NewTicker(1 * time.Minute) // Check for expired nodes every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tm.cleanupExpiredNodes()
		case <-tm.stop:
			return
		}
	}
}

// cleanupExpiredNodes scans for and deletes expired nodes.
func (tm *TTLManager) cleanupExpiredNodes() {
	var expiredKeys [][]byte
	prefix := utils.CreateExpiryIteratorPrefix()
	now := time.Now().UTC().Format(time.RFC3339)

	// Phase 1: Collect keys in a read-only transaction.
	tm.engine.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()
			keyStr := string(key)

			parts := strings.Split(keyStr, ":")
			if len(parts) >= 2 && parts[1] > now {
				break // Stop if we've passed the current time.
			}
			expiredKeys = append(expiredKeys, item.KeyCopy(nil))
		}
		return nil
	})

	// Phase 2: Delete the collected keys in a separate write transaction.
	if len(expiredKeys) > 0 {
		for _, key := range expiredKeys {
			graphID, nodeID := utils.DecodeExpiryIndexKey(key)
			if graphID != "" && nodeID != "" {
				// Each deletion gets its own transaction to ensure atomicity.
				if err := tm.engine.DeleteNode(graphID, nodeID); err != nil {
					fmt.Printf("warn: failed to delete expired node %s: %v\n", nodeID, err)
				}
			}
		}
	}
}

// AddNodeToExpiryIndex adds a node to the expiration index.
func (tm *TTLManager) AddNodeToExpiryIndex(txn *badger.Txn, graphID models.GraphID, node *models.Node) error {
	if node.ExpiresAt == nil {
		return nil
	}
	key := utils.EncodeExpiryIndexKey(graphID, node.ID, *node.ExpiresAt)
	return txn.Set(key, []byte(node.ID))
}

// RemoveNodeFromExpiryIndex removes a node from the expiration index.
func (tm *TTLManager) RemoveNodeFromExpiryIndex(txn *badger.Txn, graphID models.GraphID, node *models.Node) error {
	if node.ExpiresAt == nil {
		return nil
	}
	key := utils.EncodeExpiryIndexKey(graphID, node.ID, *node.ExpiresAt)
	return txn.Delete(key)
}
