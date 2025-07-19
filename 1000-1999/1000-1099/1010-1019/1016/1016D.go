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

	for {
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			break
		}
		a := make([]int, n)
		b := make([]int, m)
		tot := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			tot ^= a[i]
		}
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &b[j])
			tot ^= b[j]
		}
		if tot != 0 {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
			c := make([]int, m)
			ans := make([][]int, n)
			for i := 0; i < n; i++ {
				ans[i] = make([]int, m)
			}
			for i := 0; i < n-1; i++ {
				rowXor := 0
				for j := 0; j < m-1; j++ {
					ans[i][j] = 1
					rowXor ^= 1
					c[j] ^= 1
				}
				ans[i][m-1] = rowXor ^ a[i]
				c[m-1] ^= ans[i][m-1]
			}
			for j := 0; j < m; j++ {
				ans[n-1][j] = b[j] ^ c[j]
			}
			for i := 0; i < n; i++ {
				for j := 0; j < m; j++ {
					if j > 0 {
						writer.WriteByte(' ')
					}
					fmt.Fprint(writer, ans[i][j])
				}
				writer.WriteByte('\n')
			}
		}
	}
}
