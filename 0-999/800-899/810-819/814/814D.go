package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type node struct {
	x, y, r float64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	maps := make([]node, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &maps[i].x, &maps[i].y, &maps[i].r)
	}
	// sort by descending radius
	sort.Slice(maps, func(i, j int) bool {
		return maps[i].r > maps[j].r
	})
	ans := make([]int, n)
	visit := make([]bool, n)
	// helper to check overlap (one contains other or disjoint)
	overlap := func(i, j int) bool {
		dx := maps[i].x - maps[j].x
		dy := maps[i].y - maps[j].y
		sr := maps[i].r + maps[j].r
		return dx*dx+dy*dy < sr*sr
	}
	// initialize with the largest circle
	visit[0] = true
	var sum float64 = maps[0].r * maps[0].r
	// process others
	for i := 1; i < n; i++ {
		flag := false
		for j := i - 1; j >= 0; j-- {
			if overlap(i, j) {
				flag = true
				if !visit[j] {
					sum += maps[i].r * maps[i].r
					visit[i] = true
					ans[i] = ans[j]
					break
				}
				// parent visited
				if ans[j] == 0 {
					ans[i] = 1
				} else {
					ans[i] = 0
				}
				flag = false
				for k := i - 1; k >= 0; k-- {
					if ans[k] == ans[i] && overlap(i, k) {
						flag = true
						if visit[k] {
							sum -= maps[i].r * maps[i].r
						} else {
							sum += maps[i].r * maps[i].r
							visit[i] = true
						}
						break
					}
				}
				if !flag {
					visit[i] = true
					sum += maps[i].r * maps[i].r
				}
				flag = true
				break
			}
		}
		if !flag {
			ans[i] = 0
			visit[i] = true
			sum += maps[i].r * maps[i].r
		}
	}
	fmt.Printf("%.9f\n", sum*math.Pi)
}
