package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// MinHeap implements a min-heap of int64 values.
type MinHeap []int64

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
   *h = append(*h, x.(int64))
}

func (h *MinHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k1, k2, k3 int
   var t1, t2, t3 int64
   var n int
   fmt.Fscan(reader, &k1, &k2, &k3)
   fmt.Fscan(reader, &t1, &t2, &t3)
   fmt.Fscan(reader, &n)
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
   }

   // Effective number of servers per stage
   m1 := k1
   if m1 > n {
       m1 = n
   }
   m2 := k2
   if m2 > n {
       m2 = n
   }
   m3 := k3
   if m3 > n {
       m3 = n
   }

   // Initialize heaps for each stage with initial available times = 0
   h1 := &MinHeap{}
   h2 := &MinHeap{}
   h3 := &MinHeap{}
   heap.Init(h1)
   heap.Init(h2)
   heap.Init(h3)
   for i := 0; i < m1; i++ {
       heap.Push(h1, int64(0))
   }
   for i := 0; i < m2; i++ {
       heap.Push(h2, int64(0))
   }
   for i := 0; i < m3; i++ {
       heap.Push(h3, int64(0))
   }

   var maxFlow int64
   for i := 0; i < n; i++ {
       release := c[i]

       // Stage 1
       avail1 := heap.Pop(h1).(int64)
       var start1 int64
       if avail1 > release {
           start1 = avail1
       } else {
           start1 = release
       }
       comp1 := start1 + t1
       heap.Push(h1, comp1)

       // Stage 2
       avail2 := heap.Pop(h2).(int64)
       var start2 int64
       if avail2 > comp1 {
           start2 = avail2
       } else {
           start2 = comp1
       }
       comp2 := start2 + t2
       heap.Push(h2, comp2)

       // Stage 3
       avail3 := heap.Pop(h3).(int64)
       var start3 int64
       if avail3 > comp2 {
           start3 = avail3
       } else {
           start3 = comp2
       }
       comp3 := start3 + t3
       heap.Push(h3, comp3)

       // Update maximum flow time
       flow := comp3 - release
       if flow > maxFlow {
           maxFlow = flow
       }
   }

   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, maxFlow)
}
