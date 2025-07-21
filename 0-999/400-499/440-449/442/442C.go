package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Item represents a candidate removal with its score and index
type Item struct {
   score int
   idx   int
}

// MaxHeap implements a max-heap of Items
type MaxHeap []Item

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].score > h[j].score }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
   *h = append(*h, x.(Item))
}

func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[0 : n-1]
   return it
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   if n < 3 {
       fmt.Println(0)
       return
   }
   // doubly linked list neighbors
   L := make([]int, n)
   R := make([]int, n)
   removed := make([]bool, n)
   for i := 0; i < n; i++ {
       L[i] = i - 1
       R[i] = i + 1
   }
   R[n-1] = -1

   // helper to get current score for index i
   getScore := func(i int) int {
       if i < 0 || i >= n || removed[i] {
           return -1
       }
       l := L[i]
       r := R[i]
       if l < 0 || r < 0 {
           return -1
       }
       // score is min of neighbor values
       if a[l] < a[r] {
           return a[l]
       }
       return a[r]
   }

   // build initial heap of internal nodes
   h := &MaxHeap{}
   heap.Init(h)
   for i := 1; i < n-1; i++ {
       sc := getScore(i)
       if sc >= 0 {
           heap.Push(h, Item{score: sc, idx: i})
       }
   }

   var total int64
   // process removals
   for h.Len() > 0 {
       it := heap.Pop(h).(Item)
       i := it.idx
       if removed[i] {
           continue
       }
       sc := getScore(i)
       if sc != it.score {
           continue
       }
       // perform removal
       total += int64(sc)
       removed[i] = true
       l := L[i]
       r := R[i]
       // link neighbors
       if l >= 0 {
           R[l] = r
       }
       if r >= 0 {
           L[r] = l
       }
       // push updated neighbors if they become internal
       if sc2 := getScore(l); sc2 >= 0 {
           heap.Push(h, Item{score: sc2, idx: l})
       }
       if sc3 := getScore(r); sc3 >= 0 {
           heap.Push(h, Item{score: sc3, idx: r})
       }
   }
   fmt.Println(total)
}
