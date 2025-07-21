package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// an IntMaxHeap is a max-heap of ints.
type IntMaxHeap []int

func (h IntMaxHeap) Len() int           { return len(h) }
func (h IntMaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntMaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntMaxHeap) Push(x interface{}) {
   *h = append(*h, x.(int))
}
func (h *IntMaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var m, k int
   if _, err := fmt.Fscan(reader, &m, &k); err != nil {
       return
   }
   d := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &d[i])
   }
   s := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &s[i])
   }
   // Initialize
   time := 0
   fuel := 0
   h := &IntMaxHeap{}
   heap.Init(h)
   // initial at city 1
   if m > 0 {
       fuel = s[0]
       heap.Push(h, s[0])
   }
   // Traverse roads
   for i := 0; i < m; i++ {
       need := d[i]
       // wait as needed
       for fuel < need {
           if h.Len() == 0 {
               // no supply, should not happen
               break
           }
           best := (*h)[0]
           fuel += best
           time += k
       }
       // travel
       fuel -= need
       time += need
       // at city i+2, supply if exists
       if i+1 < m {
           heap.Push(h, s[i+1])
           fuel += s[i+1]
       }
   }
   // output total time
   fmt.Println(time)
}
