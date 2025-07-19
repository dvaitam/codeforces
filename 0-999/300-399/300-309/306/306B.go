package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type interval struct {
	x, end, idx int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	p := make([]interval, m)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		p[i] = interval{x: x, end: x + y - 1, idx: i + 1}
	}
	sort.Slice(p, func(i, j int) bool {
		if p[i].x != p[j].x {
			return p[i].x < p[j].x
		}
		return p[i].end > p[j].end
	})

	vis := make([]bool, m+1)
	pre, now, cnt, id := -1, -1, m, 0
	for _, iv := range p {
		if iv.x <= pre+1 {
			if now < iv.end {
				now = iv.end
				id = iv.idx
			}
		} else {
			if now > pre {
				pre = now
				cnt--
				vis[id] = true
				now = iv.end
				id = iv.idx
			} else {
				pre = iv.end
				now = -1
				cnt--
				vis[iv.idx] = true
				id = 0
			}
		}
	}
	if now > pre {
		vis[id] = true
		cnt--
	}

	// output
	fmt.Fprintln(writer, cnt)
	first := true
	for i := 1; i <= m; i++ {
		if !vis[i] {
			if !first {
				writer.WriteByte(' ')
			}
			first = false
			fmt.Fprint(writer, i)
		}
	}
	fmt.Fprintln(writer)
}
