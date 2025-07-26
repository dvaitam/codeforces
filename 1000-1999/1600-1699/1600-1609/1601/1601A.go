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
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	const maxBits = 30
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		counts := make([]int, maxBits)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			for b := 0; b < maxBits; b++ {
				if (x>>b)&1 == 1 {
					counts[b]++
				}
			}
		}
		g := 0
		for _, c := range counts {
			g = gcd(g, c)
		}
		if g == 0 {
			for k := 1; k <= n; k++ {
				fmt.Fprint(writer, k)
				if k == n {
					fmt.Fprintln(writer)
				} else {
					fmt.Fprint(writer, " ")
				}
			}
			continue
		}
		first := true
		for k := 1; k <= n; k++ {
			if g%k == 0 {
				if !first {
					fmt.Fprint(writer, " ")
				}
				fmt.Fprint(writer, k)
				first = false
			}
		}
		fmt.Fprintln(writer)
	}
}
