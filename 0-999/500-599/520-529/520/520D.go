package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

const MOD = 1000000009

// MinHeap is a min-heap of ints
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

// MaxHeap is a max-heap of ints
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
   in := bufio.NewReader(os.Stdin)
   var m int
   fmt.Fscan(in, &m)
   xs := make([]int, m)
   ys := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &xs[i], &ys[i])
   }
   // map coordinates to index
   idx := make(map[uint64]int, m)
   for i := 0; i < m; i++ {
       key := (uint64(xs[i]) << 32) | uint64(uint32(ys[i]))
       idx[key] = i
   }
   // supporters and dependents
   supporters := make([][]int, m)
   dependents := make([][]int, m)
   supportCount := make([]int, m)
   for i := 0; i < m; i++ {
       x, y := xs[i], ys[i]
       for dx := -1; dx <= 1; dx++ {
           key := (uint64(x+dx) << 32) | uint64(uint32(y-1))
           if j, ok := idx[key]; ok {
               supporters[i] = append(supporters[i], j)
           }
       }
       supportCount[i] = len(supporters[i])
       for _, j := range supporters[i] {
           dependents[j] = append(dependents[j], i)
       }
   }
   // eligibility and removed
   inSet := make([]bool, m)
   removed := make([]bool, m)
   var minh MinHeap
   var maxh MaxHeap
   heap.Init(&minh)
   heap.Init(&maxh)
   // initial eligible
   for i := 0; i < m; i++ {
       ok := true
       for _, u := range dependents[i] {
           if supportCount[u] < 2 {
               ok = false
               break
           }
       }
       if ok {
           inSet[i] = true
           heap.Push(&minh, i)
           heap.Push(&maxh, i)
       }
   }
   seq := make([]int, 0, m)
   // game
   for t := 0; t < m; t++ {
       var v int
       if t%2 == 0 {
           // Vasya, max
           for {
               top := heap.Pop(&maxh).(int)
               if !removed[top] && inSet[top] {
                   v = top
                   break
               }
           }
       } else {
           // Petya, min
           for {
               top := heap.Pop(&minh).(int)
               if !removed[top] && inSet[top] {
                   v = top
                   break
               }
           }
       }
       removed[v] = true
       inSet[v] = false
       seq = append(seq, v)
       // update dependents above v
       for _, u := range dependents[v] {
           if supportCount[u] >= 2 {
               supportCount[u]--
               if supportCount[u] == 1 {
                   // find remaining supporter
                   for _, w := range supporters[u] {
                       if !removed[w] {
                           if inSet[w] {
                               inSet[w] = false
                           }
                           break
                       }
                   }
               }
           }
       }
       // update supporters below v
       for _, w := range supporters[v] {
           if removed[w] || inSet[w] {
               continue
           }
           ok := true
           for _, u := range dependents[w] {
               if supportCount[u] < 2 {
                   ok = false
                   break
               }
           }
           if ok {
               inSet[w] = true
               heap.Push(&minh, w)
               heap.Push(&maxh, w)
           }
       }
   }
   // compute result
   var res int64
   for i := 0; i < m; i++ {
       res = (res*int64(m) + int64(seq[i])) % MOD
   }
   fmt.Println(res)
}
