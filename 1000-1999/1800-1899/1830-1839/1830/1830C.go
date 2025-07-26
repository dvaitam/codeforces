package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353
const MAXN = 300000

var fact [MAXN + 1]int64
var inv [MAXN + 1]int64

func powmod(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func init() {
	fact[0] = 1
	for i := 1; i <= MAXN; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	inv[MAXN] = powmod(fact[MAXN], MOD-2)
	for i := MAXN; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % MOD
	}
}

func catalan(n int) int64 {
	if n < 0 {
		return 0
	}
	return fact[2*n] * inv[n] % MOD * inv[n+1] % MOD
}

type Interval struct{ l, r int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		intervals := make([]Interval, k)
		ok := true
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &intervals[i].l, &intervals[i].r)
			if (intervals[i].r-intervals[i].l+1)%2 == 1 {
				ok = false
			}
		}
		if n%2 == 1 {
			ok = false
		}
		if !ok {
			fmt.Fprintln(writer, 0)
			continue
		}
		intervals = append(intervals, Interval{1, n})
		sort.Slice(intervals, func(i, j int) bool {
			if intervals[i].l == intervals[j].l {
				return intervals[i].r > intervals[j].r
			}
			return intervals[i].l < intervals[j].l
		})

		children := make(map[Interval][]Interval)
		stack := []Interval{}
		valid := true
		for _, iv := range intervals {
			for len(stack) > 0 && !(stack[len(stack)-1].l <= iv.l && iv.r <= stack[len(stack)-1].r) {
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 && iv.l <= stack[len(stack)-1].r && iv.r > stack[len(stack)-1].r {
				valid = false
				break
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				children[parent] = append(children[parent], iv)
			}
			stack = append(stack, iv)
		}
		if !valid {
			fmt.Fprintln(writer, 0)
			continue
		}

		var dfs func(Interval) int64
		dfs = func(seg Interval) int64 {
			childs := children[seg]
			sort.Slice(childs, func(i, j int) bool { return childs[i].l < childs[j].l })
			used := 0
			ans := int64(1)
			for _, c := range childs {
				ans = ans * dfs(c) % MOD
				used += (c.r - c.l + 1) / 2
			}
			free := (seg.r-seg.l+1)/2 - used
			if free < 0 {
				return 0
			}
			ans = ans * catalan(free) % MOD
			return ans
		}

		result := dfs(Interval{1, n})
		fmt.Fprintln(writer, result)
	}
}
