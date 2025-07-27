package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = int64(1e9 + 7)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m)

	max := n
	if m > max {
		max = m
	}
	max += 2

	fib := make([]int64, max+1)
	fib[1] = 1
	for i := 2; i <= max; i++ {
		fib[i] = (fib[i-1] + fib[i-2]) % mod
	}

	ans := (fib[n+1] + fib[m+1] - 1) % mod
	if ans < 0 {
		ans += mod
	}
	ans = (ans * 2) % mod
	fmt.Println(ans)
}