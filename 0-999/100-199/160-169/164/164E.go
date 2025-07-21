package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Item represents a scheduled task in the heap
type Item struct {
   t   int64
   idx int
}

// MaxHeap is a max-heap of Items by t, then idx
type MaxHeap []Item

func (h MaxHeap) Len() int { return len(h) }
// We want largest t, and for ties largest idx, at the top (Pop)
func (h MaxHeap) Less(i, j int) bool {
   if h[i].t != h[j].t {
       return h[i].t > h[j].t
   }
   return h[i].idx > h[j].idx
}
func (h MaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
   *h = append(*h, x.(Item))
}
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   item := old[n-1]
   *h = old[0 : n-1]
   return item
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   l := make([]int64, n)
   r := make([]int64, n)
   t := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &l[i], &r[i], &t[i])
   }
   ans := int64(0)
   h := &MaxHeap{}
   heap.Init(h)
   out := make([]int, n)
   for i := 0; i < n; i++ {
       // try to schedule normally
       // compute start time
       start := ans + 1
       if start < l[i] {
           start = l[i]
       }
       end := start + t[i] - 1
       if end <= r[i] {
           // schedule
           ans = end
           heap.Push(h, Item{t: t[i], idx: i + 1})
           out[i] = 0
       } else if h.Len() > 0 && (*h)[0].t > t[i] {
           // try replace largest
           top := heap.Pop(h).(Item)
           // simulate removal
           ans2 := ans - top.t
           start2 := ans2 + 1
           if start2 < l[i] {
               start2 = l[i]
           }
           end2 := start2 + t[i] - 1
           if end2 <= r[i] {
               // perform replacement
               ans = end2
               heap.Push(h, Item{t: t[i], idx: i + 1})
               out[i] = top.idx
           } else {
               // revert removal
               heap.Push(h, top)
               out[i] = -1
           }
       } else {
           out[i] = -1
       }
   }
   for i := 0; i < n; i++ {
       fmt.Fprintln(writer, out[i])
   }
}
