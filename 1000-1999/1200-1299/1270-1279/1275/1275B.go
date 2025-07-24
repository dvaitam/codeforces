package main

import (
	"bufio"
	"fmt"
	"os"
)

type commit struct {
	author int
	hash   string
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	stack := make([]commit, 0, n)
	for i := 0; i < n; i++ {
		var dev int
		var h string
		fmt.Fscan(in, &dev, &h)
		if len(stack) > 0 {
			last := stack[len(stack)-1]
			if last.author != dev {
				// review the last unreviewed commit
				stack = stack[:len(stack)-1]
			}
		}
		stack = append(stack, commit{author: dev, hash: h})
	}

	for _, c := range stack {
		fmt.Fprintln(out, c.hash)
	}
}
