package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Opt represents a queued option with its value and index
type Opt struct {
   v   int64
   num int
}

// OptHeap implements a min-heap of Opt based on v
type OptHeap []Opt

func (h OptHeap) Len() int            { return len(h) }
func (h OptHeap) Less(i, j int) bool  { return h[i].v < h[j].v }
func (h OptHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *OptHeap) Push(x interface{}) { *h = append(*h, x.(Opt)) }
func (h *OptHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var m int64
   fmt.Fscan(in, &n, &m)
   a := make([]int, n)
   ans1 := make([]int64, n)
   ans2 := make([]bool, n)
   w := make([]int64, n)
   for i := 0; i < n; i++ {
       var ai int64
       fmt.Fscan(in, &ai)
       ans1[i] = ai / 100
       a[i] = int(ai % 100)
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &w[i])
   }

   var now int64
   h := &OptHeap{}
   heap.Init(h)
   for i := 0; i < n; i++ {
       if a[i] != 0 {
           cost := int64(a[i])
           benefit := w[i] * int64(100-a[i])
           if m >= cost {
               m -= cost
               heap.Push(h, Opt{benefit, i})
           } else {
               if h.Len() == 0 {
                   now += benefit
                   m += int64(100 - a[i])
                   ans2[i] = true
               } else {
                   top := (*h)[0]
                   if top.v < benefit {
                       heap.Pop(h)
                       heap.Push(h, Opt{benefit, i})
                       now += top.v
                       ans2[top.num] = true
                       m += int64(100 - a[i])
                   } else {
                       now += benefit
                       m += int64(100 - a[i])
                       ans2[i] = true
                   }
               }
           }
       }
   }

   fmt.Fprintln(out, now)
   for i := 0; i < n; i++ {
       if ans2[i] {
           fmt.Fprintln(out, ans1[i]+1, 0)
       } else {
           fmt.Fprintln(out, ans1[i], a[i])
       }
   }
}
