package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

const maxCoord = 3000

type item struct {
	value int
	id    int
}

type minHeap []item

type maxHeap []item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].value < h[j].value }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(item))
}
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].value > h[j].value }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(item))
}
func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func getMin(h *minHeap, active []bool) int {
	for h.Len() > 0 {
		it := (*h)[0]
		if active[it.id] {
			return it.value
		}
		heap.Pop(h)
	}
	return 0
}

func getMax(h *maxHeap, active []bool) int {
	for h.Len() > 0 {
		it := (*h)[0]
		if active[it.id] {
			return it.value
		}
		heap.Pop(h)
	}
	return 0
}

var (
	n       int
	x1Arr   []int
	y1Arr   []int
	x2Arr   []int
	y2Arr   []int
	lenCol  []int
	minYCol []int
	maxYCol []int
	minXCol []int
	maxXCol []int
	contCol []bool
	writer  *bufio.Writer
)

func outputSquare(L, B, h int) bool {
	R := L + h
	T := B + h
	subset := make([]int, 0)
	areaSum := 0
	for i := 0; i < n; i++ {
		if x1Arr[i] >= L && x2Arr[i] <= R && y1Arr[i] >= B && y2Arr[i] <= T {
			subset = append(subset, i+1)
			areaSum += (x2Arr[i] - x1Arr[i]) * (y2Arr[i] - y1Arr[i])
		} else if x1Arr[i] < R && x2Arr[i] > L && y1Arr[i] < T && y2Arr[i] > B {
			return false
		}
	}
	if areaSum == h*h && len(subset) > 0 {
		fmt.Fprintf(writer, "YES %d\n", len(subset))
		for i, v := range subset {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
		return true
	}
	return false
}

func processRun(start, end, B, T int) bool {
	h := T - B
	if h <= 0 {
		return false
	}
	length := end - start
	if length < h {
		return false
	}
	type pair struct {
		idx   int
		value int
	}
	dequeMin := make([]pair, 0)
	dequeMax := make([]pair, 0)
	for x := start; x < end; x++ {
		valMin := minXCol[x]
		for len(dequeMin) > 0 && dequeMin[len(dequeMin)-1].value >= valMin {
			dequeMin = dequeMin[:len(dequeMin)-1]
		}
		dequeMin = append(dequeMin, pair{x, valMin})

		valMax := maxXCol[x]
		for len(dequeMax) > 0 && dequeMax[len(dequeMax)-1].value <= valMax {
			dequeMax = dequeMax[:len(dequeMax)-1]
		}
		dequeMax = append(dequeMax, pair{x, valMax})

		windowStart := x - h + 1
		if windowStart < start {
			continue
		}
		for len(dequeMin) > 0 && dequeMin[0].idx < windowStart {
			dequeMin = dequeMin[1:]
		}
		for len(dequeMax) > 0 && dequeMax[0].idx < windowStart {
			dequeMax = dequeMax[1:]
		}
		minVal := dequeMin[0].value
		maxVal := dequeMax[0].value
		L := windowStart
		R := L + h
		if minVal >= L && maxVal <= R {
			if outputSquare(L, B, h) {
				return true
			}
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	x1Arr = make([]int, n)
	y1Arr = make([]int, n)
	x2Arr = make([]int, n)
	y2Arr = make([]int, n)

	start := make([][]int, maxCoord+1)
	end := make([][]int, maxCoord+1)

	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x1Arr[i], &y1Arr[i], &x2Arr[i], &y2Arr[i])
		start[x1Arr[i]] = append(start[x1Arr[i]], i)
		end[x2Arr[i]] = append(end[x2Arr[i]], i)
	}

	lenCol = make([]int, maxCoord)
	minYCol = make([]int, maxCoord)
	maxYCol = make([]int, maxCoord)
	minXCol = make([]int, maxCoord)
	maxXCol = make([]int, maxCoord)
	contCol = make([]bool, maxCoord)

	active := make([]bool, n)
	sumLen := 0

	var minYHeap minHeap
	var maxYHeap maxHeap
	var minXHeap minHeap
	var maxXHeap maxHeap

	for x := 0; x <= maxCoord; x++ {
		for _, id := range end[x] {
			if active[id] {
				active[id] = false
				sumLen -= y2Arr[id] - y1Arr[id]
			}
		}
		if x == maxCoord {
			break
		}
		for _, id := range start[x] {
			active[id] = true
			sumLen += y2Arr[id] - y1Arr[id]
			heap.Push(&minYHeap, item{value: y1Arr[id], id: id})
			heap.Push(&maxYHeap, item{value: y2Arr[id], id: id})
			heap.Push(&minXHeap, item{value: x1Arr[id], id: id})
			heap.Push(&maxXHeap, item{value: x2Arr[id], id: id})
		}
		if sumLen == 0 {
			lenCol[x] = 0
			contCol[x] = false
			continue
		}
		minY := getMin(&minYHeap, active)
		maxY := getMax(&maxYHeap, active)
		minX := getMin(&minXHeap, active)
		maxX := getMax(&maxXHeap, active)
		lenCol[x] = sumLen
		minYCol[x] = minY
		maxYCol[x] = maxY
		minXCol[x] = minX
		maxXCol[x] = maxX
		if maxY-minY == sumLen {
			contCol[x] = true
		} else {
			contCol[x] = false
		}
	}

	found := false
	runStart := -1
	runB := 0
	runT := 0

	for x := 0; x < maxCoord && !found; x++ {
		if lenCol[x] > 0 && contCol[x] {
			if runStart == -1 {
				runStart = x
				runB = minYCol[x]
				runT = maxYCol[x]
			} else if minYCol[x] != runB || maxYCol[x] != runT {
				if processRun(runStart, x, runB, runT) {
					found = true
					break
				}
				runStart = x
				runB = minYCol[x]
				runT = maxYCol[x]
			}
		} else {
			if runStart != -1 {
				if processRun(runStart, x, runB, runT) {
					found = true
					break
				}
				runStart = -1
			}
		}
	}
	if !found && runStart != -1 {
		if processRun(runStart, maxCoord, runB, runT) {
			found = true
		}
	}

	if !found {
		fmt.Fprintln(writer, "NO")
	}
}
