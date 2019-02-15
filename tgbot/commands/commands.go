package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dictionary/tgbot/commands/parser"
	"github.com/dictionary/tgbot/commands/shuffle"
	"github.com/dictionary/tgbot/dict/pb"
	"github.com/dictionary/tgbot/formatter"
)

var (
	// ErrParseParams default error for parse params
	ErrParseParams = errors.New("params don't contains enough information. You should look API up")
)

// Command it's command which exec some acting
type Command interface {
	Exec(context.Context, pb.DictionaryClient, ...string) (string, error)
	Name() string
}

// NewCommandByName return Command by name
func NewCommandByName(s string) Command {
	m := map[string]Command{
		"help":         HelpCommand{},
		"pattern":      PatternCommand{},
		"add":          AddWordCommand{},
		"words":        WordsCommand{},
		"confuse":      ConfuseCommand{},
		"confuse_list": ConfuseWordListCommand{},
	}

	if command, found := m[s]; found {
		return command
	}
	return NotFoundCommand{}
}

type (
	// HelpCommand implements Command
	// for some help imformation
	HelpCommand struct {
		Command
	}

	// PatternCommand implements Command
	// for example haw need to give information for AddWordCommand
	PatternCommand struct {
		Command
	}

	// AddWordCommand implements Command
	// add new word to service
	AddWordCommand struct {
		Command
	}

	// ConfuseWordListCommand implements Command
	// add new word to service
	ConfuseWordListCommand struct {
		Command
	}

	// WordsCommand implements Command
	// gets all words from service
	WordsCommand struct {
		Command
	}

	// ConfuseCommand implements Command
	// gets several <N> words from service
	ConfuseCommand struct {
		Command
	}

	// NotFoundCommand implements Command
	// gets message that command not found
	NotFoundCommand struct {
		Command
	}
)

// Exec do command
func (c HelpCommand) Exec(ctx context.Context, dict pb.DictionaryClient, params ...string) (string, error) {
	return "Hello, I'm bot which is being developed by Maxim Zhiburt", nil
}

// Name returns name of Command
// usually for logs
func (c HelpCommand) Name() string {
	return "Help"
}

// Exec do command
func (c PatternCommand) Exec(ctx context.Context, dict pb.DictionaryClient, params ...string) (string, error) {
	return "give me string like that:\nword [transcription]\nexample\nexample", nil
}

// Name returns name of Command
// usually for logs
func (c PatternCommand) Name() string {
	return "Pattern"
}

// Exec do command
func (c AddWordCommand) Exec(ctx context.Context, dict pb.DictionaryClient, params ...string) (string, error) {
	if len(params) < 1 {
		return "", ErrParseParams
	}

	w := parser.ParseWord(params[0])
	if w == nil {
		return "", ErrParseParams
	}

	req := &pb.AddNewWordRequest{Word: w}
	_, err := dict.AddNewWord(ctx, req)
	return "Success.I've added one", err
}

// Name returns name of Command
// usually for logs
func (c AddWordCommand) Name() string {
	return "AddWord"
}

// Exec do command
func (c ConfuseWordListCommand) Exec(ctx context.Context, dict pb.DictionaryClient, params ...string) (string, error) {
	if len(params) < 1 {
		return "", ErrParseParams
	}
	if parser.IsOnlyCommand(params[0]) {
		return "", ErrParseParams
	}

	lines := strings.SplitN(params[0], "\n", -1)
	parsedWords := parser.ParseWordList(lines)

	for _, word := range parsedWords {
		req := &pb.AddNewWordRequest{Word: word}
		dict.AddNewWord(ctx, req)
	}

	confused := shuffle.ConfuseWordsN(parsedWords, strconv.FormatInt(int64(len(parsedWords)), 10))

	return fmt.Sprintf("you've given %d, I confuse %d\n%s", len(lines), len(parsedWords), formatter.WordsToStringShort(confused)), nil
}

// Name returns name of Command
// usually for logs
func (c ConfuseWordListCommand) Name() string {
	return "AddWordList"
}

// Exec do command
func (c WordsCommand) Exec(ctx context.Context, dict pb.DictionaryClient, params ...string) (string, error) {
	if len(params) > 0 && params[0] != "" {
		if search := parser.FirstWordAfterCommand(params[0]); search != "" {
			word, err := dict.GetByW(ctx, &pb.GetByWRequest{W: search})
			if err != nil {
				return "", err
			}
			return formatter.WordToString(word.Word), nil
		}
	}

	words, err := dict.Words(ctx, &pb.WordsRequest{})
	if err != nil {
		return "", err
	}
	return formatter.WordsToStringShort(words.Words), nil
}

// Name returns name of Command
// usually for logs
func (c WordsCommand) Name() string {
	return "Words"
}

// Exec do command
func (c ConfuseCommand) Exec(ctx context.Context, dict pb.DictionaryClient, params ...string) (string, error) {
	if len(params) < 1 {
		return "", ErrParseParams
	}

	resp, err := dict.Words(ctx, &pb.WordsRequest{})
	if err != nil {
		return "", err
	}

	message := params[0]
	confusedWords := shuffle.ConfuseWordsN(resp.Words, message)
	return formatter.WordsToStringShort(confusedWords), nil
}

// Name returns name of Command
// usually for logs
func (c ConfuseCommand) Name() string {
	return "Confuse"
}

// Exec do command
func (c NotFoundCommand) Exec(ctx context.Context, dict pb.DictionaryClient, params ...string) (string, error) {
	return "command not found", nil
}

// Name returns name of Command
// usually for logs
func (c NotFoundCommand) Name() string {
	return "NotFound"
}
