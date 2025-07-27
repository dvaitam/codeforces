package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	pos := make([]int, 0, k)
	var sum int64
	for i := 1; i <= n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x > n-k {
			sum += int64(x)
			pos = append(pos, i)
		}
	}

	var ways int64 = 1
	for i := 1; i < len(pos); i++ {
		diff := int64(pos[i] - pos[i-1])
		ways = (ways * diff) % mod
	}

	fmt.Fprintf(writer, "%d %d\n", sum, ways)
}
