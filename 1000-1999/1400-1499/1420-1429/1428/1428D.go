package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	r, c int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	row := n
	ans := make([]pair, 0, 2*n)
	ones := make([]pair, 0) // columns with a==1, waiting to be paired
	twos := make([]pair, 0) // columns that already have two hits, can serve for third

	ok := true
	for col := n; col >= 1 && ok; col-- {
		switch a[col] {
		case 0:
			// nothing to do
		case 1:
			if row == 0 {
				ok = false
				break
			}
			ans = append(ans, pair{row, col})
			ones = append(ones, pair{row, col})
			row--
		case 2:
			if len(ones) == 0 {
				ok = false
				break
			}
			p := ones[len(ones)-1]
			ones = ones[:len(ones)-1]
			ans = append(ans, pair{p.r, col})
			twos = append(twos, pair{p.r, col})
		case 3:
			if row == 0 {
				ok = false
				break
			}
			var p pair
			if len(twos) > 0 {
				p = twos[len(twos)-1]
				twos = twos[:len(twos)-1]
			} else if len(ones) > 0 {
				p = ones[len(ones)-1]
				ones = ones[:len(ones)-1]
			} else {
				ok = false
				break
			}
			ans = append(ans, pair{row, col})
			ans = append(ans, pair{row, p.c})
			twos = append(twos, pair{row, col})
			row--
		default:
			ok = false
		}
	}

	if !ok {
		fmt.Println(-1)
		return
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, len(ans))
	for _, pr := range ans {
		fmt.Fprintf(out, "%d %d\n", pr.r, pr.c)
	}
	out.Flush()
}
