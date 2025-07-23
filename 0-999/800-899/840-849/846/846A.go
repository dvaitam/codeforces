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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	prefixZero := make([]int, n+1)
	prefixOne := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefixZero[i] = prefixZero[i-1]
		prefixOne[i] = prefixOne[i-1]
		if a[i-1] == 0 {
			prefixZero[i]++
		} else {
			prefixOne[i]++
		}
	}

	ans := 0
	for k := 0; k <= n; k++ {
		zeros := prefixZero[k]
		ones := prefixOne[n] - prefixOne[k]
		if zeros+ones > ans {
			ans = zeros + ones
		}
	}
	fmt.Fprintln(writer, ans)
}
