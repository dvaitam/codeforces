package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	fmt.Fscan(in, &n, &m)
	if n == 0 || m == 0 {
		fmt.Fprintln(out, 0)
		return
	}
	maxR := 0
	for p := 1; (1 << p) <= n; p++ {
		maxR = p
	}
	M := maxR * m
	minD := make([]byte, M+1)
	for d := 1; d <= maxR; d++ {
		for j := 1; j <= m; j++ {
			x := d * j
			if x > M {
				break
			}
			if minD[x] == 0 {
				minD[x] = byte(d)
			}
		}
	}
	counts := make([]int64, maxR+1)
	cntPerD := make([]int64, maxR+1)
	idxR := 1
	if maxR == 0 {
		idxR = 0
	}
	for x := 1; x <= M; x++ {
		d := int(minD[x])
		if d > 0 {
			cntPerD[d]++
		}
		for idxR > 0 && idxR <= maxR && x == idxR*m {
			var total int64
			for i := 1; i <= idxR; i++ {
				total += cntPerD[i]
			}
			counts[idxR] = total
			idxR++
		}
	}
	// prefix for remaining Rs if any (when m==0 maybe?).
	for idxR <= maxR {
		var total int64
		for i := 1; i <= idxR; i++ {
			total += cntPerD[i]
		}
		counts[idxR] = total
		idxR++
	}

	isPP := make([]bool, n+1)
	for a := 2; a*a <= n; a++ {
		v := a * a
		for v <= n {
			if !isPP[v] {
				isPP[v] = true
			}
			if v > n/a {
				break
			}
			v *= a
		}
	}
	baseCount := make([]int64, maxR+1)
	for b := 2; b <= n; b++ {
		if isPP[b] {
			continue
		}
		cur := b
		r := 0
		for cur <= n {
			r++
			if cur > n/b {
				break
			}
			cur *= b
		}
		if r > maxR {
			r = maxR
		}
		baseCount[r]++
	}

	var ans int64 = 1 // for number 1
	for r := 1; r <= maxR; r++ {
		ans += counts[r] * baseCount[r]
	}
	fmt.Fprintln(out, ans)
}
