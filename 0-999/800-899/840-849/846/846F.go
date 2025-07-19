package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 1 << 20

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	p := make([]int, N)
	var s, ans int64

	for i := 1; i <= n; i++ {
		var ai int
		fmt.Fscan(reader, &ai)
		s += int64(i - p[ai])
		p[ai] = i
		ans += s
	}
	ans = ans + ans - int64(n)
	fans := float64(ans) / float64(n) / float64(n)
	fmt.Fprintf(writer, "%.12f\n", fans)
}
