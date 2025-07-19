package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// pair holds the attributes x and y of an item
type pair struct {
   x, y int64
}

// IntHeap is a min-heap of int64 values
type IntHeap []int64

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
   *h = append(*h, x.(int64))
}

func (h *IntHeap) Pop() any {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   items := make([]pair, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &items[i].x, &items[i].y)
   }
   // sort by y ascending
   sort.Slice(items, func(i, j int) bool {
       return items[i].y < items[j].y
   })
   var total, ans int64
   h := &IntHeap{}
   heap.Init(h)
   // iterate from largest y to smallest
   for i := n - 1; i >= 0; i-- {
       x := items[i].x
       y := items[i].y
       heap.Push(h, x)
       total += x
       if h.Len() > m {
           // remove smallest x
           total -= heap.Pop(h).(int64)
       }
       if total*y > ans {
           ans = total * y
       }
   }
   fmt.Println(ans)
}
