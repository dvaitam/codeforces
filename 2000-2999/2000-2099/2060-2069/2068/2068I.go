package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Node struct {
	cost int
	time int64
	idx  int
}

type PriorityQueue []Node

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].cost != pq[j].cost {
		return pq[i].cost < pq[j].cost
	}
	return pq[i].time < pq[j].time
}
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Node))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

const (
	dirUp = iota
	dirRight
	dirDown
	dirLeft
)

var dr = []int{-1, 0, 1, 0}
var dc = []int{0, 1, 0, -1}

func reflect(d int, cell byte) int {
	if cell == '/' {
		switch d {
		case dirUp:
			return dirRight
		case dirRight:
			return dirUp
		case dirDown:
			return dirLeft
		default:
			return dirDown
		}
	}
	// cell == '\\'
	switch d {
	case dirUp:
		return dirLeft
	case dirLeft:
		return dirUp
	case dirDown:
		return dirRight
	default:
		return dirDown
	}
}

func opposite(d int) int { return d ^ 2 }

func idxFrom(rcd [3]int, w int) int {
	r, c, d := rcd[0], rcd[1], rcd[2]
	return ((r*w + c) << 2) | d
}

func decodeIdx(idx, w int) (r, c, d int) {
	d = idx & 3
	cell := idx >> 2
	r = cell / w
	c = cell % w
	return
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var h, w int
	if _, err := fmt.Fscan(in, &h, &w); err != nil {
		return
	}
	grid := make([][]byte, h)
	var sr, sc int
	for i := 0; i < h; i++ {
		var line string
		fmt.Fscan(in, &line)
		grid[i] = []byte(line)
		for j := 0; j < w; j++ {
			if grid[i][j] == 'S' {
				sr, sc = i, j
				grid[i][j] = '.'
			}
		}
	}

	totalStates := h * w * 4
	const infCost = int(1e9)
	const infTime = int64(1 << 60)

	distCost := make([]int, totalStates)
	distTime := make([]int64, totalStates)
	parent := make([]int, totalStates)
	action := make([]byte, totalStates) // 0 normal/reflect, 1 destroy, 2 bounce
	for i := 0; i < totalStates; i++ {
		distCost[i] = infCost
		distTime[i] = infTime
		parent[i] = -1
	}

	pq := &PriorityQueue{}
	for d := 0; d < 4; d++ {
		idx := idxFrom([3]int{sr, sc, d}, w)
		distCost[idx] = 0
		distTime[idx] = 0
		heap.Push(pq, Node{0, 0, idx})
	}

	bestExitCost := infCost
	bestExitTime := infTime
	exitState := -1

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Node)
		if cur.cost != distCost[cur.idx] || cur.time != distTime[cur.idx] {
			continue
		}
		if cur.cost > bestExitCost || (cur.cost == bestExitCost && cur.time >= bestExitTime) {
			continue
		}
		r, c, d := decodeIdx(cur.idx, w)
		nr, nc := r+dr[d], c+dc[d]
		newTime := cur.time + 1
		// Check escape
		if nr < 0 || nr >= h || nc < 0 || nc >= w {
			if cur.cost < bestExitCost || (cur.cost == bestExitCost && newTime < bestExitTime) {
				bestExitCost = cur.cost
				bestExitTime = newTime
				exitState = cur.idx
			}
			continue
		}
		cell := grid[nr][nc]
		switch cell {
		case '#':
			nd := opposite(d)
			nidx := idxFrom([3]int{r, c, nd}, w) // stay in current cell
			if cur.cost < distCost[nidx] || (cur.cost == distCost[nidx] && newTime < distTime[nidx]) {
				distCost[nidx] = cur.cost
				distTime[nidx] = newTime
				parent[nidx] = cur.idx
				action[nidx] = 2 // bounce
				heap.Push(pq, Node{cur.cost, newTime, nidx})
			}
		case '.', 'S':
			nidx := idxFrom([3]int{nr, nc, d}, w)
			if cur.cost < distCost[nidx] || (cur.cost == distCost[nidx] && newTime < distTime[nidx]) {
				distCost[nidx] = cur.cost
				distTime[nidx] = newTime
				parent[nidx] = cur.idx
				action[nidx] = 0
				heap.Push(pq, Node{cur.cost, newTime, nidx})
			}
		case '/', '\\':
			// keep wall: reflect
			nd := reflect(d, cell)
			nidx := idxFrom([3]int{nr, nc, nd}, w)
			if cur.cost < distCost[nidx] || (cur.cost == distCost[nidx] && newTime < distTime[nidx]) {
				distCost[nidx] = cur.cost
				distTime[nidx] = newTime
				parent[nidx] = cur.idx
				action[nidx] = 0
				heap.Push(pq, Node{cur.cost, newTime, nidx})
			}
			// destroy wall: treat as free
			nidx2 := idxFrom([3]int{nr, nc, d}, w)
			ncost2 := cur.cost + 1
			if ncost2 < distCost[nidx2] || (ncost2 == distCost[nidx2] && newTime < distTime[nidx2]) {
				distCost[nidx2] = ncost2
				distTime[nidx2] = newTime
				parent[nidx2] = cur.idx
				action[nidx2] = 1 // destroy
				heap.Push(pq, Node{ncost2, newTime, nidx2})
			}
		}
	}

	if exitState == -1 {
		fmt.Println("NO")
		return
	}

	// Reconstruct path
	var ops [][3]int64
	state := exitState
	for parent[state] != -1 {
		if action[state] == 1 {
			r, c, _ := decodeIdx(state, w)
			ops = append(ops, [3]int64{distTime[state], int64(r + 1), int64(c + 1)})
		}
		state = parent[state]
	}
	// initial direction is direction of starting state
	_, _, startDir := decodeIdx(state, w)

	// reverse ops to chronological order
	for i, j := 0, len(ops)-1; i < j; i, j = i+1, j-1 {
		ops[i], ops[j] = ops[j], ops[i]
	}

	dirChars := []byte{'U', 'R', 'D', 'L'}
	fmt.Println("YES")
	fmt.Printf("%c\n", dirChars[startDir])
	fmt.Println(len(ops))
	for _, op := range ops {
		fmt.Printf("%d %d %d\n", op[0], op[1], op[2])
	}
}
