package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// MaxHeap implements a max-heap of ints.
type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // Build doubled array and intervals
   N := 2 * n
   pairs := make([]struct{ l, j int }, N)
   for j := 1; j <= N; j++ {
       aj := a[(j-1)%n]
       l := j - aj
       if l < 1 {
           l = 1
       }
       pairs[j-1] = struct{ l, j int }{l, j}
   }
   // Sort by l ascending
   sort.Slice(pairs, func(i, j int) bool {
       return pairs[i].l < pairs[j].l
   })
   // Compute R[i]
   R := make([]int, N+2)
   h := &MaxHeap{}
   heap.Init(h)
   idx := 0
   for i := 1; i <= N; i++ {
       for idx < N && pairs[idx].l <= i {
           heap.Push(h, pairs[idx].j)
           idx++
       }
       // remove intervals that end <= i
       for h.Len() > 0 {
           top := (*h)[0]
           if top <= i {
               heap.Pop(h)
           } else {
               break
           }
       }
       if h.Len() == 0 {
           R[i] = i
       } else {
           R[i] = (*h)[0]
       }
   }
   R[N+1] = N + 1
   // Binary lifting
   maxLg := 0
   for (1 << maxLg) <= N {
       maxLg++
   }
   up := make([][]int, maxLg)
   up[0] = make([]int, N+2)
   for i := 1; i <= N; i++ {
       up[0][i] = R[i]
   }
   for k := 1; k < maxLg; k++ {
       up[k] = make([]int, N+2)
       for i := 1; i <= N; i++ {
           up[k][i] = up[k-1][up[k-1][i]]
       }
   }
   // Compute result
   var total int64
   for s := 1; s <= n; s++ {
       target := s + n - 1
       cur := s
       var cnt int64
       for k := maxLg - 1; k >= 0; k-- {
           nxt := up[k][cur]
           if nxt < target {
               cur = nxt
               cnt += 1 << k
           }
       }
       if cur < target {
           cnt++
       }
       total += cnt
   }
   // Output
   fmt.Println(total)
}
