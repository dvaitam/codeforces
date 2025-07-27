package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// intHeap is a max-heap of ints
type intHeap []int

func (h intHeap) Len() int            { return len(h) }
func (h intHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h intHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *intHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *intHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

type interval struct{ l, r int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, q int
   if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
       return
   }
   matrix := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &matrix[i])
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for qi := 0; qi < q; qi++ {
       var T string
       fmt.Fscan(reader, &T)
       // build match intervals for each row
       intervals := make([]interval, 0, n)
       for i := 0; i < n; i++ {
           row := matrix[i]
           j := 0
           for j < m {
               if row[j] != T[j] {
                   j++
                   continue
               }
               l := j
               for j < m && row[j] == T[j] {
                   j++
               }
               intervals = append(intervals, interval{l, j - 1})
           }
       }
       // sort intervals by start
       sort.Slice(intervals, func(i, j int) bool {
           if intervals[i].l != intervals[j].l {
               return intervals[i].l < intervals[j].l
           }
           return intervals[i].r > intervals[j].r
       })
       // greedy cover [0..m-1]
       pos := 0
       idx := 0
       segs := 0
       h := &intHeap{}
       heap.Init(h)
       ok := true
       for pos < m {
           for idx < len(intervals) && intervals[idx].l <= pos {
               heap.Push(h, intervals[idx].r)
               idx++
           }
           if h.Len() == 0 || (*h)[0] < pos {
               ok = false
               break
           }
           bestR := heap.Pop(h).(int)
           segs++
           pos = bestR + 1
       }
       if !ok {
           fmt.Fprintln(writer, -1)
       } else {
           if segs > 0 {
               fmt.Fprintln(writer, segs-1)
           } else {
               fmt.Fprintln(writer, 0)
           }
       }
   }
}
