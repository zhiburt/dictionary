package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dictionary/tgbot/dict/pb"
)

func TestParseWord(t *testing.T) {
	cases := []struct {
		Input    string
		Expected *pb.Word
	}{
		{
			"/add test [trans]",
			&pb.Word{Word: "test", Transcription: "trans"},
		},
		{
			"/add       test       [trans]",
			&pb.Word{Word: "test", Transcription: "trans"},
		},
		{
			"/add test[trans]",
			&pb.Word{Word: "test", Transcription: "trans"},
		},
		{
			"/add test",
			&pb.Word{Word: "test"},
		},
		{
			"/add              test",
			&pb.Word{Word: "test"},
		},
		{
			"/add test [trans]\nexample1\nexample2",
			&pb.Word{Word: "test", Transcription: "trans", Examples: []string{"example1", "example2"}},
		},
		{
			"/add test\nexample1\nexample2",
			&pb.Word{Word: "test", Examples: []string{"example1", "example2"}},
		},
		{
			"/add",
			nil,
		},
		{
			"/add [test]",
			nil,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			actual := ParseWord(c.Input)
			if reflect.DeepEqual(c.Expected, actual) == false {
				t.Errorf("expected:\n%v\nactual:\n%v\n", c.Expected, actual)
			}
		})
	}
}

func TestParseWordListWord(t *testing.T) {
	cases := []struct {
		Input    string
		Expected []*pb.Word
	}{
		{"/confuse_list abc\nqwe",
			[]*pb.Word{{Word: "abc"}, {Word: "qwe"}},
		},
		{"/confuse_list qwe",
			[]*pb.Word{{Word: "qwe"}},
		},
		{"/confuse_list qwe rew\nbbb",
			[]*pb.Word{{Word: "qwe rew"}, {Word: "bbb"}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			actual := ParseWordList(strings.SplitN(c.Input, "\n", -1))
			if reflect.DeepEqual(c.Expected, actual) == false {
				t.Errorf("expected:\n%v\nactual:\n%v\n", c.Expected, actual)
			}
		})
	}
}
