package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func between(x, a, b, n int) bool {
	if a < b {
		return a < x && x < b
	}
	return a < x || x < b
}

func cross(a, b, c, d, n int) bool {
	return between(c, a, b, n) != between(d, a, b, n) && between(a, c, d, n) != between(b, c, d, n)
}

func pairings(idx []int) [][][2]int {
	if len(idx) == 0 {
		return [][][2]int{{}}
	}
	first := idx[0]
	rest := idx[1:]
	var res [][][2]int
	for i, v := range rest {
		rem := append([]int(nil), rest[:i]...)
		rem = append(rem, rest[i+1:]...)
		for _, sub := range pairings(rem) {
			res = append(res, append([][2]int{{first, v}}, sub...))
		}
	}
	return res
}

func minIntersections(coloring []byte) int {
	n := len(coloring)
	var black, white []int
	for i, c := range coloring {
		if c == 'b' {
			black = append(black, i)
		}
		if c == 'w' {
			white = append(white, i)
		}
	}
	pb := pairings(black)
	pw := pairings(white)
	best := int(^uint(0) >> 1)
	for _, bpair := range pb {
		for _, wpair := range pw {
			cnt := 0
			for _, bp := range bpair {
				for _, wp := range wpair {
					if cross(bp[0], bp[1], wp[0], wp[1], n) {
						cnt++
					}
				}
			}
			if cnt < best {
				best = cnt
			}
		}
	}
	if best == int(^uint(0)>>1) {
		return 0
	}
	return best
}

func expectedValue(s string) int {
	arr := []byte(s)
	var pos []int
	for i, c := range arr {
		if c == '?' {
			pos = append(pos, i)
		}
	}
	total := int64(0)
	count := 0
	var dfs func(int)
	dfs = func(i int) {
		if i == len(pos) {
			cb, cw := 0, 0
			for _, c := range arr {
				if c == 'b' {
					cb++
				} else if c == 'w' {
					cw++
				}
			}
			if cb%2 == 0 && cw%2 == 0 {
				count++
				total += int64(minIntersections(arr))
			}
			return
		}
		arr[pos[i]] = 'b'
		dfs(i + 1)
		arr[pos[i]] = 'w'
		dfs(i + 1)
	}
	dfs(0)
	if count == 0 {
		return 0
	}
	inv := modPow(int64(count), mod-2)
	return int(total % mod * inv % mod)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	var s string
	fmt.Fscan(in, &s)
	fmt.Fprintln(out, expectedValue(s))
}
