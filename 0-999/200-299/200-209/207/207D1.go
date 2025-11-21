package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type keyword struct {
	token  string
	weight int
}

var topicKeywords = map[int][]keyword{
	1: {
		{"culture", 3},
		{"cultural", 3},
		{"art", 2},
		{"arts", 2},
		{"artist", 2},
		{"artists", 2},
		{"painting", 2},
		{"paint", 1},
		{"music", 2},
		{"musical", 2},
		{"film", 2},
		{"cinema", 2},
		{"movie", 2},
		{"theatre", 2},
		{"theater", 2},
		{"literature", 2},
		{"literary", 2},
		{"book", 1},
		{"books", 1},
		{"poet", 2},
		{"poetry", 2},
		{"museum", 2},
		{"festival", 1},
		{"novel", 1},
		{"opera", 2},
		{"history", 1},
	},
	2: {
		{"politic", 3},
		{"government", 3},
		{"minister", 2},
		{"president", 2},
		{"parliament", 2},
		{"policy", 2},
		{"policies", 2},
		{"law", 2},
		{"state", 1},
		{"states", 1},
		{"election", 3},
		{"party", 2},
		{"parties", 2},
		{"diplomat", 2},
		{"diplomacy", 2},
		{"ministerial", 2},
		{"military", 2},
		{"army", 1},
		{"security", 1},
		{"defence", 1},
		{"defense", 1},
		{"war", 1},
		{"senate", 2},
		{"congress", 2},
		{"cabinet", 1},
		{"reform", 1},
		{"governor", 1},
	},
	3: {
		{"trade", 3},
		{"market", 3},
		{"markets", 3},
		{"econom", 3},
		{"business", 3},
		{"company", 2},
		{"companies", 2},
		{"industry", 2},
		{"industr", 2},
		{"profit", 2},
		{"profits", 2},
		{"bank", 2},
		{"finance", 2},
		{"financial", 2},
		{"investment", 3},
		{"invest", 3},
		{"import", 2},
		{"export", 2},
		{"stock", 2},
		{"stocks", 2},
		{"dollar", 1},
		{"currency", 2},
		{"sale", 1},
		{"sales", 1},
		{"production", 1},
		{"product", 1},
		{"supply", 1},
		{"demand", 1},
		{"budget", 1},
		{"tax", 1},
		{"corporate", 1},
		{"firm", 1},
		{"firms", 1},
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 0, 1<<20), 1<<25)

	var builder strings.Builder
	firstLine := true
	for scanner.Scan() {
		if !firstLine {
			builder.WriteByte(' ')
		}
		builder.WriteString(scanner.Text())
		firstLine = false
	}

	text := strings.ToLower(builder.String())
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	scores := map[int]int{1: 0, 2: 0, 3: 0}
	for _, w := range words {
		if len(w) == 0 {
			continue
		}
		for label, list := range topicKeywords {
			for _, kw := range list {
				if strings.Contains(w, kw.token) {
					scores[label] += kw.weight
				}
			}
		}
	}

	if scores[1] == 0 && scores[2] == 0 && scores[3] == 0 {
		for label, list := range topicKeywords {
			for _, kw := range list {
				if strings.Contains(text, kw.token) {
					scores[label] += kw.weight
				}
			}
		}
	}

	result := 1
	for _, label := range []int{2, 3} {
		if scores[label] > scores[result] {
			result = label
		}
	}

	fmt.Println(result)
}
