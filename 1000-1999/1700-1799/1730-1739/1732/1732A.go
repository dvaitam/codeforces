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
		arr := make([]int, n)
		g := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			g = gcd(g, arr[i])
		}
		if g == 1 {
			fmt.Fprintln(writer, 0)
			continue
		}
		ans := 3
		if gcd(g, n) == 1 {
			ans = 1
		} else if gcd(g, n-1) == 1 {
			ans = 2
		}
		fmt.Fprintln(writer, ans)
	}
}
