package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 2520

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	prefix := make([]int, mod+1)
	for i := 1; i <= mod; i++ {
		prefix[i] = prefix[i-1]
		if gcd(i, mod) == 1 {
			prefix[i]++
		}
	}

	blocks := n / mod
	rem := int(n % mod)
	ans := int64(prefix[mod])*blocks + int64(prefix[rem])
	fmt.Println(ans)
}
