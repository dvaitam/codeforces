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
	var n int
	var mod int64
	if _, err := fmt.Fscan(in, &n, &mod); err != nil {
		return
	}
	maxR := (n + 1) / 2
	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}
	fac := make([]int64, maxR+1)
	fac[0] = 1
	for i := 1; i <= maxR; i++ {
		fac[i] = (fac[i-1] * int64(i)) % mod
	}
	S := make([][]int64, n+1)
	for i := range S {
		S[i] = make([]int64, maxR+1)
	}
	S[0][0] = 1
	for i := 1; i <= n; i++ {
		for j := 1; j <= maxR && j <= i; j++ {
			S[i][j] = (S[i-1][j-1] + int64(j)*S[i-1][j]) % mod
		}
	}
	var ans int64
	for r := 1; r <= maxR; r++ {
		K := n - r + 1
		if r > K {
			break
		}
		term := pow2[K-r]
		term = (term * fac[r]) % mod
		term = (term * S[K][r]) % mod
		ans = (ans + term) % mod
	}
	fmt.Fprintln(out, ans)
}
