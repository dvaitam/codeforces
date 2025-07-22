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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	v := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &v[i])
	}
	t := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &t[i])
	}

	prefix := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + t[i]
	}

	cntDiff := make([]int64, n+2)
	add := make([]int64, n+2)

	for i := 1; i <= n; i++ {
		l := i
		r := n
		pos := n + 1
		for l <= r {
			m := (l + r) / 2
			if prefix[m]-prefix[i-1] >= v[i] {
				pos = m
				r = m - 1
			} else {
				l = m + 1
			}
		}
		cntDiff[i]++
		if pos <= n {
			cntDiff[pos]--
			add[pos] += v[i] - (prefix[pos-1] - prefix[i-1])
		} else {
			cntDiff[n+1]--
		}
	}

	cur := int64(0)
	for i := 1; i <= n; i++ {
		cur += cntDiff[i]
		melted := cur*t[i] + add[i]
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, melted)
	}
	writer.WriteByte('\n')
}
