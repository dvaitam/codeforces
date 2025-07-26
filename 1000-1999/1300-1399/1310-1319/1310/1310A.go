package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Item represents a category with initial count a and time cost t per added publication
type Item struct {
   a int64
   t int64
}

// MaxHeap is a max-heap of Items by t
type MaxHeap []Item

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].t > h[j].t }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
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
   fmt.Fscan(reader, &n)
   items := make([]Item, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &items[i].a)
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &items[i].t)
   }
   sort.Slice(items, func(i, j int) bool { return items[i].a < items[j].a })
   h := &MaxHeap{}
   heap.Init(h)
   var idx int
   var k int64
   if n > 0 {
       k = items[0].a
   }
   var total int64
   for idx < n || h.Len() > 0 {
       if h.Len() == 0 && idx < n && k < items[idx].a {
           k = items[idx].a
       }
       for idx < n && items[idx].a <= k {
           heap.Push(h, items[idx])
           idx++
       }
       top := heap.Pop(h).(Item)
       total += (k - top.a) * top.t
       k++
   }
   fmt.Fprintln(writer, total)
}
