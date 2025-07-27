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
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		dCount := 0
		kCount := 0
		mp := make(map[[2]int]int)
		for i := 0; i < n; i++ {
			if s[i] == 'D' {
				dCount++
			} else {
				kCount++
			}
			g := gcd(dCount, kCount)
			key := [2]int{dCount / g, kCount / g}
			mp[key]++
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, mp[key])
		}
		writer.WriteByte('\n')
	}
}
