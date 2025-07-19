package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	conn := make([][]int, n)
	degree := make([]int, n)
	for i := 1; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		conn[a] = append(conn[a], b)
		conn[b] = append(conn[b], a)
		degree[a]++
		degree[b]++
	}
	order := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if degree[i] <= 1 {
			order = append(order, i)
		}
	}
	// prune leaves
	for i := 0; i < len(order); i++ {
		u := order[i]
		degree[u]--
		for _, v := range conn[u] {
			degree[v]--
			if degree[v] == 1 {
				order = append(order, v)
			}
		}
	}
	answer := make([]int, n)
	for i := range answer {
		answer[i] = -1
	}
	maxval := 0
	for _, u := range order {
		count := 0
		maximum := 0
		for _, v := range conn[u] {
			if answer[v] < 0 {
				continue
			}
			if answer[v] > maximum {
				maximum = answer[v]
			}
			count++
		}
		var val int
		if count == 0 {
			val = 1
		} else if count == 1 {
			val = maximum + 1
		} else {
			// next power of two greater than maximum
			// find highest bit of maximum
			hb := 1 << (31 - bits.LeadingZeros32(uint32(maximum)))
			val = hb << 1
		}
		answer[u] = val
		if val > maxval {
			maxval = val
		}
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	if maxval >= (1 << 26) {
		fmt.Fprintln(out, "Impossible!")
	} else {
		for i := 0; i < n; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			// compute character: 'Z' - ctz(answer)
			tz := bits.TrailingZeros32(uint32(answer[i]))
			out.WriteByte(byte('Z' - tz))
		}
		out.WriteByte('\n')
	}
}
