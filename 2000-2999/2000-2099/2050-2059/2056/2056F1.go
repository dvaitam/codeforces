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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for i := 0; i < t; i++ {
		var k_bits, m int
		fmt.Fscan(reader, &k_bits, &m)
		var s string
		fmt.Fscan(reader, &s)

		w := 0
		for j := 0; j < len(s); j++ {
			if s[j] == '1' {
				w++
			}
		}

		if w == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}

		S := make([][]int, w+1)
		for j := 0; j <= w; j++ {
			S[j] = make([]int, w+1)
		}
		S[0][0] = 1
		for i := 1; i <= w; i++ {
			for j := 1; j <= i; j++ {
				S[i][j] = ((j & 1) * S[i-1][j]) ^ S[i-1][j-1]
			}
		}

		ans := 0
		for x := 0; x < m; x++ {
			parity := 0
			for k := 1; k <= w; k++ {
				if S[w][k] == 1 {
					if (x & (k - 1)) == (k - 1) {
						parity ^= 1
					}
				}
			}
			if parity == 1 {
				ans ^= x
			}
		}
		fmt.Fprintln(writer, ans)
	}
}