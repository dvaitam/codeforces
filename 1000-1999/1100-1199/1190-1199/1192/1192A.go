package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, t int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   fmt.Fscan(in, &t)
   coords := make([][2]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &coords[i][0], &coords[i][1])
   }
   if n == 1 {
       fmt.Fprintln(out, "YES")
       fmt.Fprintln(out, 1)
       return
   }
   // map coordinates to index (1-based)
   idx := make(map[int64]int, n)
   for i, c := range coords {
       key := (int64(c[0]) + 1000000000) << 32 | int64(c[1]+1000000000)
       idx[key] = i
   }
   // neighbor lists
   nbr4 := make([][]int, n)
   nbr8 := make([][]int, n)
   cnt4 := make([]int, n)
   cnt8 := make([]int, n)
   for i, c := range coords {
       x, y := c[0], c[1]
       for dx := -1; dx <= 1; dx++ {
           for dy := -1; dy <= 1; dy++ {
               if dx == 0 && dy == 0 {
                   continue
               }
               key := (int64(x+dx) + 1000000000) << 32 | int64(y+dy+1000000000)
               if j, ok := idx[key]; ok {
                   cnt8[i]++
                   nbr8[i] = append(nbr8[i], j)
                   if dx*dx+dy*dy == 1 {
                       cnt4[i]++
                       nbr4[i] = append(nbr4[i], j)
                   }
               }
           }
       }
   }
   removed := make([]bool, n)
   inCand := make([]bool, n)
   s := &intHeap{cmpMax: true}
   heap.Init(s)
   // initial candidates
   for i := 0; i < n; i++ {
       if cnt4[i] < 4 && (cnt8[i] > 0) {
           heap.Push(s, i)
           inCand[i] = true
       }
   }
   var remOrder []int
   remOrder = make([]int, 0, n)
   remS := n
   for len(remOrder) < n {
       if s.Len() == 0 {
           break
       }
       i := heap.Pop(s).(int)
       if removed[i] {
           continue
       }
       // check current valid
       if !(cnt4[i] < 4 && (cnt8[i] > 0 || remS == 1)) {
           continue
       }
       // remove i
       removed[i] = true
       remOrder = append(remOrder, i)
       remS--
       // update neighbors
       for _, j := range nbr4[i] {
           if removed[j] {
               continue
           }
           cnt4[j]--
       }
       for _, j := range nbr8[i] {
           if removed[j] {
               continue
           }
           cnt8[j]--
       }
       // push neighbors that may become candidates
       if remS > 0 {
           for _, j := range append(nbr4[i], nbr8[i]...) {
               if removed[j] || inCand[j] {
                   continue
               }
               if cnt4[j] < 4 && (cnt8[j] > 0 || remS == 1) {
                   heap.Push(s, j)
                   inCand[j] = true
               }
           }
       }
   }
   if len(remOrder) < n {
       fmt.Fprintln(out, "NO")
       return
   }
   // build answer: reverse removal order and convert to 1-based indices
   sOrder := make([]int, n)
   for i, v := range remOrder {
       sOrder[n-1-i] = v + 1
   }
   fmt.Fprintln(out, "YES")
   for _, v := range sOrder {
       fmt.Fprintln(out, v)
   }
}

// max-heap of ints
type intHeap struct {
   data   []int
   cmpMax bool
}

func (h intHeap) Len() int           { return len(h.data) }
func (h intHeap) Less(i, j int) bool {
   if h.cmpMax {
       return h.data[i] > h.data[j]
   }
   return h.data[i] < h.data[j]
}
func (h intHeap) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }
func (h *intHeap) Push(x interface{}) { h.data = append(h.data, x.(int)) }
func (h *intHeap) Pop() interface{} {
   old := h.data
   n := len(old)
   x := old[n-1]
   h.data = old[:n-1]
   return x
}
