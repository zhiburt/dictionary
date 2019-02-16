package word

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
)

func TestAddWordInto(t *testing.T) {
	testcases := []struct {
		existWords []Word
		word       Word
		expected   error
	}{
		{
			[]Word{
				Word{ID: "1", W: "1"},
				Word{ID: "2", W: "2"},
			},
			Word{ID: "3", W: "3"},
			nil,
		},
		{
			[]Word{
				Word{ID: "1", W: "1"},
			},
			Word{ID: "2", W: "1"},
			ErrDublicateWord,
		},
		{nil, Word{}, ErrRepository},
	}

	t.Run("", func(t *testing.T) {
		for _, c := range testcases {
			dir, err := ioutil.TempDir(".", "badger-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)
			repo := confBadgerRepo(t, dir, c.existWords)

			if err := repo.AddWordInto(context.Background(), c.word); c.expected != err {
				t.Errorf("expected error %v but was\nactual %v\n", c.expected, err)
			}
		}
	})
}

func TestWords(t *testing.T) {
	testcases := []struct {
		addedWords    []Word
		wordsExpected []Word
		errExpected   error
	}{
		{
			[]Word{
				Word{ID: "1", W: "1"},
				Word{ID: "2", W: "2"},
			},
			[]Word{
				Word{ID: "1", W: "1"},
				Word{ID: "2", W: "2"},
			},
			nil,
		},
		{nil, nil, nil},
	}

	t.Run("", func(t *testing.T) {
		for _, c := range testcases {
			dir, err := ioutil.TempDir(".", "badger-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)
			repo := confBadgerRepo(t, dir, c.addedWords)

			words, err := repo.Words(context.Background())
			if err != c.errExpected {
				t.Errorf("expected error %v but was\nactual %v\n", c.errExpected, err)
			}
			if !arraySortedEqual(words, c.wordsExpected) {
				t.Errorf("expected words %v but was\nactual %v\n", c.wordsExpected, words)
			}
		}
	})
}

func confBadgerRepo(t *testing.T, dir string, words []Word) *BadgerRepository {
	b := &BadgerRepository{}
	opt := badger.DefaultOptions
	opt.Dir = dir
	opt.ValueDir = dir
	opt.Logger = nil

	db, err := badger.Open(opt)
	if err != nil {
		panic(err)
	}
	b.logger = log.NewNopLogger()
	b.db = db

	for _, w := range words {
		require.NoError(t, b.AddWordInto(context.Background(), w))
	}

	return b
}

func arraySortedEqual(a, b []Word) bool {
	if len(a) != len(b) {
		return false
	}

	acpy := make([]Word, len(a))
	bcpy := make([]Word, len(b))

	copy(acpy, a)
	copy(bcpy, b)

	sort.Slice(acpy, func(i, j int) bool { return acpy[i].ID < acpy[j].ID })
	sort.Slice(bcpy, func(i, j int) bool { return bcpy[i].ID < bcpy[j].ID })

	return reflect.DeepEqual(acpy, bcpy)
}
