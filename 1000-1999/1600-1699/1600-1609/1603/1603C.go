package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

type state struct {
	val  int
	cnt  int64
	cost int64
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		fmt.Fprintln(out, solve(a))
	}
}

func solve(a []int) int64 {
	var ans int64
	vec := []state{}
	for i := len(a) - 1; i >= 0; i-- {
		x := a[i]
		mp := make(map[int]state)
		mp[x] = state{val: x, cnt: 1, cost: 0}
		for _, st := range vec {
			k := (x + st.val - 1) / st.val
			newVal := x / k
			cnt := st.cnt
			cost := (st.cost + cnt*int64(k-1)) % mod
			if cur, ok := mp[newVal]; ok {
				cur.cnt = (cur.cnt + cnt) % mod
				cur.cost = (cur.cost + cost) % mod
				mp[newVal] = cur
			} else {
				mp[newVal] = state{val: newVal, cnt: cnt % mod, cost: cost}
			}
		}
		vec = vec[:0]
		for _, st := range mp {
			vec = append(vec, st)
			ans = (ans + st.cost) % mod
		}
	}
	return ans % mod
}
