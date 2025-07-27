package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	ns := make([]int, t)
	maxN := 0
	for i := 0; i < t; i++ {
		fmt.Fscan(reader, &ns[i])
		if ns[i] > maxN {
			maxN = ns[i]
		}
	}

	// Precompute factorials up to 2*maxN
	fac := make([]int64, 2*maxN+1)
	fac[0] = 1
	for i := 1; i <= 2*maxN; i++ {
		fac[i] = fac[i-1] * int64(i) % MOD
	}

	inv2 := int64((MOD + 1) / 2) // modular inverse of 2 modulo MOD
	for _, n := range ns {
		ans := fac[2*n] * inv2 % MOD
		fmt.Fprintln(writer, ans)
	}
}
