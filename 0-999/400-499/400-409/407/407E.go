package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// min-heap of int64
type IntHeap []int64
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *IntHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

// max-heap of int64
type RevHeap []int64
func (h RevHeap) Len() int           { return len(h) }
func (h RevHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h RevHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *RevHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *RevHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   var d int64
   fmt.Fscan(in, &n, &k, &d)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   bestL, bestR := 1, 1
   bestLen := 1
   if d == 0 {
       // find longest run of equal numbers
       curL := 1
       for i := 1; i < n; i++ {
           if a[i] != a[i-1] {
               curL = i + 1
           }
           if i+1-curL+1 > bestLen {
               bestLen = i + 1 - curL + 1
               bestL = curL
               bestR = i + 1
           }
       }
       fmt.Printf("%d %d", bestL, bestR)
       return
   }
   // compute mods and b values
   mod := make([]int64, n)
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       m := a[i] % d
       if m < 0 {
           m += d
       }
       mod[i] = m
       b[i] = (a[i] - m) / d
   }
   // process segments of equal mod
   for s := 0; s < n; {
       e := s
       for e+1 < n && mod[e+1] == mod[s] {
           e++
       }
       // sliding window on [s..e]
       freq := make(map[int64]int)
       lazyMin := make(map[int64]int)
       lazyMax := make(map[int64]int)
       minH := &IntHeap{}
       maxH := &RevHeap{}
       heap.Init(minH)
       heap.Init(maxH)
       l := s
       for r := s; r <= e; r++ {
           x := b[r]
           heap.Push(minH, x)
           heap.Push(maxH, x)
           freq[x]++
           // remove duplicates
           for freq[x] > 1 {
               y := b[l]
               freq[y]--
               lazyMin[y]++
               lazyMax[y]++
               l++
           }
           // clean heaps
           for minH.Len() > 0 {
               top := (*minH)[0]
               if lazyMin[top] > 0 {
                   heap.Pop(minH)
                   lazyMin[top]--
               } else {
                   break
               }
           }
           for maxH.Len() > 0 {
               top := (*maxH)[0]
               if lazyMax[top] > 0 {
                   heap.Pop(maxH)
                   lazyMax[top]--
               } else {
                   break
               }
           }
           // update best if missing <= k
           if minH.Len() > 0 {
               // clean heaps before computing
               // (they are already clean from duplicates removal)
               minv := (*minH)[0]
               maxv := (*maxH)[0]
               wlen := int64(r - l + 1)
               missing := (maxv - minv + 1) - wlen
               if missing <= k {
                   currLen := r - l + 1
                   if currLen > bestLen {
                       bestLen = currLen
                       bestL = l + 1
                       bestR = r + 1
                   }
               }
           }
       }
       s = e + 1
   }
   fmt.Printf("%d %d", bestL, bestR)
}
