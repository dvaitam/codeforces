package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"sort"
)

// intHeap implements heap.Interface for ints (min-heap)
type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// returnItem is used for tracking when a passenger returns
// after receiving water

type returnItem struct {
	time int64
	idx  int
}

type retHeap []returnItem

func (h retHeap) Len() int            { return len(h) }
func (h retHeap) Less(i, j int) bool  { return h[i].time < h[j].time }
func (h retHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *retHeap) Push(x interface{}) { *h = append(*h, x.(returnItem)) }
func (h *retHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	var p int64
	if _, err := fmt.Fscan(reader, &n, &p); err != nil {
		return
	}
	times := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &times[i])
	}

	// passengers sorted by time of leaving seat
	type pair struct {
		t   int64
		idx int
	}
	arr := make([]pair, n)
	for i := 0; i < n; i++ {
		arr[i] = pair{t: times[i], idx: i}
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].t == arr[j].t {
			return arr[i].idx < arr[j].idx
		}
		return arr[i].t < arr[j].t
	})

	var avail intHeap       // indices whose time has come but not yet left
	var empty intHeap       // indices currently empty (min-heap)
	var returns retHeap     // return events
	queue := make([]int, 0) // queue at the tank
	emptyMark := make([]bool, n)

	heap.Init(&avail)
	heap.Init(&empty)
	heap.Init(&returns)

	res := make([]int64, n)

	var time int64
	var tankBusy int64
	pos := 0
	served := 0

	const inf int64 = 1 << 60

	for served < n {
		// add arrivals
		for pos < n && arr[pos].t <= time {
			heap.Push(&avail, arr[pos].idx)
			pos++
		}
		// process returns
		for returns.Len() > 0 && returns[0].time <= time {
			item := heap.Pop(&returns).(returnItem)
			emptyMark[item.idx] = false
			heap.Push(&avail, item.idx) // after returning, passenger cannot leave again, but pushing to avail won't hurt since we'll check marks
		}
		// clean up empty heap
		for empty.Len() > 0 {
			x := empty[0]
			if emptyMark[x] {
				break
			}
			heap.Pop(&empty)
		}
		minEmpty := n + 1
		if empty.Len() > 0 {
			minEmpty = empty[0]
		}
		// try to move one passenger from avail to queue
		moved := false
		for avail.Len() > 0 {
			idx := heap.Pop(&avail).(int)
			if emptyMark[idx] {
				// already left earlier, skip
				continue
			}
			if idx+1 < minEmpty { // seat numbering is 1-based in description
				emptyMark[idx] = true
				heap.Push(&empty, idx)
				queue = append(queue, idx)
				moved = true
			} else {
				heap.Push(&avail, idx)
			}
			break
		}
		if moved {
			// only one passenger leaves per moment
			time++
			continue
		}
		// if tank is free and someone is waiting
		if tankBusy <= time && len(queue) > 0 {
			idx := queue[0]
			queue = queue[1:]
			start := time
			if tankBusy > start {
				start = tankBusy
			}
			finish := start + p
			res[idx] = finish
			tankBusy = finish
			heap.Push(&returns, returnItem{time: finish, idx: idx})
			served++
			time = start
			continue
		}

		// advance time to next event
		next := inf
		if pos < n {
			if arr[pos].t < next {
				next = arr[pos].t
			}
		}
		if tankBusy > time && tankBusy < next {
			next = tankBusy
		}
		if returns.Len() > 0 && returns[0].time < next {
			next = returns[0].time
		}
		if next == inf {
			break
		}
		if next > time {
			time = next
		} else {
			time++
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, res[i])
	}
	fmt.Fprintln(writer)
}
