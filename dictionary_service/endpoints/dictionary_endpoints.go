package endpoints

import "github.com/dictionary/dictionary_service/models/word"

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

// WordsResponse holds the response values for the WordsResponse method.
type WordsResponse struct {
	Words []word.Word `json:"words"`
	Err   error       `json:"error,omitempty"`
}

// WordsRequest holds the response values for the WordsRequest method.
type WordsRequest struct {
}

// GetByWRequest holds the request parameters for the GetByW method.
type GetByWRequest struct {
	W string
}

// GetByWResponse holds the response values for the GetByW method.
type GetByWResponse struct {
	Word word.Word `json:"word"`
	Err  error     `json:"error,omitempty"`
}
