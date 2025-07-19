package main

import (
	"bufio"
	"fmt"
	"os"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	// Read n*n gcd table entries
	S := make(map[int]int)
	for i := 0; i < n*n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		S[x]++
	}

	ans := make([]int, 0, n)
	tmp := make([]int, 0, n)
	prev := make([]int, 0, n)
	for len(ans) < n {
		// Remove counts of gcds with previously chosen values
		for _, t := range tmp {
			for _, a := range ans {
				d := gcd(t, a)
				S[d]--
			}
			for _, p := range prev {
				d := gcd(t, p)
				S[d]--
			}
		}
		tmp = tmp[:0]

		// Copy current answers to prev
		prev = append(prev[:0], ans...)

		// Find the largest remaining number
		maxk := 0
		for k, v := range S {
			if v > 0 && k > maxk {
				maxk = k
			}
		}
		tmp = append(tmp, maxk)
		ans = append(ans, maxk)
	}

	// Output the reconstructed array
	for i, v := range ans {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, v)
	}
	writer.WriteByte('\n')
}
