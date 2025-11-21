package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		pCount := 0
		for i := 0; i < n; i++ {
			if s[i] == '1' {
				pCount++
			}
		}

		countZeros := func() int {
			row := make([]byte, n)
			copy(row, s)
			zeros := 0
			for j := n - 1; j >= 0; j-- {
				row[j] = flip(row[j])
				for k := 0; k < n; k++ {
					if row[k] == '0' {
						zeros++
					}
				}
				row[j] = flip(row[j])
			}
			return zeros
		}

		if pCount <= n/2 {
			fmt.Fprintln(out, countZeros())
		} else {
			fmt.Fprintln(out, n)
		}
	}
}

func flip(b byte) byte {
	if b == '0' {
		return '1'
	}
	return '0'
}
