package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "strconv"
)

// Item represents a bamboo in the priority queue
type Item struct {
   id int
   ai int64
}

// A MaxHeap implements heap.Interface for Items, ordering by ai descending
type MaxHeap []Item

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].ai > h[j].ai }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   var p int64
   fmt.Fscan(in, &n, &m, &k, &p)
   hArr := make([]int64, n)
   aArr := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &hArr[i], &aArr[i])
   }
   // compute upper bound for H
   var hi int64
   for i := 0; i < n; i++ {
       endH := hArr[i] + aArr[i]*int64(m)
       if endH > hi {
           hi = endH
       }
   }
   // binary search on H
   lo := int64(-1)
   // check function
   can := func(H int64) bool {
       // events[d]: list of bamboo ids to consider on day d
       events := make([][]int, m+2)
       times := make([]int, n)
       // initial events
       for i := 0; i < n; i++ {
           var t0 int
           if hArr[i] > H {
               t0 = 1
           } else {
               // days until exceeding: floor((H - h)/a) + 2
               t0 = int((H - hArr[i]) / aArr[i]) + 2
           }
           if t0 <= m {
               events[t0] = append(events[t0], i)
           }
       }
       // max-heap by a_i
       hq := &MaxHeap{}
       heap.Init(hq)
       for d := 1; d <= m; d++ {
           // add events
           for _, i := range events[d] {
               heap.Push(hq, Item{i, aArr[i]})
           }
           // use up to k hits
           for j := 0; j < k; j++ {
               if hq.Len() == 0 {
                   break
               }
               it := heap.Pop(hq).(Item)
               id := it.id
               times[id]++
               // compute height at start of next day after growth
               // current pre-growth height = h + a*(d-1) - times*p
               // after growth: +a
               hStart := hArr[id] + aArr[id]*int64(d) - int64(times[id])*p
               var tnext int
               if hStart > H {
                   tnext = d + 1
               } else {
                   // floor((H - hStart)/a) + 2 days after day d
                   tnext = d + int((H - hStart)/aArr[id]) + 2
               }
               if tnext <= m {
                   events[tnext] = append(events[tnext], id)
               }
           }
           // if any bamboo still exceeds before hit limit
           if hq.Len() > 0 {
               return false
           }
       }
       return true
   }
   for lo+1 < hi {
       mid := (lo + hi) / 2
       if can(mid) {
           hi = mid
       } else {
           lo = mid
       }
   }
   // output result
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, hi)
}
