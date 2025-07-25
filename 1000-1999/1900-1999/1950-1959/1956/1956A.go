package main

import (
	"bufio"
	"fmt"
	"os"
)

func winners(n int, a []int) int {
	for n >= a[0] {
		cnt := 0
		for _, v := range a {
			if v <= n {
				cnt++
			} else {
				break
			}
		}
		if cnt == 0 {
			break
		}
		n -= cnt
	}
	return n
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var k, q int
		fmt.Fscan(in, &k, &q)
		a := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(in, &a[i])
		}
		res := make([]int, q)
		for i := 0; i < q; i++ {
			fmt.Fscan(in, &res[i])
			res[i] = winners(res[i], a)
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
