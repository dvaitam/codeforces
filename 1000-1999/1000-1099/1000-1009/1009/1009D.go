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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	if m < n-1 {
		fmt.Fprint(writer, "Impossible")
		return
	}
	if n < 1000 {
		t := 0
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				if gcd(i, j) == 1 {
					t++
				}
			}
		}
		if t < m {
			fmt.Fprint(writer, "Impossible")
			return
		}
	}
	fmt.Fprintln(writer, "Possible")
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if gcd(i, j) == 1 {
				fmt.Fprintf(writer, "%d %d\n", i, j)
				m--
				if m == 0 {
					return
				}
			}
		}
	}
}
