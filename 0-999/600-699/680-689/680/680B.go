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

	var n, a int
	if _, err := fmt.Fscan(reader, &n, &a); err != nil {
		return
	}
	criminals := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &criminals[i])
	}

	ans := 0
	maxDist := a - 1
	if n-a > maxDist {
		maxDist = n - a
	}
	for d := 0; d <= maxDist; d++ {
		l := a - d
		r := a + d
		if l >= 1 && r <= n {
			if l == r {
				ans += criminals[l]
			} else if criminals[l] == 1 && criminals[r] == 1 {
				ans += 2
			}
		} else if l >= 1 && l <= n && (r < 1 || r > n) {
			ans += criminals[l]
		} else if r >= 1 && r <= n && (l < 1 || l > n) {
			ans += criminals[r]
		}
	}

	fmt.Fprintln(writer, ans)
}
