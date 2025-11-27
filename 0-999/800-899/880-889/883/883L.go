package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	infVal = int64(3e18)
	infPos = int(1e9)
)

type carInfo struct {
	time int64
	id   int
}

type carHeap []carInfo

func (h carHeap) Len() int { return len(h) }
func (h carHeap) Less(i, j int) bool {
	if h[i].time != h[j].time {
		return h[i].time < h[j].time
	}
	return h[i].id < h[j].id
}
func (h carHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *carHeap) Push(x interface{}) { *h = append(*h, x.(carInfo)) }
func (h *carHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type busyCar struct {
	finish int64
	id     int
	pos    int
}

type busyHeap []busyCar

func (h busyHeap) Len() int { return len(h) }
func (h busyHeap) Less(i, j int) bool {
	return h[i].finish < h[j].finish
}
func (h busyHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *busyHeap) Push(x interface{}) { *h = append(*h, x.(busyCar)) }
func (h *busyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type nodeRight struct {
	pos  int
	time int64
	id   int
}
type nodeLeft struct {
	negPos int
	time   int64
	id     int
}

var (
	treeRight []nodeRight
	treeLeft  []nodeLeft
	locHeaps  []carHeap
	nHouses   int
)

func lessRight(a, b nodeRight) bool {
	if a.pos != b.pos {
		return a.pos < b.pos
	}
	if a.time != b.time {
		return a.time < b.time
	}
	return a.id < b.id
}

func lessLeft(a, b nodeLeft) bool {
	if a.negPos != b.negPos {
		return a.negPos < b.negPos
	}
	if a.time != b.time {
		return a.time < b.time
	}
	return a.id < b.id
}

func build(idx, l, r int) {
	if l == r {
		treeRight[idx] = nodeRight{pos: infPos, time: infVal, id: infPos}
		treeLeft[idx] = nodeLeft{negPos: infPos, time: infVal, id: infPos}
		return
	}
	mid := (l + r) / 2
	build(2*idx, l, mid)
	build(2*idx+1, mid+1, r)
	treeRight[idx] = treeRight[2*idx]
	treeLeft[idx] = treeLeft[2*idx]
}

func update(idx, l, r, pos int) {
	if l == r {
		if len(locHeaps[pos]) > 0 {
			best := locHeaps[pos][0]
			treeRight[idx] = nodeRight{pos: pos, time: best.time, id: best.id}
			treeLeft[idx] = nodeLeft{negPos: -pos, time: best.time, id: best.id}
		} else {
			treeRight[idx] = nodeRight{pos: infPos, time: infVal, id: infPos}
			treeLeft[idx] = nodeLeft{negPos: infPos, time: infVal, id: infPos}
		}
		return
	}
	mid := (l + r) / 2
	if pos <= mid {
		update(2*idx, l, mid, pos)
	} else {
		update(2*idx+1, mid+1, r, pos)
	}
	if lessRight(treeRight[2*idx], treeRight[2*idx+1]) {
		treeRight[idx] = treeRight[2*idx]
	} else {
		treeRight[idx] = treeRight[2*idx+1]
	}
	if lessLeft(treeLeft[2*idx], treeLeft[2*idx+1]) {
		treeLeft[idx] = treeLeft[2*idx]
	} else {
		treeLeft[idx] = treeLeft[2*idx+1]
	}
}

func queryR(idx, l, r, ql, qr int) nodeRight {
	if ql > r || qr < l {
		return nodeRight{pos: infPos, time: infVal, id: infPos}
	}
	if ql <= l && r <= qr {
		return treeRight[idx]
	}
	mid := (l + r) / 2
	v1 := queryR(2*idx, l, mid, ql, qr)
	v2 := queryR(2*idx+1, mid+1, r, ql, qr)
	if lessRight(v1, v2) {
		return v1
	}
	return v2
}

func queryL(idx, l, r, ql, qr int) nodeLeft {
	if ql > r || qr < l {
		return nodeLeft{negPos: infPos, time: infVal, id: infPos}
	}
	if ql <= l && r <= qr {
		return treeLeft[idx]
	}
	mid := (l + r) / 2
	v1 := queryL(2*idx, l, mid, ql, qr)
	v2 := queryL(2*idx+1, mid+1, r, ql, qr)
	if lessLeft(v1, v2) {
		return v1
	}
	return v2
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

var scanner *bufio.Scanner

func scanInt() int {
	scanner.Scan()
	x, _ := strconv.Atoi(scanner.Text())
	return x
}

func scanInt64() int64 {
	scanner.Scan()
	x, _ := strconv.ParseInt(scanner.Text(), 10, 64)
	return x
}

func main() {
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	n := scanInt()
	k := scanInt()
	m := scanInt()

	nHouses = n
	locHeaps = make([]carHeap, n+1)
	treeRight = make([]nodeRight, 4*n+100)
	treeLeft = make([]nodeLeft, 4*n+100)

	build(1, 1, n)

	for i := 1; i <= k; i++ {
		pos := scanInt()
		c := carInfo{time: 0, id: i}
		heap.Push(&locHeaps[pos], c)
		update(1, 1, n, pos)
	}

	busyHeap := &busyHeap{}
	heap.Init(busyHeap)

	for i := 0; i < m; i++ {
		tj := scanInt64()
		aj := scanInt()
		bj := scanInt()

		currTime := tj

		// Release cars that have finished by current time.
		for busyHeap.Len() > 0 {
			top := (*busyHeap)[0]
			if top.finish <= currTime {
				heap.Pop(busyHeap)
				c := carInfo{time: top.finish, id: top.id}
				heap.Push(&locHeaps[top.pos], c)
				update(1, 1, n, top.pos)
			} else {
				break
			}
		}

		// If no cars are available, fast-forward to next finish time and release.
		if treeRight[1].id == infPos && busyHeap.Len() > 0 {
			currTime = (*busyHeap)[0].finish
			for busyHeap.Len() > 0 {
				top := (*busyHeap)[0]
				if top.finish <= currTime {
					heap.Pop(busyHeap)
					c := carInfo{time: top.finish, id: top.id}
					heap.Push(&locHeaps[top.pos], c)
					update(1, 1, n, top.pos)
				} else {
					break
				}
			}
		}

		candR := queryR(1, 1, n, aj, n)
		candL := queryL(1, 1, n, 1, aj)

		bestIsR := false
		distR := int64(math.MaxInt64)
		if candR.id != infPos {
			distR = int64(candR.pos - aj)
		}
		distL := int64(math.MaxInt64)
		if candL.id != infPos {
			distL = int64(aj + candL.negPos)
		}

		if candR.id == infPos {
			bestIsR = false
		} else if candL.id == infPos {
			bestIsR = true
		} else {
			if distR < distL {
				bestIsR = true
			} else if distL < distR {
				bestIsR = false
			} else {
				if candR.time < candL.time {
					bestIsR = true
				} else if candL.time < candR.time {
					bestIsR = false
				} else {
					if candR.id < candL.id {
						bestIsR = true
					} else {
						bestIsR = false
					}
				}
			}
		}

		var selPos, selID int
		var selTime int64

		if bestIsR {
			selPos = candR.pos
			selID = candR.id
			selTime = candR.time
		} else {
			selPos = -candL.negPos
			selID = candL.id
			selTime = candL.time
		}

		// Remove from heap; the chosen car might have a later availability.
		heap.Pop(&locHeaps[selPos])
		update(1, 1, n, selPos)

		pickupTime := max64(currTime+int64(abs(selPos-aj)), selTime+int64(abs(selPos-aj)))
		waitTime := pickupTime - tj
		dropTime := pickupTime + int64(abs(aj-bj))

		fmt.Fprintf(writer, "%d %d\n", selID, waitTime)

		heap.Push(busyHeap, busyCar{finish: dropTime, id: selID, pos: bj})
	}
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
