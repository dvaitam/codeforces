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

	var S, T string
	fmt.Fscan(in, &S)
	fmt.Fscan(in, &T)

	n := len(S)
	m := len(T)

	const INF = int(1e9)
	pos := make([]int, 26)
	for i := range pos {
		pos[i] = INF
	}
	for i := 1; i < n; i++ {
		idx := int(S[i] - 'a')
		if pos[idx] > i {
			pos[idx] = i
		}
	}

	bestLen := int(1 << 60)
	best := ""
	for s := 0; s <= m-2; s++ {
		ch := int(T[s] - 'a')
		idx := pos[ch]
		if idx == INF {
			continue
		}
		total := idx + (m - s)
		if total < bestLen {
			bestLen = total
			best = S[:idx] + T[s:]
		}
	}

	if bestLen == int(1<<60) {
		fmt.Fprintln(out, "-1")
	} else {
		fmt.Fprintln(out, best)
	}
}
