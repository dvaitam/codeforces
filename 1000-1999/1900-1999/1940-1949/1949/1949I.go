package main

import (
	"bufio"
	"fmt"
	"os"
)

type Disk struct {
	x, y int64
	r    int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	disks := make([]Disk, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &disks[i].x, &disks[i].y, &disks[i].r)
	}

	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := disks[i].x - disks[j].x
			dy := disks[i].y - disks[j].y
			dist2 := dx*dx + dy*dy
			sr := disks[i].r + disks[j].r
			if dist2 == sr*sr {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}

	color := make([]int, n)
	for i := range color {
		color[i] = -1
	}

	for i := 0; i < n; i++ {
		if color[i] != -1 {
			continue
		}
		q := []int{i}
		color[i] = 0
		cnt := [2]int{1, 0}
		bip := true
		for h := 0; h < len(q); h++ {
			v := q[h]
			for _, u := range adj[v] {
				if color[u] == -1 {
					color[u] = color[v] ^ 1
					cnt[color[u]]++
					q = append(q, u)
				} else if color[u] == color[v] {
					bip = false
				}
			}
		}
		if bip && cnt[0] != cnt[1] {
			fmt.Fprintln(writer, "YES")
			return
		}
	}
	fmt.Fprintln(writer, "NO")
}
