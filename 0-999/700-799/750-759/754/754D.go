package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

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

// Interval represents a segment with left, right, and original id.
type Interval struct {
   l, r, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   intervals := make([]Interval, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &intervals[i].l, &intervals[i].r)
       intervals[i].id = i + 1
   }
   sort.Slice(intervals, func(i, j int) bool {
       return intervals[i].l < intervals[j].l
   })
   rmin := make([]int, n)
   h := &MinHeap{}
   heap.Init(h)
   tot := 0
   const negInf = -1234567890
   for i := 0; i < n; i++ {
       if tot == k {
           t := (*h)[0]
           if t <= intervals[i].r {
               heap.Pop(h)
               heap.Push(h, intervals[i].r)
           }
           rmin[i] = (*h)[0]
       } else {
           heap.Push(h, intervals[i].r)
           tot++
           if tot != k {
               rmin[i] = negInf
           } else {
               rmin[i] = (*h)[0]
           }
       }
   }
   ans := 0
   ansi := 0
   for i := 0; i < n; i++ {
       cur := rmin[i] - intervals[i].l + 1
       if cur < 0 {
           cur = 0
       }
       if cur > ans {
           ans = cur
           ansi = i
       }
   }
   // Select intervals up to ansi
   selected := make([]Interval, ansi+1)
   copy(selected, intervals[:ansi+1])
   // Sort by r descending
   sort.Slice(selected, func(i, j int) bool {
       return selected[i].r > selected[j].r
   })
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
   if ans == 0 {
       for i := 1; i <= k; i++ {
           fmt.Fprintf(writer, "%d ", i)
       }
       return
   }
   for i := 0; i < k; i++ {
       fmt.Fprintf(writer, "%d ", selected[i].id)
   }
}
