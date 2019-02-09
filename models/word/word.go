package word

// Word represents a word
type Word struct {
	ID            string   `json:"id,omitempty"`
	W             string   `json:"text"`
	Transcription string   `json:"transcription"`
	Examples      []string `json:"examples,omitempty"`
	Timestamp     int64    `json:"timestamp,omitempty"`
}
