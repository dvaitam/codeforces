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

	var h, w, q int
	if _, err := fmt.Fscan(reader, &h, &w, &q); err != nil {
		return
	}
	n := h * w
	perm := make([]bool, n)
	for i := 0; i < h; i++ {
		var line string
		fmt.Fscan(reader, &line)
		for j := 0; j < w; j++ {
			if line[j] == '#' {
				perm[i*w+j] = true
			}
		}
	}

	start := 0
	target := (h-1)*w + (w - 1)

	block := make([]int, n)
	used := make([]int, n)
	visit := make([]int, n)
	parent := make([]int, n)
	queue := make([]int, n)

	visitIter := 0
	usedIter := 0
	dayIter := 0

	bfs := func(avoidUsed int) bool {
		visitIter++
		head, tail := 0, 0
		queue[tail] = start
		tail++
		visit[start] = visitIter
		parent[start] = -1
		for head < tail {
			v := queue[head]
			head++
			if v == target {
				return true
			}
			r := v / w
			c := v % w
			// up
			if r > 0 {
				u := v - w
				if !perm[u] && block[u] != dayIter && (avoidUsed == 0 || used[u] != avoidUsed) && visit[u] != visitIter {
					visit[u] = visitIter
					parent[u] = v
					queue[tail] = u
					tail++
				}
			}
			// down
			if r+1 < h {
				u := v + w
				if !perm[u] && block[u] != dayIter && (avoidUsed == 0 || used[u] != avoidUsed) && visit[u] != visitIter {
					visit[u] = visitIter
					parent[u] = v
					queue[tail] = u
					tail++
				}
			}
			// left
			if c > 0 {
				u := v - 1
				if !perm[u] && block[u] != dayIter && (avoidUsed == 0 || used[u] != avoidUsed) && visit[u] != visitIter {
					visit[u] = visitIter
					parent[u] = v
					queue[tail] = u
					tail++
				}
			}
			// right
			if c+1 < w {
				u := v + 1
				if !perm[u] && block[u] != dayIter && (avoidUsed == 0 || used[u] != avoidUsed) && visit[u] != visitIter {
					visit[u] = visitIter
					parent[u] = v
					queue[tail] = u
					tail++
				}
			}
		}
		return false
	}

	for day := 0; day < q; day++ {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return
		}
		dayIter++
		for i := 0; i < k; i++ {
			var r, c int
			fmt.Fscan(reader, &r, &c)
			r--
			c--
			block[r*w+c] = dayIter
		}

		if !bfs(0) {
			fmt.Fprintln(writer, "NO")
			writer.Flush()
			continue
		}

		usedIter++
		// mark path from target to start
		u := target
		for u != start {
			if u != start && u != target {
				used[u] = usedIter
			}
			u = parent[u]
		}
		if bfs(usedIter) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
		writer.Flush()
	}
}
