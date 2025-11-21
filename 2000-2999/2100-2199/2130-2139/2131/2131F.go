package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1e9 + 7

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var aStr, bStr string
		fmt.Fscan(in, &aStr)
		fmt.Fscan(in, &bStr)

		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			if aStr[i] != bStr[i] {
				prefix[i+1] = prefix[i] + 1
			} else {
				prefix[i+1] = prefix[i]
			}
		}

		answer := int64(0)
		for i := 1; i <= n; i++ {
			answer += prefix[i]
		}
		fmt.Fprintln(out, answer)
	}
}

