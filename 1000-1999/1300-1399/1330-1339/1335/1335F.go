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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		colors := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &colors[i])
		}
		dirs := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &dirs[i])
		}
		solve(n, m, colors, dirs, writer)
	}
}

func solve(n, m int, colors, dirs []string, writer *bufio.Writer) {
	size := n * m
	to := make([]int, size)
	rev := make([][]int, size)
	indeg := make([]int, size)
	black := make([]bool, size)

	idx := func(i, j int) int { return i*m + j }

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			id := idx(i, j)
			black[id] = colors[i][j] == '0'
			switch dirs[i][j] {
			case 'U':
				to[id] = idx(i-1, j)
			case 'D':
				to[id] = idx(i+1, j)
			case 'L':
				to[id] = idx(i, j-1)
			case 'R':
				to[id] = idx(i, j+1)
			}
			rev[to[id]] = append(rev[to[id]], id)
			indeg[to[id]]++
		}
	}

	queue := make([]int, 0, size)
	removed := make([]bool, size)
	for i := 0; i < size; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
			removed[i] = true
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		w := to[v]
		indeg[w]--
		if indeg[w] == 0 && !removed[w] {
			removed[w] = true
			queue = append(queue, w)
		}
	}

	cycleId := make([]int, size)
	for i := range cycleId {
		cycleId[i] = -1
	}
	pos := make([]int, size)
	dist := make([]int, size)
	phase := make([]int, size)

	cycleLens := []int{}
	cycleBlack := [][]bool{}
	bfs := make([]int, 0, size)
	id := 0

	for i := 0; i < size; i++ {
		if !removed[i] && cycleId[i] == -1 {
			cur := i
			nodes := []int{}
			for {
				cycleId[cur] = id
				nodes = append(nodes, cur)
				cur = to[cur]
				if cur == i {
					break
				}
			}
			L := len(nodes)
			cycleLens = append(cycleLens, L)
			cycleBlack = append(cycleBlack, make([]bool, L))
			for idxNode, node := range nodes {
				pos[node] = idxNode
				dist[node] = 0
				phase[node] = idxNode
				if black[node] {
					cycleBlack[id][phase[node]] = true
				}
				bfs = append(bfs, node)
			}
			id++
		}
	}

	for head := 0; head < len(bfs); head++ {
		v := bfs[head]
		cid := cycleId[v]
		L := cycleLens[cid]
		for _, u := range rev[v] {
			if cycleId[u] == -1 {
				cycleId[u] = cid
				dist[u] = dist[v] + 1
				pos[u] = pos[v]
				phase[u] = (phase[v] - 1 + L) % L
				if black[u] {
					cycleBlack[cid][phase[u]] = true
				}
				bfs = append(bfs, u)
			}
		}
	}

	robots := 0
	blackCells := 0
	for cid := 0; cid < id; cid++ {
		robots += cycleLens[cid]
		for _, b := range cycleBlack[cid] {
			if b {
				blackCells++
			}
		}
	}
	fmt.Fprintln(writer, robots, blackCells)
}
