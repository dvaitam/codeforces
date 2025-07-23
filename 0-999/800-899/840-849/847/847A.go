package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	l := make([]int, n+1)
	r := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &l[i], &r[i])
	}

	var heads []int
	for i := 1; i <= n; i++ {
		if l[i] == 0 {
			heads = append(heads, i)
		}
	}

	// connect lists sequentially
	for i := 0; i < len(heads)-1; i++ {
		tail := heads[i]
		for r[tail] != 0 {
			tail = r[tail]
		}
		nextHead := heads[i+1]
		r[tail] = nextHead
		l[nextHead] = tail
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i := 1; i <= n; i++ {
		fmt.Fprintf(out, "%d %d\n", l[i], r[i])
	}
}
