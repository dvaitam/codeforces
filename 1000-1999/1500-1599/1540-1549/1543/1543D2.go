package main

import (
	"bufio"
	"fmt"
	"os"
)

// xorK performs digit-wise addition modulo k (k-itwise XOR).
func xorK(a, b, k int) int {
	res := 0
	mul := 1
	for a > 0 || b > 0 {
		res += ((a%k + b%k) % k) * mul
		a /= k
		b /= k
		mul *= k
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		prev := 0
		for i := 0; i < n; i++ {
			q := xorK(prev, i, k)
			fmt.Fprintln(writer, q)
			writer.Flush()
			var ans int
			if _, err := fmt.Fscan(reader, &ans); err != nil {
				return
			}
			if ans == 1 {
				break
			}
			prev = xorK(prev, q, k)
		}
	}
}
