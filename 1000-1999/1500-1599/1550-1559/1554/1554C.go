package main

import (
	"bufio"
	"fmt"
	"os"
)

func mexXor(n, m int) int {
	if n > m {
		return 0
	}
	const inf int = 1<<31 - 1
	dp := [2]int{0, inf}
	for i := 30; i >= 0; i-- {
		ndp := [2]int{inf, inf}
		bitN := (n >> i) & 1
		bitM := (m >> i) & 1
		for gt := 0; gt < 2; gt++ {
			val := dp[gt]
			if val >= inf {
				continue
			}
			for b := 0; b < 2; b++ {
				var newGt int
				if gt == 0 {
					if b < bitM {
						continue
					}
					if b > bitM {
						newGt = 1
					} else {
						newGt = 0
					}
				} else {
					newGt = 1
				}
				newVal := val | ((b ^ bitN) << i)
				if newVal < ndp[newGt] {
					ndp[newGt] = newVal
				}
			}
		}
		dp = ndp
	}
	return dp[1]
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		ans := mexXor(n, m)
		fmt.Fprintln(writer, ans)
	}
}
