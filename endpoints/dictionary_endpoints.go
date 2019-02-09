package endpoints

import "github.com/dictionary/models/word"

// AddNewWordRequest holds the request parameters for the AddNewWord method.
type AddNewWordRequest struct {
	Word word.Word
}

// AddNewWordResponse holds the response values for the AddNewWord method.
type AddNewWordResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

// GetByIDRequest holds the request parameters for the GetByID method.
type GetByIDRequest struct {
	ID string
}

// GetByIDResponse holds the response values for the GetByID method.
type GetByIDResponse struct {
	Word word.Word `json:"word"`
	Err  error     `json:"error,omitempty"`
}