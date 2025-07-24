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
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		pos := make([]int, 1001)
		for i := 1; i <= n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			pos[x] = i
		}
		ans := -1
		for i := 1; i <= 1000; i++ {
			if pos[i] == 0 {
				continue
			}
			for j := 1; j <= 1000; j++ {
				if pos[j] == 0 {
					continue
				}
				if gcd(i, j) == 1 {
					sum := pos[i] + pos[j]
					if sum > ans {
						ans = sum
					}
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
