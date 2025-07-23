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

	var n, d int
	if _, err := fmt.Fscan(in, &n, &d); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	const inf = int(1e9)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = inf
	}
	dist[0] = 0

	for i := 0; i < n; i++ {
		if dist[i] == inf {
			continue
		}
		for j := i + 1; j <= i+d && j < n; j++ {
			if s[j] == '1' && dist[i]+1 < dist[j] {
				dist[j] = dist[i] + 1
			}
		}
	}

	if dist[n-1] == inf {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, dist[n-1])
	}
}
