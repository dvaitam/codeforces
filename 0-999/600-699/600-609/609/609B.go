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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	cnt := make([]int, m+1)
	for i := 0; i < n; i++ {
		var g int
		fmt.Fscan(reader, &g)
		cnt[g]++
	}
	rem := n
	var ans int64
	for i := 1; i <= m; i++ {
		rem -= cnt[i]
		ans += int64(cnt[i] * rem)
	}
	fmt.Fprintln(writer, ans)
}
