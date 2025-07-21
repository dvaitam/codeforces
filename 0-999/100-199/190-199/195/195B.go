package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// basket represents a basket in the priority queue
type basket struct {
   cnt   int // number of balls currently in basket
   dist2 int // twice the distance to center: |m+1-2*i|
   idx   int // basket index (1-based)
}

// basketHeap implements heap.Interface, ordering by cnt, then dist2, then idx
type basketHeap []basket

func (h basketHeap) Len() int { return len(h) }
func (h basketHeap) Less(i, j int) bool {
   if h[i].cnt != h[j].cnt {
       return h[i].cnt < h[j].cnt
   }
   if h[i].dist2 != h[j].dist2 {
       return h[i].dist2 < h[j].dist2
   }
   return h[i].idx < h[j].idx
}
func (h basketHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *basketHeap) Push(x interface{}) {
   *h = append(*h, x.(basket))
}
func (h *basketHeap) Pop() interface{} {
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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }

   // initialize heap with m baskets
   h := make(basketHeap, 0, m)
   // center factor for distance calculation
   center := m + 1
   for i := 1; i <= m; i++ {
       d2 := center - 2*i
       if d2 < 0 {
           d2 = -d2
       }
       h = append(h, basket{cnt: 0, dist2: d2, idx: i})
   }
   heap.Init(&h)

   // simulate placing balls
   for i := 0; i < n; i++ {
       b := heap.Pop(&h).(basket)
       // output basket index for ball i+1
       fmt.Fprintln(writer, b.idx)
       b.cnt++
       heap.Push(&h, b)
   }
}
