package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type state struct {
	pos int
	str string
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k, m int
		fmt.Fscan(reader, &n, &k, &m)
		var s string
		fmt.Fscan(reader, &s)
		next := make([][]int, m+1)
		for i := range next {
			next[i] = make([]int, k)
			for j := range next[i] {
				next[i][j] = m
			}
		}
		for j := 0; j < k; j++ {
			next[m][j] = m
		}
		for i := m - 1; i >= 0; i-- {
			copy(next[i], next[i+1])
			c := int(s[i] - 'a')
			if c < k {
				next[i][c] = i
			}
		}
		q := []state{{0, ""}}
		visited := make([][]bool, n+1)
		for i := range visited {
			visited[i] = make([]bool, m+1)
		}
		visited[0][0] = true
		found := false
		ans := ""
		for len(q) > 0 && !found {
			cur := q[0]
			q = q[1:]
			if len(cur.str) == n {
				continue
			}
			for c := 0; c < k; c++ {
				nextPos := next[cur.pos][c]
				newStr := cur.str + string('a'+byte(c))
				if nextPos == m {
					ans = newStr + strings.Repeat("a", n-len(newStr))
					found = true
					break
				}
				newPos := nextPos + 1
				l := len(cur.str) + 1
				if !visited[l][newPos] {
					visited[l][newPos] = true
					q = append(q, state{newPos, newStr})
				}
			}
		}
		if found {
			fmt.Fprintln(writer, "NO")
			fmt.Fprintln(writer, ans)
		} else {
			fmt.Fprintln(writer, "YES")
		}
	}
}
