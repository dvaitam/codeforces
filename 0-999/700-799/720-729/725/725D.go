package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// min-heap of int64
type IntHeap []int64

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *IntHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   var t1, w1 int64
   // read team 1
   fmt.Fscan(in, &t1, &w1)
   // opponents
   opp := make([]struct{ t, w int64 }, 0, n-1)
   for i := 2; i <= n; i++ {
       var t, w int64
       fmt.Fscan(in, &t, &w)
       opp = append(opp, struct{ t, w int64 }{t, w})
   }
   // sort by t descending
   sort.Slice(opp, func(i, j int) bool { return opp[i].t > opp[j].t })

   // heap of costs to disqualify
   h := &IntHeap{}
   heap.Init(h)
   // initial T and pointer
   T := t1
   j := 0
   // push all teams with t > T
   for j < len(opp) && opp[j].t > T {
       cost := opp[j].w - opp[j].t + 1
       heap.Push(h, cost)
       j++
   }
   // initial answer
   ans := h.Len() + 1
   used := int64(0)
   // try disqualifying more
   for h.Len() > 0 {
       cost := heap.Pop(h).(int64)
       if used+cost > t1 {
           break
       }
       used += cost
       T = t1 - used
       // add newly exceeding teams
       for j < len(opp) && opp[j].t > T {
           c := opp[j].w - opp[j].t + 1
           heap.Push(h, c)
           j++
       }
       cur := h.Len() + 1
       if cur < ans {
           ans = cur
       }
   }
   // print result
   fmt.Println(ans)
}
