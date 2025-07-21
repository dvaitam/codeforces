package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// MaxHeap implements a max-heap of int64 values.
type MaxHeap []int64

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
   *h = append(*h, x.(int64))
}

func (h *MaxHeap) Pop() interface{} {
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

   var n, m, k int
   var p int64
   fmt.Fscan(reader, &n, &m, &k, &p)

   rowSum := make([]int64, n)
   colSum := make([]int64, m)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           var a int64
           fmt.Fscan(reader, &a)
           rowSum[i] += a
           colSum[j] += a
       }
   }

   // bestR[i] = max pleasure by doing i row operations
   bestR := make([]int64, k+1)
   // bestC[j] = max pleasure by doing j column operations
   bestC := make([]int64, k+1)

   // Row operations
   hR := make(MaxHeap, len(rowSum))
   copy(hR, rowSum)
   heap.Init(&hR)
   for i := 1; i <= k; i++ {
       x := heap.Pop(&hR).(int64)
       bestR[i] = bestR[i-1] + x
       heap.Push(&hR, x - int64(m)*p)
   }

   // Column operations
   hC := make(MaxHeap, len(colSum))
   copy(hC, colSum)
   heap.Init(&hC)
   for i := 1; i <= k; i++ {
       x := heap.Pop(&hC).(int64)
       bestC[i] = bestC[i-1] + x
       heap.Push(&hC, x - int64(n)*p)
   }

   // Combine
   var ans int64 = -1 << 60
   for i := 0; i <= k; i++ {
       j := k - i
       // overlap penalty: i * j * p
       val := bestR[i] + bestC[j] - int64(i)*int64(j)*p
       if val > ans {
           ans = val
       }
   }
   fmt.Fprintln(writer, ans)
}
