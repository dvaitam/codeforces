package main

import (
	"bufio"
	"fmt"
	"os"
)

func minOps(s string) int {
	n := len(s)
	prefixOps := make([]int, n+1)
	prefixPar := make([]int, n+1)
	parity := 0
	op := 0
	for i := 1; i <= n; i++ {
		bit := int(s[i-1]-'0') ^ parity
		if bit == 1 {
			op++
			parity ^= 1
		}
		prefixOps[i] = op
		prefixPar[i] = parity
	}
	suffix0 := make([]int, n+2)
	suffix1 := make([]int, n+2)
	op = 0
	parity = 0
	for i := n; i >= 1; i-- {
		bit := int(s[i-1]-'0') ^ parity
		if bit == 0 {
			op++
			parity ^= 1
		}
		suffix0[i] = op
	}
	op = 0
	parity = 1
	for i := n; i >= 1; i-- {
		bit := int(s[i-1]-'0') ^ parity
		if bit == 0 {
			op++
			parity ^= 1
		}
		suffix1[i] = op
	}
	best := n + 5
	for k := 0; k <= n; k++ {
		p := prefixPar[k]
		ops := prefixOps[k]
		if p == 0 {
			ops += suffix0[k+1]
		} else {
			ops += suffix1[k+1]
		}
		if ops < best {
			best = ops
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		fmt.Fprintln(writer, minOps(s))
	}
}
