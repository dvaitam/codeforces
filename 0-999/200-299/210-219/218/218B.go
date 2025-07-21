package main

import (
   "container/heap"
   "fmt"
)

// MaxHeap implements a max-heap of ints.
type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
   *h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

// MinHeap implements a min-heap of ints.
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
   *h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   seats := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Scan(&seats[i])
   }
   maxH := &MaxHeap{}
   minH := &MinHeap{}
   for _, v := range seats {
       heap.Push(maxH, v)
       heap.Push(minH, v)
   }
   heap.Init(maxH)
   heap.Init(minH)
   maxSum, minSum := 0, 0
   for i := 0; i < n; i++ {
       x := heap.Pop(maxH).(int)
       maxSum += x
       if x-1 > 0 {
           heap.Push(maxH, x-1)
       }
       y := heap.Pop(minH).(int)
       minSum += y
       if y-1 > 0 {
           heap.Push(minH, y-1)
       }
   }
   fmt.Printf("%d %d\n", maxSum, minSum)
}
