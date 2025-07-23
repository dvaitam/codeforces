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
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	t := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &t[i])
	}

	base := 0
	bonus := make([]int, n)
	for i := 0; i < n; i++ {
		if t[i] == 1 {
			base += arr[i]
		} else {
			bonus[i] = arr[i]
		}
	}

	cur := 0
	for i := 0; i < k && i < n; i++ {
		cur += bonus[i]
	}
	best := cur
	for i := k; i < n; i++ {
		cur += bonus[i] - bonus[i-k]
		if cur > best {
			best = cur
		}
	}

	fmt.Fprintln(writer, base+best)
}
