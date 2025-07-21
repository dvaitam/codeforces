// Package main solves the DH video recompression scheduling problem.
// See problemD.txt for details.
package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// A min-heap of server finish times.
type minHeap []int64

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(int64)) }
func (h *minHeap) Pop() interface{} {
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

   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   s := make([]int64, n)
   m := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &s[i], &m[i])
   }
   res := make([]int64, n)
   // busy heap holds finish times
   h := &minHeap{}
   heap.Init(h)
   free := k
   // waiting queue: indices of jobs waiting
   type job struct{ idx int; m int64 }
   queue := make([]job, 0, n)

   // process all finish events up to time t
   process := func(t int64) {
       for h.Len() > 0 {
           ft := (*h)[0]
           if ft > t {
               break
           }
           heap.Pop(h)
           if len(queue) > 0 {
               j := queue[0]
               queue = queue[1:]
               start := ft
               finish := start + j.m
               res[j.idx] = finish
               heap.Push(h, finish)
           } else {
               free++
           }
       }
   }

   for i := 0; i < n; i++ {
       ti := s[i]
       process(ti)
       if free > 0 {
           free--
           start := ti
           finish := start + m[i]
           res[i] = finish
           heap.Push(h, finish)
       } else {
           queue = append(queue, job{i, m[i]})
       }
   }
   // finish remaining queued jobs
   for len(queue) > 0 {
       ft := heap.Pop(h).(int64)
       j := queue[0]
       queue = queue[1:]
       start := ft
       finish := start + j.m
       res[j.idx] = finish
       heap.Push(h, finish)
   }
   // output
   for i := 0; i < n; i++ {
       fmt.Fprintln(out, res[i])
   }
}
