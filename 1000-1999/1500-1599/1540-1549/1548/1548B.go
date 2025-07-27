package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

type pair struct {
	g int64
	l int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if n == 1 {
			fmt.Fprintln(out, 1)
			continue
		}
		diff := make([]int64, n-1)
		for i := 1; i < n; i++ {
			d := a[i] - a[i-1]
			if d < 0 {
				d = -d
			}
			diff[i-1] = d
		}
		ans := 1
		cur := make([]pair, 0)
		for _, d := range diff {
			next := make([]pair, 0, len(cur)+1)
			if d > 1 {
				next = append(next, pair{d, 1})
				if ans < 2 {
					ans = 2
				}
			}
			for _, p := range cur {
				g := gcd(p.g, d)
				if g <= 1 {
					continue
				}
				l := p.l + 1
				if len(next) > 0 && next[len(next)-1].g == g {
					if l > next[len(next)-1].l {
						next[len(next)-1].l = l
					}
				} else {
					next = append(next, pair{g, l})
				}
				if ans < l+1 {
					ans = l + 1
				}
			}
			cur = next
		}
		fmt.Fprintln(out, ans)
	}
}
