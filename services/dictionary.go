package services

import (
	"context"
	"time"

	"github.com/dictionary/models/word"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
)

// Dictionary service discribes dictionary
type Dictionary interface {
	AddNewWord(ctx context.Context, w word.Word) (string, error)
	GetByID(ctx context.Context, id string) (word.Word, error)
	GetByW(ctx context.Context, w string) (word.Word, error)
	Words(ctx context.Context) ([]word.Word, error)
}

// dictionary implements the Dictionary Service
type dictionary struct {
	logger     log.Logger
	repository word.Repository
}

// NewDictionary create and return new Dictionary
func NewDictionary(rep word.Repository, logger log.Logger) Dictionary {
	return &dictionary{
		logger:     logger,
		repository: rep,
	}
}

func (d *dictionary) AddNewWord(ctx context.Context, w word.Word) (string, error) {
	logger := log.With(d.logger, "method", "AddNewWord")
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	w.ID = id
	w.Timestamp = time.Now().Unix()

	if err := d.repository.AddWordInto(ctx, w); err != nil {
		level.Error(logger).Log("error", err)
		return "", word.ErrRepository
	}
	return id, nil
}

func (d *dictionary) GetByID(ctx context.Context, id string) (word.Word, error) {
	logger := log.With(d.logger, "method", "GetByID")

	w, err := d.repository.GetWordByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("erorr", err)
		return w, err
	}
	return w, nil
}

func (d *dictionary) GetByW(ctx context.Context, w string) (word.Word, error) {
	logger := log.With(d.logger, "method", "GetByID")

	word, err := d.repository.GetWordByW(ctx, w)
	if err != nil {
		level.Error(logger).Log("erorr", err)
		return word, err
	}
	return word, nil
}

func (d *dictionary) Words(ctx context.Context) ([]word.Word, error) {
	logger := log.With(d.logger, "method", "GetByID")

	w, err := d.repository.Words(ctx)
	if err != nil {
		level.Error(logger).Log("erorr", err)
		return w, err
	}
	return w, nil
}
