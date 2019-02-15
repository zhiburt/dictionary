package wrapbadger

import "github.com/dgraph-io/badger"

// DB interface budger for mocking in tests
type DB interface {
	View(func(*badger.Txn) error) error
	NewTransaction(bool) *badger.Txn
}
