package parser

import (
	"regexp"
	"strings"

	"github.com/dictionary/tgbot/dict/pb"
)

// FirstWordAfterCommand returns first word after command(/asd)
// if there's no word sreturns ""
func FirstWordAfterCommand(message string) string {
	r := regexp.MustCompile(`/\w+\s+(\w+)`)
	if r.MatchString(message) {
		return r.ReplaceAllString(message, "$1")
	}

	return ""
}

// ParseWordList parses lines such as
// /command word1
// word2
// word3
// and return all these words
//
// todo: if there're no words return command
func ParseWordList(lines []string) []*pb.Word {
	var words []*pb.Word
	r := regexp.MustCompile(`(?:/\w+\s+)?([\w\s]+)`)
	for _, line := range lines {
		if r.MatchString(line) {
			word := r.ReplaceAllString(line, "$1")
			words = append(words, &pb.Word{Word: word})
		}
	}

	return words
}

// IsOnlyCommand checks if exists a word after command
func IsOnlyCommand(s string) bool {
	r := regexp.MustCompile(`/\w+\s+$`)
	return r.MatchString(s)
}

// ParseWord parse string like that
// /command word [transcription]
// example1
// example2
// and return pb.Word
func ParseWord(message string) *pb.Word {
	lines := strings.SplitN(message, "\n", -1)
	if len(lines) == 0 {
		return nil
	}

	wordAndTranscription := parsedWordAndTranscription(lines[0])
	if wordAndTranscription == nil {
		return nil
	}

	w := &pb.Word{}
	w.Word = wordAndTranscription[0]
	if wordAndTranscription[1] != "" {
		w.Transcription = trimBracets(wordAndTranscription[1])
	}
	if len(lines) > 1 {
		w.Examples = lines[1:]
	}

	return w
}

func parsedWordAndTranscription(firstLine string) []string {
	pattern := regexp.MustCompile(`/\w+\s+(?P<word>\w+)\s*(?P<trans>\[[\wʃʊθʧʒʤðɑɪeəɛɔɪɑ:ʌæ\']+\])?`)
	if pattern.MatchString(firstLine) == false {
		return nil
	}

	result := []byte{}
	submatches := pattern.FindAllStringSubmatchIndex(firstLine, -1)[0]
	result = pattern.ExpandString(result, "$word:$trans", firstLine, submatches)
	return strings.Split(string(result), ":")
}

func trimBracets(s string) string {
	return s[1 : len(s)-1]
}
