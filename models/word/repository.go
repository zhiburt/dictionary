package word

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/go-kit/kit/log/level"

	"github.com/dgraph-io/badger"
	"github.com/go-kit/kit/log"
)

// ErrRepository is the default erorr in repo
var ErrRepository = errors.New("Repository cannot do that you want0")

// ErrWordNotFound erorr when word is not found
var ErrWordNotFound = errors.New("Repository cannot do that you want0")

// Repository wrapp db
type Repository interface {
	AddWordInto(ctx context.Context, word Word) error
	GetWordByID(ctx context.Context, id string) (Word, error)
}

// BadgerRepository is implementation Repository by Badger
type BadgerRepository struct {
	db     *badger.DB
	logger log.Logger
	sync.Mutex
}

// NewBadgerRepository creates and returns new repository
func NewBadgerRepository(dir string, logger log.Logger) Repository {
	badger.UseDefaultLogger()
	opt := badger.DefaultOptions
	opt.Dir = dir
	opt.ValueDir = dir
	db, err := badger.Open(opt)
	if err != nil {
		panic(err)
	}

	return &BadgerRepository{
		db:     db,
		logger: logger,
	}
}

// GetWordByID gets word from badger by ID
func (br *BadgerRepository) GetWordByID(ctx context.Context, id string) (Word, error) {
	w := Word{}
	err := br.db.View(func(txn *badger.Txn) error {
		val, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		valCopy, err := val.ValueCopy(nil)
		handleError(err)

		handleError(unmarshalWord(valCopy, &w))

		return nil
	})

	if err != nil {
		level.Warn(br.logger).Log("method", "GetWordByID")
		return w, ErrWordNotFound
	}

	return w, nil
}

// AddWordInto adds word into badger if this one doesn't exiest there
func (br *BadgerRepository) AddWordInto(ctx context.Context, word Word) error {
	txn := br.db.NewTransaction(true)
	b, err := marshalWord(word)
	handleError(err)
	if err = txn.Set([]byte(word.ID), b); err != nil {
		level.Error(br.logger).Log("method", "AddWordInto")
		return ErrRepository
	}
	if err = txn.Commit(); err != nil {
		level.Error(br.logger).Log("method", "AddWordInto")
		return ErrRepository
	}

	return nil
}

func marshalWord(w Word) ([]byte, error) {
	return json.Marshal(w)
}

func unmarshalWord(b []byte, w *Word) error {
	return json.Unmarshal(b, w)
}

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}
