package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	const inf = int(1e9)
	dist := make([]int, n)
	last := -inf
	for i := 0; i < n; i++ {
		if arr[i] == 0 {
			last = i
		}
		dist[i] = i - last
	}
	last = inf
	for i := n - 1; i >= 0; i-- {
		if arr[i] == 0 {
			last = i
		}
		if d := last - i; d < dist[i] {
			dist[i] = d
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, dist[i])
	}
	fmt.Fprintln(out)
}
