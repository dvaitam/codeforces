package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func extgcd(a, b int64) (g, x, y int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := extgcd(b, a%b)
	return g, y1, x1 - (a/b)*y1
}

func modInverse(a, mod int64) int64 {
	g, x, _ := extgcd(a, mod)
	if g != 1 {
		return 0
	}
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	diffs := make([]int64, n)
	var bsum int64
	for i := 0; i < n; i++ {
		var a, b int64
		fmt.Fscan(reader, &a, &b)
		diffs[i] = a - b
		bsum += b
	}
	sort.Slice(diffs, func(i, j int) bool { return diffs[i] > diffs[j] })
	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + diffs[i-1]
	}
	values := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		values[i] = bsum + prefix[i]
	}

	S := int64(1)
	for S*S <= int64(n) {
		S++
	}
	pre := make([][]int64, S+1)
	for s := int64(1); s <= S; s++ {
		pre[s] = make([]int64, s)
		for i := 0; i < int(s); i++ {
			pre[s][i] = -1 << 63
		}
		for r := int64(0); r < s; r++ {
			maxv := int64(-1 << 63)
			for x := r; x <= int64(n); x += s {
				v := values[x]
				if v > maxv {
					maxv = v
				}
			}
			pre[s][r] = maxv
		}
	}

	var m int
	fmt.Fscan(reader, &m)
	for ; m > 0; m-- {
		var xj, yj int64
		fmt.Fscan(reader, &xj, &yj)
		g := gcd(xj, yj)
		if int64(n)%g != 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		lcm := xj / g * yj
		x1 := xj / g
		y1 := yj / g
		n1 := int64(n) / g
		inv := modInverse(x1%y1, y1)
		t0 := (n1 % y1) * inv % y1
		r0 := xj * t0
		if r0 > int64(n) {
			fmt.Fprintln(writer, -1)
			continue
		}
		if lcm <= S {
			res := pre[lcm][r0%lcm]
			if res == -1<<63 {
				fmt.Fprintln(writer, -1)
			} else {
				fmt.Fprintln(writer, res)
			}
		} else {
			maxv := int64(-1 << 63)
			for R := r0; R <= int64(n); R += lcm {
				if v := values[R]; v > maxv {
					maxv = v
				}
			}
			if maxv == -1<<63 {
				fmt.Fprintln(writer, -1)
			} else {
				fmt.Fprintln(writer, maxv)
			}
		}
	}
}
