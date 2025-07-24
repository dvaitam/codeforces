package main

import (
	"bufio"
	"fmt"
	"os"
)

func bullsAndCows(a, b string) (int, int) {
	bulls := 0
	for i := 0; i < 4; i++ {
		if a[i] == b[i] {
			bulls++
		}
	}
	var usedA [10]bool
	var usedB [10]bool
	for i := 0; i < 4; i++ {
		usedA[a[i]-'0'] = true
		usedB[b[i]-'0'] = true
	}
	cows := 0
	for d := 0; d < 10; d++ {
		if usedA[d] && usedB[d] {
			cows++
		}
	}
	cows -= bulls
	return bulls, cows
}

func main() {
	writer := bufio.NewWriter(os.Stdout)
	reader := bufio.NewReader(os.Stdin)

	// generate all possible codes (permutations of 4 distinct digits)
	codes := make([]string, 0, 5040)
	for d1 := byte('0'); d1 <= '9'; d1++ {
		for d2 := byte('0'); d2 <= '9'; d2++ {
			if d2 == d1 {
				continue
			}
			for d3 := byte('0'); d3 <= '9'; d3++ {
				if d3 == d1 || d3 == d2 {
					continue
				}
				for d4 := byte('0'); d4 <= '9'; d4++ {
					if d4 == d1 || d4 == d2 || d4 == d3 {
						continue
					}
					codes = append(codes, string([]byte{d1, d2, d3, d4}))
				}
			}
		}
	}

	remaining := codes
	for q := 0; q < 7 && len(remaining) > 0; q++ {
		guess := remaining[0]
		fmt.Fprintln(writer, guess)
		writer.Flush()

		var bulls, cows int
		if _, err := fmt.Fscan(reader, &bulls, &cows); err != nil {
			return
		}
		if bulls == 4 {
			return
		}
		next := make([]string, 0, len(remaining))
		for _, code := range remaining {
			b, c := bullsAndCows(guess, code)
			if b == bulls && c == cows {
				next = append(next, code)
			}
		}
		remaining = next
	}
}
