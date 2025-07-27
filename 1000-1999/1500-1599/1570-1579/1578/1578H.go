package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	s   string
	pos int
)

// parseType parses a type expression starting at position pos
// according to the grammar described in problemH.txt and
// returns its order.
func parseType() int {
	left := parseFactor()
	if pos+1 < len(s) && s[pos] == '-' && s[pos+1] == '>' {
		pos += 2
		right := parseType()
		if left+1 > right {
			return left + 1
		}
		return right
	}
	return left
}

// parseFactor parses either the unit type "()" or a parenthesized
// type expression.
func parseFactor() int {
	// assume s[pos] == '(' according to grammar
	if pos+1 < len(s) && s[pos] == '(' && s[pos+1] == ')' {
		pos += 2
		return 0
	}
	// '(' T ')'
	pos++ // consume '('
	val := parseType()
	pos++ // consume ')'
	return val
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &s)
	pos = 0
	ans := parseType()
	fmt.Println(ans)
}
