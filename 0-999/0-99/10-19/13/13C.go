package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// MinHeap implements a min-heap of ints.
type MinHeap []int
func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

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

// Block holds values for isotonic regression block.
type Block struct {
   low     *MaxHeap // lower half (including median)
   high    *MinHeap // upper half
   sumLow  int64    // sum of values in low
   sumHigh int64    // sum of values in high
}

// NewBlock creates a block with one element x.
func NewBlock(x int) *Block {
   low := &MaxHeap{}
   high := &MinHeap{}
   heap.Init(low)
   heap.Init(high)
   b := &Block{low: low, high: high}
   b.add(x)
   return b
}

// add inserts a value into the block, rebalancing heaps.
func (b *Block) add(x int) {
   // push into appropriate heap
   if b.low.Len() == 0 || x <= (*b.low)[0] {
       heap.Push(b.low, x)
       b.sumLow += int64(x)
   } else {
       heap.Push(b.high, x)
       b.sumHigh += int64(x)
   }
   // rebalance: keep low.Len() >= high.Len() and difference <= 1
   if b.low.Len() < b.high.Len() {
       y := heap.Pop(b.high).(int)
       b.sumHigh -= int64(y)
       heap.Push(b.low, y)
       b.sumLow += int64(y)
   } else if b.low.Len() > b.high.Len()+1 {
       y := heap.Pop(b.low).(int)
       b.sumLow -= int64(y)
       heap.Push(b.high, y)
       b.sumHigh += int64(y)
   }
}

// median returns current median of the block.
func (b *Block) median() int {
   return (*b.low)[0]
}

// cost computes sum of absolute deviations to median.
func (b *Block) cost() int64 {
   m := int64(b.median())
   return m*int64(b.low.Len()) - b.sumLow + b.sumHigh - m*int64(b.high.Len())
}

// merge combines another block into this one.
func (b *Block) merge(other *Block) {
   for _, x := range *other.low {
       b.add(x)
   }
   for _, x := range *other.high {
       b.add(x)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   blocks := make([]*Block, 0, n)
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       b := NewBlock(a)
       blocks = append(blocks, b)
       // enforce non-decreasing medians
       for len(blocks) >= 2 {
           last := blocks[len(blocks)-1]
           prev := blocks[len(blocks)-2]
           if prev.median() <= last.median() {
               break
           }
           // merge prev and last
           blocks = blocks[:len(blocks)-2]
           prev.merge(last)
           blocks = append(blocks, prev)
       }
   }
   var ans int64
   for _, b := range blocks {
       ans += b.cost()
   }
   fmt.Println(ans)
}
