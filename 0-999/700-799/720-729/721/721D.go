package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// item represents a value and its position.
type item struct {
   v   int64
   pos int
}

// maxHeap implements a max-heap of item based on v.
type maxHeap []item

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].v > h[j].v }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *maxHeap) Push(x interface{}) {
   *h = append(*h, x.(item))
}

func (h *maxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[:n-1]
   return it
}

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var x int64
   if _, err := fmt.Fscan(in, &n, &m, &x); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   num := 0
   h := &maxHeap{}
   heap.Init(h)
   for i, v := range a {
       if v < 0 {
           num++
       }
       heap.Push(h, item{v: abs(v), pos: i})
   }
   for k := 0; k < m; k++ {
       it := heap.Pop(h).(item)
       p := it.pos
       if a[p] < 0 {
           num--
       }
       if num&1 == 1 {
           a[p] += x
       } else {
           a[p] -= x
       }
       if a[p] < 0 {
           num++
       }
       heap.Push(h, item{v: abs(a[p]), pos: p})
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, v := range a {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprintf(out, "%d", v)
   }
   out.WriteByte('\n')
}
