package main

import (
	"bufio"
	"fmt"
	"os"
)

func orRange(k int) int {
	res := 0
	for i := 0; i < k; i++ {
		res |= i
	}
	return res
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

		pos := 0
		for pos < 31 && ((x>>pos)&1) == 1 {
			pos++
		}
		m0 := 1 << pos
		k := m0
		if k > n {
			k = n
		}
		if k == n {
			if orRange(k) != x {
				k--
			}
		}

		ans := make([]int, 0, n)
		curOr := 0
		for i := 0; i < k; i++ {
			ans = append(ans, i)
			curOr |= i
		}

		remaining := n - len(ans)
		if remaining > 0 {
			missing := x &^ curOr
			if missing != 0 {
				ans = append(ans, missing)
				curOr |= missing
				remaining--
			}
			for remaining > 0 {
				ans = append(ans, 0)
				remaining--
			}
		}

		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
