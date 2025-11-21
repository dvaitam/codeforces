package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type term struct {
	root   string
	weight int
}

var terms = map[int][]term{
	1: {
		{"cultur", 6}, {"art", 4}, {"museum", 4}, {"poet", 4}, {"music", 4},
		{"novel", 3}, {"book", 3}, {"literatur", 4}, {"theatre", 3}, {"theater", 3},
		{"movie", 3}, {"film", 3}, {"festival", 2}, {"paint", 3}, {"artist", 3},
		{"opera", 3}, {"gallery", 2}, {"sculpt", 2}, {"drama", 2}, {"ballet", 3},
		{"author", 2}, {"story", 2}, {"heritage", 3}, {"classic", 2},
	},
	2: {
		{"government", 5}, {"minister", 4}, {"president", 5}, {"parliament", 4}, {"policy", 3},
		{"elect", 4}, {"polit", 5}, {"state", 2}, {"party", 3}, {"law", 2},
		{"constitut", 4}, {"senate", 3}, {"congress", 3}, {"cabinet", 3}, {"kremlin", 5},
		{"duma", 4}, {"administrat", 3}, {"security", 2}, {"diplomat", 3}, {"military", 2},
		{"war", 2}, {"conflict", 2}, {"reform", 2}, {"prime", 3}, {"governor", 3},
		{"opposition", 3}, {"referendum", 3}, {"authority", 3}, {"campaign", 3}, {"policy", 3},
		{"official", 3}, {"regulation", 2}, {"senator", 3},
	},
	3: {
		{"econom", 6}, {"market", 5}, {"trade", 5}, {"company", 4}, {"business", 5},
		{"industr", 4}, {"invest", 4}, {"financ", 5}, {"bank", 4}, {"profit", 4},
		{"price", 3}, {"export", 4}, {"import", 4}, {"currency", 3}, {"capital", 3},
		{"product", 3}, {"factory", 3}, {"income", 3}, {"revenue", 4}, {"earn", 3},
		{"budget", 3}, {"inflation", 3}, {"tax", 3}, {"loan", 3}, {"credit", 3},
		{"retail", 2}, {"sale", 3}, {"stock", 4}, {"share", 3}, {"dollar", 4},
		{"euro", 4}, {"ruble", 3}, {"yen", 3}, {"gdp", 4}, {"oil", 3},
		{"gas", 3}, {"metal", 2}, {"energy", 3}, {"analyst", 2}, {"contract", 2},
		{"investor", 3}, {"portfolio", 2}, {"dividend", 3}, {"logistic", 2}, {"freight", 2},
	},
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	data, err := io.ReadAll(reader)
	if err != nil {
		return
	}
	text := strings.ToLower(string(data))
	words := tokenize(text)

	scores := map[int]int{1: 0, 2: 0, 3: 0}
	for subject, list := range terms {
		for _, t := range list {
			for _, w := range words {
				if strings.Contains(w, t.root) {
					scores[subject] += t.weight
				}
			}
		}
	}

	// Numeric-heavy texts usually describe economics/trade topics.
	digitScore := 0
	symbolScore := 0
	for _, ch := range text {
		if ch >= '0' && ch <= '9' {
			digitScore++
		}
		if ch == '$' || ch == '%' {
			symbolScore += 2
		}
	}
	if digitScore > len(words) {
		scores[3] += 3
	}
	if symbolScore > 0 {
		scores[3] += symbolScore
	}

	result := 3
	if scores[3] >= scores[2] && scores[3] >= scores[1] {
		result = 3
	} else if scores[2] >= scores[3] && scores[2] >= scores[1] {
		result = 2
	} else {
		result = 1
	}
	fmt.Println(result)
}

func tokenize(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}
