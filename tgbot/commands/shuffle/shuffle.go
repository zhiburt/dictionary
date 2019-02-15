package shuffle

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/dictionary/tgbot/dict/pb"
)

// ConfuseWordsN gets N or less if exists less than N and
// confuse ones on random order and
// returns those words
func ConfuseWordsN(words []*pb.Word, msg string) []*pb.Word {
	var resp []*pb.Word
	m := make(map[int]struct{})
	n := amauntWordsNeed(msg)
	if n > uint64(len(words)) {
		n = uint64(len(words))
	}

	for i := uint64(0); i < n; i++ {
		r := randIntn(m, len(words))
		resp = append(resp, words[r])
	}

	return resp
}

func amauntWordsNeed(msg string) uint64 {
	var n uint64
	r := regexp.MustCompile(`(?:/[\w]+)?\s*(\d+)`)
	if r.MatchString(msg) {
		msg = r.ReplaceAllString(msg, "$1")
		n, _ = strconv.ParseUint(msg, 10, 64)
	}

	return n
}

func randIntn(m map[int]struct{}, l int) int {
	for {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(l)
		if _, ok := m[n]; !ok {
			m[n] = struct{}{}
			return n
		}
	}
}
