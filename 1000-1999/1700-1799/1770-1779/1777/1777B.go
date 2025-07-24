package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1e9 + 7

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

	fact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}

	for _, n := range ns {
		ans := fact[n] * int64(n) % mod * int64(n-1) % mod
		fmt.Fprintln(writer, ans)
	}
}
