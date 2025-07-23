package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pair struct {
	sign int
	val  int
}

func squareFree(x int) int {
	if x < 0 {
		x = -x
	}
	res := 1
	for p := 2; p*p <= x; p++ {
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt ^= 1
		}
		if cnt == 1 {
			res *= p
		}
	}
	if x > 1 {
		res *= x
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	keys := make([]Pair, n)
	for i, v := range arr {
		if v == 0 {
			keys[i] = Pair{0, 0}
			continue
		}
		sign := 1
		if v < 0 {
			sign = -1
		}
		sf := squareFree(v)
		keys[i] = Pair{sign, sf}
	}

	ans := make([]int, n+1)
	for l := 0; l < n; l++ {
		mp := make(map[Pair]int)
		groups := 0
		for r := l; r < n; r++ {
			k := keys[r]
			if k.sign != 0 {
				if mp[k] == 0 {
					groups++
				}
				mp[k]++
			}
			g := groups
			if g == 0 {
				g = 1
			}
			ans[g]++
		}
	}
	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	out.Flush()
}
