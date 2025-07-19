package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

const mod = 1000000007

// Min-heap for ints
type MinHeap []int
func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

// Max-heap for ints
type MaxHeap []int
func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N int
   var X, Y int64
   fmt.Fscan(reader, &N, &X, &Y)
   intervals := make([][2]int, N)
   for i := 0; i < N; i++ {
       var f, s int
       fmt.Fscan(reader, &f, &s)
       intervals[i][0] = f
       intervals[i][1] = s
   }
   sort.Slice(intervals, func(i, j int) bool {
       if intervals[i][0] != intervals[j][0] {
           return intervals[i][0] < intervals[j][0]
       }
       return intervals[i][1] < intervals[j][1]
   })

   rem := &MinHeap{}
   avail := &MaxHeap{}
   heap.Init(rem)
   heap.Init(avail)

   var ans int64
   for _, iv := range intervals {
       f, s := iv[0], iv[1]
       // move finished to avail
       for rem.Len() > 0 {
           top := (*rem)[0]
           if top < f {
               heap.Pop(rem)
               heap.Push(avail, top)
           } else {
               break
           }
       }
       if avail.Len() == 0 {
           // no available, start new
           heap.Push(avail, f)
           ans = (ans + X) % mod
       }
       // reuse or the one just pushed
       topAvail := (*avail)[0]
       delta := int64(f-topAvail) * Y
       if delta > X {
           delta = X
       }
       ans = (ans + delta) % mod
       // duration cost
       ans = (ans + int64(s-f)*Y % mod) % mod
       // assign this interval
       heap.Pop(avail)
       heap.Push(rem, s)
   }
   fmt.Fprintln(writer, ans)
}
