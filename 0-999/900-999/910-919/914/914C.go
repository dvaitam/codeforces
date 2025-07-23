package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func popcount(x int) int {
	cnt := 0
	for x > 0 {
		cnt += x & 1
		x >>= 1
	}
	return cnt
}

var steps [1005]int

func calc(x int) int {
	if steps[x] != -1 {
		return steps[x]
	}
	if x == 1 {
		steps[x] = 0
		return 0
	}
	steps[x] = calc(popcount(x)) + 1
	return steps[x]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n string
	var k int
	fmt.Fscan(in, &n)
	fmt.Fscan(in, &k)

	if k == 0 {
		fmt.Fprintln(out, 1)
		return
	}

	for i := 0; i < len(steps); i++ {
		steps[i] = -1
	}
	steps[1] = 0
	for i := 2; i < len(steps); i++ {
		calc(i)
	}

	L := len(n)
	comb := make([][]int64, L+1)
	for i := 0; i <= L; i++ {
		comb[i] = make([]int64, L+1)
	}
	for i := 0; i <= L; i++ {
		comb[i][0] = 1
		comb[i][i] = 1
		for j := 1; j < i; j++ {
			comb[i][j] = (comb[i-1][j-1] + comb[i-1][j]) % mod
		}
	}

	count := func(r int) int64 {
		if r < 0 {
			return 0
		}
		ones := 0
		var ans int64
		for i := 0; i < L; i++ {
			if n[i] == '1' {
				rem := L - i - 1
				need := r - ones
				if need >= 0 && need <= rem {
					ans = (ans + comb[rem][need]) % mod
				}
				ones++
			}
		}
		if ones == r {
			ans = (ans + 1) % mod
		}
		return ans
	}

	var ans int64
	for r := 1; r <= 1000; r++ {
		if steps[r] == k-1 {
			ans = (ans + count(r)) % mod
		}
	}
	if k == 1 {
		ans = (ans - 1 + mod) % mod
	}
	fmt.Fprintln(out, ans)
}
