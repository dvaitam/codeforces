package main

import (
	"bufio"
	"fmt"
	"os"
)

type Key struct {
	v    int
	mask uint64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	blue := make([]int, n+1)
	red := make([]int, n+1)
	for i := 1; i <= n-1; i++ {
		fmt.Fscan(reader, &blue[i], &red[i])
	}

	var q int
	fmt.Fscan(reader, &q)

	queries := make(map[Key][]int)
	answers := make([]int64, q)

	for i := 0; i < q; i++ {
		var v int
		var s string
		fmt.Fscan(reader, &v, &s)
		var mask uint64
		for j := 0; j < len(s); j++ {
			if s[j] == 'R' || s[j] == 'r' {
				mask |= 1 << uint(j)
			}
		}
		k := Key{v, mask}
		queries[k] = append(queries[k], i)
		answers[i] = -1
	}

	v := 1
	var mask uint64
	var t int64

	check := func(v int, mask uint64, t int64) {
		k := Key{v, mask}
		if idxs, ok := queries[k]; ok {
			for _, idx := range idxs {
				if answers[idx] == -1 {
					answers[idx] = t
				}
			}
		}
	}

	check(v, mask, t)

	const maxSteps = int64(1e7)
	for v != n && t < maxSteps {
		bit := (mask >> uint(v-1)) & 1
		mask ^= 1 << uint(v-1)
		if bit == 0 {
			v = red[v]
		} else {
			v = blue[v]
		}
		t++
		check(v, mask, t)
	}

	for i := 0; i < q; i++ {
		fmt.Fprintln(writer, answers[i])
	}
}
