package main

import (
	"bufio"
	"fmt"
	"os"
)

func makePerm(n int) ([]int, bool) {
	if n == 1 {
		return []int{1}, true
	}
	if n <= 3 {
		return nil, false
	}
	res := make([]int, 0, n)
	for i := 2; i <= n; i += 2 {
		res = append(res, i)
	}
	for i := 1; i <= n; i += 2 {
		res = append(res, i)
	}
	return res, true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	if n == 1 && m == 1 {
		fmt.Fprintln(out, "YES")
		fmt.Fprintln(out, 1)
		return
	}
	pr, ok1 := makePerm(n)
	pc, ok2 := makePerm(m)
	if !ok1 || !ok2 {
		fmt.Fprintln(out, "NO")
		return
	}
	fmt.Fprintln(out, "YES")
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, (pr[i]-1)*m+pc[j])
		}
		fmt.Fprintln(out)
	}
}
