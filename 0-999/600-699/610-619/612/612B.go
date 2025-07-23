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
	pos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		var f int
		fmt.Fscan(reader, &f)
		pos[f] = i
	}

	var ans int64
	for i := 2; i <= n; i++ {
		diff := pos[i] - pos[i-1]
		if diff < 0 {
			diff = -diff
		}
		ans += int64(diff)
	}

	fmt.Fprintln(writer, ans)
}
