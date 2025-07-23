package main

import (
	"bufio"
	"fmt"
	"os"
)

var n, k int
var grid [][]byte
var comp [][]int
var compSize []int
var used []int
var tag int

var dirs = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func bfs() {
	comp = make([][]int, n)
	for i := 0; i < n; i++ {
		comp[i] = make([]int, n)
	}
	compSize = []int{0}
	cid := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '.' && comp[i][j] == 0 {
				cid++
				q := make([][2]int, 0)
				q = append(q, [2]int{i, j})
				comp[i][j] = cid
				size := 0
				for h := 0; h < len(q); h++ {
					x, y := q[h][0], q[h][1]
					size++
					for _, d := range dirs {
						nx, ny := x+d[0], y+d[1]
						if nx >= 0 && nx < n && ny >= 0 && ny < n && grid[nx][ny] == '.' && comp[nx][ny] == 0 {
							comp[nx][ny] = cid
							q = append(q, [2]int{nx, ny})
						}
					}
				}
				compSize = append(compSize, size)
			}
		}
	}
	used = make([]int, cid+1)
}

func evaluate(i, j int, cnt []int) int {
	size := k * k
	tag++
	for r := i; r < i+k; r++ {
		if j > 0 && grid[r][j-1] == '.' {
			id := comp[r][j-1]
			if used[id] != tag {
				used[id] = tag
				size += compSize[id] - cnt[id]
			}
		}
		if j+k < n && grid[r][j+k] == '.' {
			id := comp[r][j+k]
			if used[id] != tag {
				used[id] = tag
				size += compSize[id] - cnt[id]
			}
		}
	}
	for c := j; c < j+k; c++ {
		if i > 0 && grid[i-1][c] == '.' {
			id := comp[i-1][c]
			if used[id] != tag {
				used[id] = tag
				size += compSize[id] - cnt[id]
			}
		}
		if i+k < n && grid[i+k][c] == '.' {
			id := comp[i+k][c]
			if used[id] != tag {
				used[id] = tag
				size += compSize[id] - cnt[id]
			}
		}
	}
	return size
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	grid = make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}
	bfs()
	cnt := make([]int, len(used))
	ans := 0
	for i := 0; i+k <= n; i++ {
		for t := range cnt {
			cnt[t] = 0
		}
		for r := i; r < i+k; r++ {
			for c := 0; c < k; c++ {
				id := comp[r][c]
				if id > 0 {
					cnt[id]++
				}
			}
		}
		cur := evaluate(i, 0, cnt)
		if cur > ans {
			ans = cur
		}
		for j := 1; j+k <= n; j++ {
			for r := i; r < i+k; r++ {
				id := comp[r][j-1]
				if id > 0 {
					cnt[id]--
				}
				id = comp[r][j+k-1]
				if id > 0 {
					cnt[id]++
				}
			}
			cur = evaluate(i, j, cnt)
			if cur > ans {
				ans = cur
			}
		}
	}
	fmt.Println(ans)
}
