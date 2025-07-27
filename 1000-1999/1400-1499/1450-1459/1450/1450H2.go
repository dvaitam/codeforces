package main

import (
	"bufio"
	"fmt"
	"os"
)

func minIntersections(coloring []byte) int {
	blackEven, blackOdd := 0, 0
	for i, c := range coloring {
		if c == 'b' {
			if (i+1)%2 == 0 {
				blackEven++
			} else {
				blackOdd++
			}
		}
	}
	return abs(blackEven-blackOdd) / 2
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func expectedValue(s string) int {
	mod := int64(998244353)
	arr := []byte(s)
	var pos []int
	for i, c := range arr {
		if c == '?' {
			pos = append(pos, i)
		}
	}
	total := int64(0)
	cnt := int64(0)
	var dfs func(int)
	dfs = func(idx int) {
		if idx == len(pos) {
			if countBlack(arr)%2 == 0 {
				cnt++
				total += int64(minIntersections(arr))
			}
			return
		}
		arr[pos[idx]] = 'b'
		dfs(idx + 1)
		arr[pos[idx]] = 'w'
		dfs(idx + 1)
	}
	dfs(0)
	if cnt == 0 {
		return 0
	}
	// compute total * inv(cnt) mod
	inv := modPow(cnt%mod, mod-2, mod)
	return int(total % mod * inv % mod)
}

func countBlack(arr []byte) int {
	c := 0
	for _, ch := range arr {
		if ch == 'b' {
			c++
		}
	}
	return c
}

func modPow(a, b, mod int64) int64 {
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

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	_ = n
	var s string
	fmt.Fscan(in, &s)
	fmt.Fprintln(out, expectedValue(s))
	bs := []byte(s)
	for i := 0; i < m; i++ {
		var idx int
		var chStr string
		fmt.Fscan(in, &idx, &chStr)
		idx--
		bs[idx] = chStr[0]
		fmt.Fprintln(out, expectedValue(string(bs)))
	}
}
