package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func query(w *bufio.Writer, r *bufio.Reader, s string) int {
	fmt.Fprintln(w, s)
	w.Flush()
	var res int
	fmt.Fscan(r, &res)
	if res == -1 {
		os.Exit(0)
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	sA := strings.Repeat("a", 300)
	dA := query(out, in, sA)

	sB := strings.Repeat("b", 300)
	dB := query(out, in, sB)

	n := 600 - dA - dB
	totalB := dA - (300 - n)

	guess := make([]byte, n)
	for i := range guess {
		guess[i] = 'a'
	}
	dist := totalB

	for i := 0; i < n && dist > 0; i++ {
		guess[i] = 'b'
		cur := query(out, in, string(guess))
		if cur < dist {
			dist = cur
		} else {
			guess[i] = 'a'
		}
	}

	if dist != 0 {
		query(out, in, string(guess))
	}
}
