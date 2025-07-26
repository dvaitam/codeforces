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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		var m int
		fmt.Fscan(reader, &m)
		var l, r string
		fmt.Fscan(reader, &l)
		fmt.Fscan(reader, &r)

		n := len(s)
		nextPos := make([][10]int, n+1)
		for d := 0; d < 10; d++ {
			nextPos[n][d] = n
		}
		for i := n - 1; i >= 0; i-- {
			for d := 0; d < 10; d++ {
				nextPos[i][d] = nextPos[i+1][d]
			}
			digit := int(s[i] - '0')
			nextPos[i][digit] = i
		}

		cur := map[int]struct{}{-1: {}}
		found := false
		for i := 0; i < m && !found; i++ {
			nextSet := make(map[int]struct{})
			lo := int(l[i] - '0')
			hi := int(r[i] - '0')
			for pos := range cur {
				start := pos + 1
				if start > n {
					start = n
				}
				for d := lo; d <= hi; d++ {
					nxt := nextPos[start][d]
					if nxt == n {
						found = true
						break
					}
					nextSet[nxt] = struct{}{}
				}
				if found {
					break
				}
			}
			cur = nextSet
		}
		if found {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
