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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		ans := 0
		for i := 1; i <= n && i <= m; i++ {
			for j := 1; i+j < m && i*j < n; j++ {
				a := (n - i*j) / (i + j)
				b := m - i - j
				if a < b {
					ans += a
				} else {
					ans += b
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
