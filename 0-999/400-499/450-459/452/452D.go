package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
   *h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k, n1, n2, n3, t1, t2, t3 int
   if _, err := fmt.Fscan(reader, &k, &n1, &n2, &n3, &t1, &t2, &t3); err != nil {
       return
   }
   // Stage 1: washing
   wheap := &IntHeap{}
   heap.Init(wheap)
   for i := 0; i < n1; i++ {
       heap.Push(wheap, 0)
   }
   wfin := make([]int, 0, k)
   for i := 0; i < k; i++ {
       avail := heap.Pop(wheap).(int)
       finish := avail + t1
       wfin = append(wfin, finish)
       heap.Push(wheap, finish)
   }
   // Stage 2: drying
   dheap := &IntHeap{}
   heap.Init(dheap)
   for i := 0; i < n2; i++ {
       heap.Push(dheap, 0)
   }
   dfin := make([]int, 0, k)
   // wfin is non-decreasing by construction
   for _, wf := range wfin {
       avail := heap.Pop(dheap).(int)
       start := wf
       if avail > start {
           start = avail
       }
       finish := start + t2
       dfin = append(dfin, finish)
       heap.Push(dheap, finish)
   }
   // Stage 3: folding
   fheap := &IntHeap{}
   heap.Init(fheap)
   for i := 0; i < n3; i++ {
       heap.Push(fheap, 0)
   }
   var ans int
   // dfin is non-decreasing
   for _, df := range dfin {
       avail := heap.Pop(fheap).(int)
       start := df
       if avail > start {
           start = avail
       }
       finish := start + t3
       if finish > ans {
           ans = finish
       }
       heap.Push(fheap, finish)
   }
   // output result
   fmt.Println(ans)
