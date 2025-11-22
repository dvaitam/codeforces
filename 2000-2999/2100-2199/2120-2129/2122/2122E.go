package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

type option struct {
	delta int
	cnt   int64
}

func buildOptions(top, bottom []int, k int) [][]option {
	n := len(top)
	opts := make([][]option, n-1)
	for i := 0; i < n-1; i++ {
		a := top[i+1]
		b := bottom[i]
		list := make([]option, 0)
		for d := -(k - 1); d <= k-1; d++ {
			var ways int64
			switch {
			case a != -1 && b != -1:
				if a-b == d {
					ways = 1
				}
			case a != -1:
				need := a - d
				if 1 <= need && need <= k {
					ways = 1
				}
			case b != -1:
				need := d + b
				if 1 <= need && need <= k {
					ways = 1
				}
			default:
				// both unknown
				diff := d
				if diff < 0 {
					diff = -diff
				}
				if diff < k {
					ways = int64(k - diff)
				}
			}
			if ways > 0 {
				list = append(list, option{delta: d, cnt: ways})
			}
		}
		opts[i] = list
	}
	return opts
}

func solveCase(n, k int, top, bottom []int) int64 {
	// Special case n == 1: no constraints.
	if n == 1 {
		ans := int64(1)
		if top[0] == -1 {
			ans = (ans * int64(k)) % mod
		}
		if bottom[0] == -1 {
			ans = (ans * int64(k)) % mod
		}
		return ans
	}

	opts := buildOptions(top, bottom, k)

	dp0 := map[int]int64{0: 1}      // uncapped states: pref -> ways
	dp1 := make([]map[int]int64, k) // capped states by slack x (0..k-1)

	for _, arr := range opts {
		ndp0 := make(map[int]int64)
		ndp1 := make([]map[int]int64, k)

		// Transition from uncapped states.
		for pref, waysState := range dp0 {
			for _, op := range arr {
				val := (waysState * op.cnt) % mod
				pref2 := pref + op.delta
				if op.delta < 0 {
					x := -op.delta
					if ndp1[x] == nil {
						ndp1[x] = make(map[int]int64)
					}
					ndp1[x][pref2] = (ndp1[x][pref2] + val) % mod
				} else {
					ndp0[pref2] = (ndp0[pref2] + val) % mod
				}
			}
		}

		// Transition from capped states.
		for x, mp := range dp1 {
			if len(mp) == 0 {
				continue
			}
			for pref, waysState := range mp {
				for _, op := range arr {
					val := (waysState * op.cnt) % mod
					if op.delta >= 0 {
						if op.delta <= x {
							pref2 := pref + op.delta
							x2 := x - op.delta
							if ndp1[x2] == nil {
								ndp1[x2] = make(map[int]int64)
							}
							ndp1[x2][pref2] = (ndp1[x2][pref2] + val) % mod
						}
					} else {
						pref2 := pref + op.delta
						x2 := -op.delta
						if ndp1[x2] == nil {
							ndp1[x2] = make(map[int]int64)
						}
						ndp1[x2][pref2] = (ndp1[x2][pref2] + val) % mod
					}
				}
			}
		}

		dp0 = ndp0
		dp1 = ndp1
	}

	var ans int64
	for _, v := range dp0 {
		ans = (ans + v) % mod
	}
	for _, mp := range dp1 {
		for _, v := range mp {
			ans = (ans + v) % mod
		}
	}

	// Multiply by free choices of unused cells.
	if top[0] == -1 {
		ans = (ans * int64(k)) % mod
	}
	if bottom[n-1] == -1 {
		ans = (ans * int64(k)) % mod
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		top := make([]int, n)
		bottom := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &top[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &bottom[i])
		}
		fmt.Fprintln(out, solveCase(n, k, top, bottom))
	}
}
