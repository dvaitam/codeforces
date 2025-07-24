package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	pairs := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pairs[i][0], &pairs[i][1])
	}

	fac := make([]int64, n+1)
	fac[0] = 1
	for i := 1; i <= n; i++ {
		fac[i] = fac[i-1] * int64(i) % mod
	}

	// count permutations sorted by first elements
	tmp := make([][2]int, n)
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i][0] < tmp[j][0] })
	cntA := int64(1)
	for i := 0; i < n; {
		j := i
		for j < n && tmp[j][0] == tmp[i][0] {
			j++
		}
		cntA = cntA * fac[j-i] % mod
		i = j
	}

	// count permutations sorted by second elements
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool { return tmp[i][1] < tmp[j][1] })
	cntB := int64(1)
	for i := 0; i < n; {
		j := i
		for j < n && tmp[j][1] == tmp[i][1] {
			j++
		}
		cntB = cntB * fac[j-i] % mod
		i = j
	}

	// count permutations sorted by both
	copy(tmp, pairs)
	sort.Slice(tmp, func(i, j int) bool {
		if tmp[i][0] == tmp[j][0] {
			return tmp[i][1] < tmp[j][1]
		}
		return tmp[i][0] < tmp[j][0]
	})
	valid := true
	for i := 1; i < n; i++ {
		if tmp[i][1] < tmp[i-1][1] {
			valid = false
			break
		}
	}
	cntAB := int64(0)
	if valid {
		cntAB = 1
		for i := 0; i < n; {
			j := i
			for j < n && tmp[j] == tmp[i] {
				j++
			}
			cntAB = cntAB * fac[j-i] % mod
			i = j
		}
	}

	ans := (fac[n] - cntA - cntB + cntAB) % mod
	if ans < 0 {
		ans += mod
	}
	fmt.Fprintln(out, ans)
}
