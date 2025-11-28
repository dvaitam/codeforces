package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
)

// item holds the cost difference and position for a '?' placeholder
type item struct {
	diff int // cost change from ')' to '(' (A - B)
	idx  int // position in result
}

// minHeap implements a min-heap of items based on diff
type minHeap []item

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].diff < h[j].diff }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return
	}
	pattern := scanner.Text()

	n := len(pattern)
	if n%2 != 0 {
		fmt.Println("-1")
		return
	}

	res := []byte(pattern)
	h := &minHeap{}
	heap.Init(h)

	var totalCost int64
	balance := 0
	
	// We need to read pairs of costs for each '?'
	// Iterate through the pattern and process '?' as we go
	
	for i := 0; i < n; i++ {
		ch := res[i]
		if ch == '(' {
			balance++
		} else if ch == ')' {
			balance--
		} else { // ch == '?'
			// Read costs for this '?'
			if !scanner.Scan() {
				break // Should not happen based on valid input
			}
			aStr := scanner.Text()
			if !scanner.Scan() {
				break
			}
			bStr := scanner.Text()

			a, _ := strconv.Atoi(aStr)
			b, _ := strconv.Atoi(bStr)

			// Greedy strategy: Assume ')' first
			res[i] = ')'
			balance--
			totalCost += int64(b)

			// Push option to swap to '('
			// Cost change: new_cost - old_cost = a - b
			heap.Push(h, item{diff: a - b, idx: i})
		}

		// If balance is negative, we MUST flip a previous ')' to '('
		if balance < 0 {
			if h.Len() == 0 {
				fmt.Println("-1")
				return
			}
			it := heap.Pop(h).(item)
			
			// Apply flip
			res[it.idx] = '('
			totalCost += int64(it.diff)
			balance += 2 // ')' (-1) -> '(' (+1) is a change of +2
		}
	}

	if balance != 0 {
		fmt.Println("-1")
	} else {
		fmt.Println(totalCost)
		fmt.Println(string(res))
	}
}

