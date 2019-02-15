package formatter

import (
	"fmt"
	"strings"

	"github.com/dictionary/tgbot/dict/pb"
)

// WordToString format simple word to string
func WordToString(w *pb.Word) string {
	builder := strings.Builder{}
	builder.WriteString(w.Word)
	if w.Transcription != "" {
		builder.WriteString(" [" + w.Transcription + "]")
	}
	builder.WriteRune('\n')
	if w.Examples != nil {
		for _, example := range w.Examples {
			builder.WriteString(example + "\n")
		}
	}

	return builder.String()
}

// WordsToStringShort format slice to string
func WordsToStringShort(words []*pb.Word) string {
	b := &strings.Builder{}
	for _, w := range words {
		fmt.Fprintf(b, "%s\n", w.Word)
	}
	if len(words) > 0 {
		fmt.Fprintf(b, "----[%d]\n", len(words))
	}

	return b.String()
}
