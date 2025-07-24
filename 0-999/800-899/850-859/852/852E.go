package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func modPow(base, exp int64) int64 {
	result := int64(1)
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % MOD
		}
		base = base * base % MOD
		exp >>= 1
	}
	return result
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	deg := make([]int, n+1)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		deg[a]++
		deg[b]++
	}

	leaves := 0
	for i := 1; i <= n; i++ {
		if deg[i] <= 1 {
			leaves++
		}
	}

	pow := modPow(2, int64(n-leaves))
	ans := int64(n+leaves) % MOD * pow % MOD
	fmt.Fprintln(writer, ans)
}
