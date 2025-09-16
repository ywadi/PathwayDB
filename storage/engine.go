package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v3"
)

// BadgerEngine implements the StorageEngine interface using Badger v3
type BadgerEngine struct {
	db         *badger.DB
	path       string
	ttlManager *TTLManager
}

// NewBadgerEngine creates a new BadgerEngine instance
func NewBadgerEngine() *BadgerEngine {
	engine := &BadgerEngine{}
	engine.ttlManager = NewTTLManager(engine)
	return engine
}

// Open initializes the Badger database
func (e *BadgerEngine) Open(path string) error {
	e.path = path
	
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // Disable badger logging for cleaner output
	
	var err error
	e.db, err = badger.Open(opts)
	if err != nil {
		return fmt.Errorf("failed to open badger database: %w", err)
	}
	
	log.Printf("Badger database opened at: %s", path)

	// Start the TTL manager
	e.ttlManager.Start()

	return nil
}

// Close closes the Badger database
func (e *BadgerEngine) Close() error {
	// Stop the TTL manager first
	if e.ttlManager != nil {
		e.ttlManager.Stop()
	}

	if e.db != nil {
		err := e.db.Close()
		if err != nil {
			return fmt.Errorf("failed to close badger database: %w", err)
		}
		log.Printf("Badger database closed")
	}
	return nil
}

// Backup creates a backup of the database
// Cleanup is a test helper to manually trigger TTL cleanup.
func (e *BadgerEngine) Cleanup() {
	if e.ttlManager != nil {
		e.ttlManager.Cleanup()
	}
}

func (e *BadgerEngine) Backup(backupPath string) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}
	
	backupFile := filepath.Join(backupPath, "backup.db")
	f, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer f.Close()
	
	_, err = e.db.Backup(f, 0)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	
	log.Printf("Database backup created at: %s", backupFile)
	return nil
}

// RunTransaction executes a function within a Badger transaction
func (e *BadgerEngine) RunTransaction(fn TransactionFunc) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}
	
	return e.db.Update(func(txn *badger.Txn) error {
		tx := &BadgerTransaction{txn: txn}
		return fn(tx)
	})
}

// RunReadOnlyTransaction executes a read-only function within a Badger transaction
func (e *BadgerEngine) RunReadOnlyTransaction(fn func(*badger.Txn) error) error {
	if e.db == nil {
		return fmt.Errorf("database not opened")
	}
	
	return e.db.View(fn)
}

// set is a helper method to set a key-value pair
func (e *BadgerEngine) set(key []byte, value []byte) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// get is a helper method to get a value by key
func (e *BadgerEngine) get(key []byte) ([]byte, error) {
	var value []byte
	err := e.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		
		value, err = item.ValueCopy(nil)
		return err
	})
	
	return value, err
}

// delete is a helper method to delete a key
func (e *BadgerEngine) delete(key []byte) error {
	return e.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// iterateWithPrefix iterates over keys with a given prefix
func (e *BadgerEngine) iterateWithPrefix(prefix []byte, fn func(key []byte, value []byte) error) error {
	return e.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := item.Key()
			
			err := item.Value(func(value []byte) error {
				return fn(key, value)
			})
			if err != nil {
				return err
			}
		}
		
		return nil
	})
}

// BadgerTransaction wraps a Badger transaction to implement the Transaction interface
type BadgerTransaction struct {
	txn *badger.Txn
}

// Commit commits the transaction
func (t *BadgerTransaction) Commit() error {
	return t.txn.Commit()
}

// Discard discards the transaction
func (t *BadgerTransaction) Discard() {
	t.txn.Discard()
}

// set is a helper method for setting values within a transaction
func (t *BadgerTransaction) set(key []byte, value []byte) error {
	return t.txn.Set(key, value)
}

// setWithTTL is a helper method for setting values with a TTL within a transaction
func (t *BadgerTransaction) setWithTTL(key []byte, value []byte, ttl time.Duration) error {
	e := badger.NewEntry(key, value).WithTTL(ttl)
	return t.txn.SetEntry(e)
}

// get is a helper method for getting values within a transaction
func (t *BadgerTransaction) get(key []byte) ([]byte, error) {
	item, err := t.txn.Get(key)
	if err != nil {
		return nil, err
	}
	
	return item.ValueCopy(nil)
}

// delete is a helper method for deleting keys within a transaction
func (t *BadgerTransaction) delete(key []byte) error {
	return t.txn.Delete(key)
}
