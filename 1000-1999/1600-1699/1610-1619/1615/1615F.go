package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	inv2 := (MOD + 1) / 2
	for ; T > 0; T-- {
		var n int
		var s, t string
		fmt.Fscan(in, &n)
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)

		size := 4*n + 1
		offset := 2 * n
		f := make([]int64, size)
		g := make([]int64, size)
		f[offset] = 1

		for i := 0; i < n; i++ {
			nf := make([]int64, size)
			ng := make([]int64, size)

			for x := 0; x <= 1; x++ {
				if s[i] != '?' && s[i] != byte('0'+x) {
					continue
				}
				for y := 0; y <= 1; y++ {
					if t[i] != '?' && t[i] != byte('0'+y) {
						continue
					}
					sign := 1
					if i%2 == 1 {
						sign = -1
					}
					vx := -1
					if x == 1 {
						vx = 1
					}
					vy := -1
					if y == 1 {
						vy = 1
					}
					delta := sign * (vx - vy)

					start := offset - 2*i
					end := offset + 2*i
					for j := start; j <= end; j++ {
						nj := j + delta
						nf[nj] = (nf[nj] + f[j]) % MOD
						ng[nj] = (ng[nj] + g[j]) % MOD
					}
				}
			}

			for j := 0; j < size; j++ {
				if nf[j] != 0 {
					val := int64(abs(j - offset))
					ng[j] = (ng[j] + nf[j]*val) % MOD
				}
			}
			f = nf
			g = ng
		}

		ans := g[offset] * inv2 % MOD
		fmt.Fprintln(out, ans)
	}
}
