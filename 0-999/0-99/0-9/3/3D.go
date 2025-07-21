package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// item holds the cost difference and position for a '?' placeholder
type item struct {
   diff int // cost to change from '(' to ')'
   idx  int // position in result
}

// minHeap implements a min-heap of items based on diff
type minHeap []item

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].diff < h[j].diff }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[:n-1]
   return it
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read pattern
   var pattern string
   if _, err := fmt.Fscanln(reader, &pattern); err != nil {
       fmt.Println(-1)
       return
   }
   // count '?' to know how many cost pairs
   qCount := 0
   for _, ch := range pattern {
       if ch == '?' {
           qCount++
       }
   }
   // read cost pairs
   costs := make([][2]int, qCount)
   for i := 0; i < qCount; i++ {
       if _, err := fmt.Fscan(reader, &costs[i][0], &costs[i][1]); err != nil {
           fmt.Println(-1)
           return
       }
   }
   // prepare result as runes
   res := []rune(pattern)
   h := &minHeap{}
   heap.Init(h)
   balance := 0
   totalCost := 0
   qi := 0
   // initial assignment: treat all '?' as '(' and push diff
   for i, ch := range res {
       switch ch {
       case '(':
           balance++
       case ')':
           balance--
       case '?':
           a := costs[qi][0]
           b := costs[qi][1]
           totalCost += a
           res[i] = '('
           // diff is extra cost to flip to ')'
           heap.Push(h, item{diff: b - a, idx: i})
           balance++
           qi++
       }
       if balance < 0 {
           // need to flip a previous '?' to ')'
           if h.Len() == 0 {
               fmt.Println(-1)
               return
           }
           it := heap.Pop(h).(item)
           res[it.idx] = ')'
           totalCost += it.diff
           balance += 2
       }
   }
   // close extra '(' by flipping
   for balance > 0 {
       if h.Len() == 0 {
           fmt.Println(-1)
           return
       }
       it := heap.Pop(h).(item)
       res[it.idx] = ')'
       totalCost += it.diff
       balance -= 2
   }
   if balance != 0 {
       fmt.Println(-1)
       return
   }
   // output result
   fmt.Println(totalCost)
   fmt.Println(string(res))
}
