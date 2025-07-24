package main

import (
	"bufio"
	"fmt"
	"os"
)

func buildFunnyPerm(n, x int) []int {
	if n%x != 0 {
		return nil
	}
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	p[1] = x
	p[n] = 1
	cur := x
	for i := x + 1; i < n; i++ {
		if i%cur == 0 && n%i == 0 {
			p[cur], p[i] = p[i], p[cur]
			cur = i
		}
	}
	if cur != n {
		p[cur] = n
	}
	return p[1:]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, x int
		fmt.Fscan(in, &n, &x)
		res := buildFunnyPerm(n, x)
		if res == nil {
			fmt.Fprintln(out, -1)
			continue
		}
		for i, v := range res {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, v)
		}
		out.WriteByte('\n')
	}
}
