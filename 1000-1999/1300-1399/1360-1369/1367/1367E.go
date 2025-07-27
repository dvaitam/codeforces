package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)
		freq := [26]int{}
		for _, ch := range s {
			freq[ch-'a']++
		}
		best := 0
		for length := 1; length <= n; length++ {
			g := gcd(length, k)
			cycleLen := length / g
			cycles := 0
			for _, cnt := range freq {
				cycles += cnt / cycleLen
			}
			if cycles >= g && length > best {
				best = length
			}
		}
		fmt.Fprintln(out, best)
	}
}
