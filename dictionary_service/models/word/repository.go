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
var ErrRepository = errors.New("Repository cannot do that you want")

// ErrWordNotFound erorr when word is not found
var ErrWordNotFound = errors.New("Repository cannot has found one")

// ErrDublicateWord erorr when word is not found
var ErrDublicateWord = errors.New("Repository has found the same word")

type Repository interface {
	Words(ctx context.Context) ([]Word, error)
	AddWordInto(ctx context.Context, word Word) error
	GetWordByID(ctx context.Context, id string) (Word, error)
	GetWordByW(ctx context.Context, w string) (Word, error)
}

// BadgerRepository is implementation Repository by Badger
type BadgerRepository struct {
	db     *badger.DB
	logger log.Logger
	sync.Mutex
}

// NewBadgerRepository creates and returns new repository
func NewBadgerRepository(dir string, logger log.Logger) Repository {
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

// Words
// todo: refactor
func (br *BadgerRepository) Words(ctx context.Context) ([]Word, error) {
	var words []Word
	err := br.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				w := Word{}
				err := unmarshalWord(v, &w)
				handleError(err)
				words = append(words, w)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return words, err
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

// GetWordByW gets word from badger by W
func (br *BadgerRepository) GetWordByW(ctx context.Context, word string) (Word, error) {
	w := Word{}
	err := br.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			val, err := it.Item().ValueCopy([]byte{})
			handleError(err)
			wCopy := Word{}
			handleError(unmarshalWord(val, &wCopy))
			if wCopy.W == word {
				w = wCopy
				return nil
			}
		}

		return ErrWordNotFound
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

	if ok, _ := containsWord(ctx, br, word.W); ok {
		return ErrDublicateWord
	}

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

func containsWord(ctx context.Context, br *BadgerRepository, w string) (bool, error) {
	words, err := br.Words(ctx)
	if err != nil {
		return false, err
	}

	for _, word := range words {
		if word.W == w {
			return true, nil
		}
	}

	return false, nil
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
