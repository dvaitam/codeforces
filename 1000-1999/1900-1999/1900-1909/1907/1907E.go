package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	// Precompute counts of digit triples without carry
	var triple [10]int64
	for a := 0; a <= 9; a++ {
		for b := 0; b <= 9; b++ {
			for c := 0; c <= 9; c++ {
				s := a + b + c
				if s <= 9 {
					triple[s]++
				}
			}
		}
	}

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		ans := int64(1)
		if n == 0 {
			ans = triple[0]
		} else {
			for n > 0 {
				d := n % 10
				ans *= triple[d]
				n /= 10
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
