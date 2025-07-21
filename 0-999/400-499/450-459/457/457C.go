package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// IntHeap is a min-heap of ints
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) {
   *h = append(*h, x.(int))
}
func (h *IntHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   // Map candidate to list of bribe costs
   costsMap := make(map[int][]int)
   c0 := 0
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       if a == 0 {
           c0++
       } else {
           costsMap[a] = append(costsMap[a], b)
       }
   }
   // Bucket costs by number of supporters
   maxCount := 0
   buckets := make([][]int, n+1)
   for _, costs := range costsMap {
       cnt := len(costs)
       if cnt > maxCount {
           maxCount = cnt
       }
       sort.Ints(costs)
       // add all costs of this candidate to bucket[cnt]
       buckets[cnt] = append(buckets[cnt], costs...)
   }
   // Min-heap of available bribe costs
   h := &IntHeap{}
   heap.Init(h)
   bribes := 0
   var costSum int64
   ans := int64(1<<63 - 1)
   // Iterate threshold of maximum opponent votes
   for d := maxCount; d >= 0; d-- {
       // add supporters of candidates with initial votes == d
       for _, cost := range buckets[d] {
           heap.Push(h, cost)
       }
       // bribe until we have more votes than d
       // need c0 + bribes > d => bribes >= d - c0 + 1 - bribes
       needed := d - (c0 + bribes) + 1
       for needed > 0 && h.Len() > 0 {
           minCost := heap.Pop(h).(int)
           costSum += int64(minCost)
           bribes++
           needed--
       }
       if c0+bribes > d && costSum < ans {
           ans = costSum
       }
   }
   fmt.Fprint(writer, ans)
