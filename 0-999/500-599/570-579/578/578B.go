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

	var n, k int
	var x int64
	if _, err := fmt.Fscan(reader, &n, &k, &x); err != nil {
		return
	}
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] | arr[i]
	}
	suffix := make([]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		suffix[i] = suffix[i+1] | arr[i]
	}

	pow := int64(1)
	for i := 0; i < k; i++ {
		pow *= x
	}

	var ans int64
	for i := 0; i < n; i++ {
		val := prefix[i] | (arr[i] * pow) | suffix[i+1]
		if val > ans {
			ans = val
		}
	}

	fmt.Fprintln(writer, ans)
}
