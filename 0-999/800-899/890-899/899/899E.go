package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
)

type Segment struct {
	length int
	idx    int
	val    int
	left   int
	right  int
	alive  bool
}

var segments []Segment

type PQ []int

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
	a := segments[pq[i]]
	b := segments[pq[j]]
	if a.length != b.length {
		return a.length > b.length
	}
	return a.idx < b.idx
}
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(int)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	// build initial segments
	segments = make([]Segment, 0)
	var pq PQ
	start := 0
	for start < n {
		j := start
		for j < n && arr[j] == arr[start] {
			j++
		}
		segID := len(segments)
		segments = append(segments, Segment{
			length: j - start,
			idx:    start,
			val:    arr[start],
			left:   segID - 1,
			right:  -1,
			alive:  true,
		})
		if segID > 0 {
			segments[segID-1].right = segID
		}
		heap.Push(&pq, segID)
		start = j
	}

	ops := 0
	for pq.Len() > 0 {
		id := heap.Pop(&pq).(int)
		if !segments[id].alive {
			continue
		}
		ops++
		left := segments[id].left
		right := segments[id].right
		segments[id].alive = false
		if left != -1 {
			segments[left].right = right
		}
		if right != -1 {
			segments[right].left = left
		}
		if left != -1 && right != -1 && segments[left].alive && segments[right].alive && segments[left].val == segments[right].val {
			// merge left and right
			newID := len(segments)
			segments[left].alive = false
			segments[right].alive = false
			segments = append(segments, Segment{
				length: segments[left].length + segments[right].length,
				idx:    segments[left].idx,
				val:    segments[left].val,
				left:   segments[left].left,
				right:  segments[right].right,
				alive:  true,
			})
			if segments[newID].left != -1 {
				segments[segments[newID].left].right = newID
			}
			if segments[newID].right != -1 {
				segments[segments[newID].right].left = newID
			}
			heap.Push(&pq, newID)
		}
	}
	fmt.Fprintln(writer, ops)
}
