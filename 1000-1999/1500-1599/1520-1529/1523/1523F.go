package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

const inf int64 = 1 << 60

type State struct {
	time   int64
	mask   int
	pos    int
	quests int
}

type PQ []State

func (p PQ) Len() int            { return len(p) }
func (p PQ) Less(i, j int) bool  { return p[i].time < p[j].time }
func (p PQ) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *PQ) Push(x interface{}) { *p = append(*p, x.(State)) }
func (p *PQ) Pop() interface{} {
	old := *p
	x := old[len(old)-1]
	*p = old[:len(old)-1]
	return x
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	xa := make([]int, n)
	ya := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xa[i], &ya[i])
	}
	type Quest struct{ x, y, t int }
	qs := make([]Quest, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &qs[i].x, &qs[i].y, &qs[i].t)
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].t < qs[j].t })

	// precompute distances
	distTT := make([][]int, n)
	for i := 0; i < n; i++ {
		distTT[i] = make([]int, n)
		for j := 0; j < n; j++ {
			distTT[i][j] = abs(xa[i]-xa[j]) + abs(ya[i]-ya[j])
		}
	}
	distTQ := make([][]int, n)
	for i := 0; i < n; i++ {
		distTQ[i] = make([]int, m)
		for j := 0; j < m; j++ {
			distTQ[i][j] = abs(xa[i]-qs[j].x) + abs(ya[i]-qs[j].y)
		}
	}
	distQT := make([][]int, m)
	for j := 0; j < m; j++ {
		distQT[j] = make([]int, n)
		for i := 0; i < n; i++ {
			distQT[j][i] = abs(qs[j].x-xa[i]) + abs(qs[j].y-ya[i])
		}
	}
	distQQ := make([][]int, m)
	for i := 0; i < m; i++ {
		distQQ[i] = make([]int, m)
		for j := 0; j < m; j++ {
			distQQ[i][j] = abs(qs[i].x-qs[j].x) + abs(qs[i].y-qs[j].y)
		}
	}
	// precompute distMaskTower and distMaskQuest
	maskCnt := 1 << n
	distMaskTower := make([][]int, maskCnt)
	distMaskQuest := make([][]int, maskCnt)
	for mask := 0; mask < maskCnt; mask++ {
		distMaskTower[mask] = make([]int, n)
		distMaskQuest[mask] = make([]int, m)
		for i := 0; i < n; i++ {
			if mask == 0 {
				distMaskTower[mask][i] = 0
			} else {
				d := int(1e9)
				for j := 0; j < n; j++ {
					if mask&(1<<j) != 0 {
						if distTT[j][i] < d {
							d = distTT[j][i]
						}
					}
				}
				distMaskTower[mask][i] = d
			}
		}
		for j := 0; j < m; j++ {
			if mask == 0 {
				distMaskQuest[mask][j] = int(1e9)
			} else {
				d := int(1e9)
				for i := 0; i < n; i++ {
					if mask&(1<<i) != 0 {
						if distTQ[i][j] < d {
							d = distTQ[i][j]
						}
					}
				}
				distMaskQuest[mask][j] = d
			}
		}
	}

	towerPos := m // index for tower state
	posCnt := m + 1

	bestTime := make([][]int64, maskCnt)
	bestQuest := make([][]int, maskCnt)
	for i := range bestTime {
		bestTime[i] = make([]int64, posCnt)
		bestQuest[i] = make([]int, posCnt)
		for j := range bestTime[i] {
			bestTime[i][j] = inf
			bestQuest[i][j] = -1
		}
	}

	pq := &PQ{}
	heap.Init(pq)

	// initial state: spawn anywhere without towers
	if bestTime[0][towerPos] > 0 {
		bestTime[0][towerPos] = 0
		bestQuest[0][towerPos] = 0
		heap.Push(pq, State{time: 0, mask: 0, pos: towerPos, quests: 0})
	}
	for i := 0; i < n; i++ {
		mask := 1 << i
		if bestTime[mask][towerPos] > 0 {
			bestTime[mask][towerPos] = 0
			bestQuest[mask][towerPos] = 0
			heap.Push(pq, State{time: 0, mask: mask, pos: towerPos, quests: 0})
		}
	}
	for j := 0; j < m; j++ {
		if int64(qs[j].t) < bestTime[0][j] {
			bestTime[0][j] = int64(qs[j].t)
			bestQuest[0][j] = 1
			heap.Push(pq, State{time: int64(qs[j].t), mask: 0, pos: j, quests: 1})
		}
	}

	ans := 0
	for pq.Len() > 0 {
		s := heap.Pop(pq).(State)
		if s.time != bestTime[s.mask][s.pos] || s.quests != bestQuest[s.mask][s.pos] {
			continue
		}
		if s.quests > ans {
			ans = s.quests
		}
		if s.pos == towerPos {
			// move to quests
			for j := 0; j < m; j++ {
				arrive := s.time + int64(distMaskQuest[s.mask][j])
				if arrive <= int64(qs[j].t) {
					if s.quests+1 > bestQuest[s.mask][j] || (s.quests+1 == bestQuest[s.mask][j] && int64(qs[j].t) < bestTime[s.mask][j]) {
						bestQuest[s.mask][j] = s.quests + 1
						bestTime[s.mask][j] = int64(qs[j].t)
						heap.Push(pq, State{time: int64(qs[j].t), mask: s.mask, pos: j, quests: s.quests + 1})
					}
				}
			}
			// activate new towers
			for i := 0; i < n; i++ {
				if s.mask&(1<<i) == 0 {
					arrive := s.time + int64(distMaskTower[s.mask][i])
					newMask := s.mask | (1 << i)
					if arrive < bestTime[newMask][towerPos] || (arrive == bestTime[newMask][towerPos] && s.quests > bestQuest[newMask][towerPos]) {
						bestTime[newMask][towerPos] = arrive
						bestQuest[newMask][towerPos] = s.quests
						heap.Push(pq, State{time: arrive, mask: newMask, pos: towerPos, quests: s.quests})
					}
				}
			}
		} else {
			j := s.pos
			// teleport to towers
			if s.time < bestTime[s.mask][towerPos] || (s.time == bestTime[s.mask][towerPos] && s.quests > bestQuest[s.mask][towerPos]) {
				bestTime[s.mask][towerPos] = s.time
				bestQuest[s.mask][towerPos] = s.quests
				heap.Push(pq, State{time: s.time, mask: s.mask, pos: towerPos, quests: s.quests})
			}
			// move directly to later quests
			for k := j + 1; k < m; k++ {
				arrive := s.time + int64(distQQ[j][k])
				if arrive <= int64(qs[k].t) {
					if s.quests+1 > bestQuest[s.mask][k] || (s.quests+1 == bestQuest[s.mask][k] && int64(qs[k].t) < bestTime[s.mask][k]) {
						bestQuest[s.mask][k] = s.quests + 1
						bestTime[s.mask][k] = int64(qs[k].t)
						heap.Push(pq, State{time: int64(qs[k].t), mask: s.mask, pos: k, quests: s.quests + 1})
					}
				}
			}
			// activate new towers from this quest
			for i := 0; i < n; i++ {
				if s.mask&(1<<i) == 0 {
					d1 := distQT[j][i]
					if distMaskTower[s.mask][i] < d1 {
						d1 = distMaskTower[s.mask][i]
					}
					arrive := s.time + int64(d1)
					newMask := s.mask | (1 << i)
					if arrive < bestTime[newMask][towerPos] || (arrive == bestTime[newMask][towerPos] && s.quests > bestQuest[newMask][towerPos]) {
						bestTime[newMask][towerPos] = arrive
						bestQuest[newMask][towerPos] = s.quests
						heap.Push(pq, State{time: arrive, mask: newMask, pos: towerPos, quests: s.quests})
					}
				}
			}
		}
	}

	fmt.Fprintln(out, ans)
}
