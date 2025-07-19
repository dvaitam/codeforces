package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Segment represents an interval with start l, end r, and original index
type Segment struct {
   l, r int
   idx  int
}

// Item is an element in the priority queue
type Item struct {
   r   int // end
   idx int // original index
}

// A MinHeap implements heap.Interface and holds Items.
type MinHeap []Item

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].r < h[j].r }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
   *h = append(*h, x.(Item))
}

func (h *MinHeap) Pop() interface{} {
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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   segs := make([]Segment, m)
   for i := 0; i < m; i++ {
       var r, c int
       fmt.Fscan(reader, &r, &c)
       // map to 1-based l = n-r+1, end = c
       segs[i] = Segment{l: n - r + 1, r: c, idx: i + 1}
   }
   sort.Slice(segs, func(i, j int) bool {
       return segs[i].l < segs[j].l
   })

   h := &MinHeap{}
   heap.Init(h)
   ans := make([]int, 0, m)
   j := 0
   for l := 1; l <= n; l++ {
       for j < m && segs[j].l == l {
           heap.Push(h, Item{r: segs[j].r, idx: segs[j].idx})
           j++
       }
       // remove expired
       for h.Len() > 0 && (*h)[0].r < l {
           heap.Pop(h)
       }
       if h.Len() > 0 {
           itm := heap.Pop(h).(Item)
           ans = append(ans, itm.idx)
       }
   }
   // output
   fmt.Fprintln(writer, len(ans))
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
