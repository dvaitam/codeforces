package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math"
   "os"
)

// Item represents an adjacent boy-girl pair
type Item struct {
   diff int
   pos  int // left index of pair
   // heap index
   idx int
}

// Priority queue of Items
type PQ []*Item

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
   if pq[i].diff != pq[j].diff {
       return pq[i].diff < pq[j].diff
   }
   return pq[i].pos < pq[j].pos
}
func (pq PQ) Swap(i, j int) {
   pq[i], pq[j] = pq[j], pq[i]
   pq[i].idx = i
   pq[j].idx = j
}
func (pq *PQ) Push(x interface{}) {
   it := x.(*Item)
   it.idx = len(*pq)
   *pq = append(*pq, it)
}
func (pq *PQ) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   it.idx = -1
   *pq = old[0 : n-1]
   return it
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   s := make([]byte, n)
   fmt.Fscan(in, &s)
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // linked list
   L := make([]int, n+2)
   R := make([]int, n+2)
   alive := make([]bool, n+2)
   for i := 1; i <= n; i++ {
       L[i] = i - 1
       R[i] = i + 1
       alive[i] = true
   }
   R[0] = 1
   L[0] = 0
   R[n] = 0

   pq := &PQ{}
   heap.Init(pq)
   // push initial candidate pairs
   for i := 1; i < n; i++ {
       if s[i-1] != s[i] {
           diff := int(math.Abs(float64(a[i] - a[i+1])))
           heap.Push(pq, &Item{diff: diff, pos: i})
       }
   }
   var res [][2]int

   for pq.Len() > 0 {
       it := heap.Pop(pq).(*Item)
       i := it.pos
       // check validity
       if !alive[i] {
           continue
       }
       j := R[i]
       if j == 0 || !alive[j] {
           continue
       }
       // gender must differ
       if s[i-1] == s[j-1] {
           continue
       }
       // record pair
       if i < j {
           res = append(res, [2]int{i, j})
       } else {
           res = append(res, [2]int{j, i})
       }
       // remove i and j
       alive[i], alive[j] = false, false
       li := L[i]
       rj := R[j]
       if li >= 1 {
           R[li] = rj
       }
       if rj >= 1 {
           L[rj] = li
       }
       // new adjacency
       if li >= 1 && rj >= 1 && alive[li] && alive[rj] && s[li-1] != s[rj-1] {
           diff := int(math.Abs(float64(a[li] - a[rj])))
           heap.Push(pq, &Item{diff: diff, pos: li})
       }
   }
   // output
   fmt.Fprintln(out, len(res))
   for _, p := range res {
       fmt.Fprintln(out, p[0], p[1])
   }
}
