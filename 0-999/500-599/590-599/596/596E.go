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

	var n, m, q int
	fmt.Fscan(reader, &n, &m, &q)
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}
	a := make([]int, 10)
	b := make([]int, 10)
	for i := 0; i < 10; i++ {
		fmt.Fscan(reader, &a[i], &b[i])
	}

	N := n * m
	digit := make([]int, N)
	nxt := make([]int, N)
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			id := x*m + y
			d := int(grid[x][y] - '0')
			digit[id] = d
			nx := x + a[d]
			ny := y + b[d]
			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				nxt[id] = id
			} else {
				nxt[id] = nx*m + ny
			}
		}
	}

	cycleID := make([]int, N)
	isCycle := make([]bool, N)
	color := make([]int, N)
	cycleMask := []int{0}
	var stack []int
	var curID int
	var dfs func(int)
	dfs = func(v int) {
		color[v] = 1
		stack = append(stack, v)
		to := nxt[v]
		if color[to] == 0 {
			dfs(to)
		} else if color[to] == 1 {
			curID++
			mask := 0
			for i := len(stack) - 1; ; i-- {
				w := stack[i]
				cycleID[w] = curID
				isCycle[w] = true
				mask |= 1 << uint(digit[w])
				if w == to {
					break
				}
			}
			if curID >= len(cycleMask) {
				cycleMask = append(cycleMask, mask)
			} else {
				cycleMask = append(cycleMask, 0)
				cycleMask[curID] = mask
			}
		}
		stack = stack[:len(stack)-1]
		color[v] = 2
		if !isCycle[v] {
			cycleID[v] = cycleID[to]
		}
	}

	for i := 0; i < N; i++ {
		if color[i] == 0 {
			dfs(i)
		}
	}

	// ensure cycleMask length = number of cycles +1
	if len(cycleMask) <= curID {
		cycleMask = append(cycleMask, make([]int, curID-len(cycleMask)+1)...)
	}

	// collect digits on cycles completely
	for i := 0; i < N; i++ {
		if isCycle[i] {
			id := cycleID[i]
			cycleMask[id] |= 1 << uint(digit[i])
		}
	}

	for ; q > 0; q-- {
		var s string
		fmt.Fscan(reader, &s)
		L := len(s)
		suffix := make([]int, L+1)
		for i := L - 1; i >= 0; i-- {
			suffix[i] = suffix[i+1] | (1 << uint(s[i]-'0'))
		}
		index := make([]int, N)
		visited := make([]bool, N)
		var dfsIdx func(int) int
		dfsIdx = func(v int) int {
			if visited[v] {
				return index[v]
			}
			visited[v] = true
			if isCycle[v] {
				index[v] = 0
			} else {
				idx := dfsIdx(nxt[v])
				if idx < L && int(s[idx]-'0') == digit[v] {
					idx++
				}
				index[v] = idx
			}
			return index[v]
		}
		for i := 0; i < N; i++ {
			if !visited[i] {
				dfsIdx(i)
			}
		}
		ok := false
		for i := 0; i < N && !ok; i++ {
			p := index[i]
			mask := cycleMask[cycleID[i]]
			if suffix[p]&^mask == 0 {
				ok = true
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
