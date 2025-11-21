package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	waiting := 0
	results := make([]string, 0)

	for i := 0; i < n; i++ {
		var typ string
		fmt.Fscan(in, &typ)
		if typ == "P" {
			var p int
			fmt.Fscan(in, &p)
			waiting += p
		} else if typ == "B" {
			var b int
			fmt.Fscan(in, &b)
			board := b
			if waiting < board {
				board = waiting
			}
			waiting -= board
			if b > board {
				results = append(results, "YES")
			} else {
				results = append(results, "NO")
			}
		}
	}

	for i, r := range results {
		if i > 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprint(out, r)
	}
}
