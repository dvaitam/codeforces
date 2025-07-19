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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s, t string
	fmt.Fscan(reader, &s, &t)

	// convert chars to integer values 0..25
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = int(s[i] - 'a')
		b[i] = int(t[i] - 'a')
	}
	// ans has an extra position for carry
	ans := make([]int, n+1)

	// first pass: split sums and handle odd with base carry
	for i := n - 1; i >= 0; i-- {
		val := a[i] + b[i]
		if val%2 == 0 {
			ans[i] = val / 2
		} else {
			ans[i] = (val - 1) / 2
			ans[i+1] += 13
		}
	}
	// second pass: propagate base-26 carry from least significant
	carry := 0
	for i := n - 1; i >= 0; i-- {
		ans[i] += carry
		carry = ans[i] / 26
		ans[i] %= 26
	}

	// build result string
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = byte(ans[i]) + 'a'
	}
	writer.Write(res)
	writer.WriteByte('\n')
}
