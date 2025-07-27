package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxVal = 1000000

var spf [maxVal + 1]int

func init() {
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxVal; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
}

func squareFree(x int) int {
	res := 1
	for x > 1 {
		p := spf[x]
		if p == 0 {
			p = x
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt++
		}
		if cnt%2 == 1 {
			res *= p
		}
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
		var n int
		fmt.Fscan(in, &n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			f := squareFree(x)
			freq[f]++
		}
		ans0 := 0
		merge := 0
		for v, c := range freq {
			if c > ans0 {
				ans0 = c
			}
			if v == 1 || c%2 == 0 {
				merge += c
			}
		}
		ans1 := ans0
		if merge > ans1 {
			ans1 = merge
		}

		var q int
		fmt.Fscan(in, &q)
		for i := 0; i < q; i++ {
			var w int64
			fmt.Fscan(in, &w)
			if w == 0 {
				fmt.Fprintln(out, ans0)
			} else {
				fmt.Fprintln(out, ans1)
			}
		}
	}
}
