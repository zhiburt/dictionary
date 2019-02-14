package word

import (
	"context"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/mock"
)

type MockBadgerRepo struct {
	mock.Mock
}

func (m *MockBadgerRepo) View(f func(*badger.Txn) error) error {
	args := m.Called(f)
	return args.Error(0)
}

func (m *MockBadgerRepo) NewTransaction(b bool) *badger.Txn {
	return &badger.Txn{}
}

func TestWords(t *testing.T) {
	b := BadgerRepository{}
	b.logger = log.NewNopLogger()

	testcases := []struct {
		param    error
		expected error
	}{
		{nil, nil},
		{ErrRepository, ErrRepository},
	}

	t.Run("", func(t *testing.T) {
		for _, c := range testcases {
			m := &MockBadgerRepo{}
			m.On("View", mock.AnythingOfType("func(*badger.Txn) error")).Return(c.param)
			b.db = m

			if _, err := b.Words(context.Background()); err != c.expected {
				t.Error()
			}
		}
	})
}
