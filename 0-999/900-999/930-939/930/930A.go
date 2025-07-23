package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	depth := make([]int, n+1)
	count := make([]int, n+1)
	count[0] = 1 // node 1 at depth 0
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(reader, &p)
		depth[i] = depth[p] + 1
		d := depth[i]
		if d >= len(count) {
			temp := make([]int, d+1)
			copy(temp, count)
			count = temp
		}
		count[d]++
	}
	ans := 0
	for _, c := range count {
		if c%2 == 1 {
			ans++
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, ans)
	writer.Flush()
}
