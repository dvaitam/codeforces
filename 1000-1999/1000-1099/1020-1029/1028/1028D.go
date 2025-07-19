package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

const INF = 1000000007

// MaxHeap is a max-heap of ints
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

// MinHeap is a min-heap of ints
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

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   readInt := func() int {
       var x int
       var neg bool
       b, err := in.ReadByte()
       if err != nil {
           return 0
       }
       for (b < '0' || b > '9') && b != '-' {
           b, _ = in.ReadByte()
       }
       if b == '-' {
           neg = true
           b, _ = in.ReadByte()
       }
       for b >= '0' && b <= '9' {
           x = x*10 + int(b-'0')
           b, _ = in.ReadByte()
       }
       if neg {
           return -x
       }
       return x
   }

   readCmd := func() byte {
       for {
           b, err := in.ReadByte()
           if err != nil {
               return 0
           }
           if b == 'D' || b == 'C' {
               return b
           }
       }
   }

   n := readInt()
   ans := int64(1)
   var l int = 0
   var r int = INF
   tmp := make([]int, 0, n)
   var k1 MaxHeap
   var k2 MinHeap
   heap.Init(&k1)
   heap.Init(&k2)

   for i := 0; i < n; i++ {
       cmd := readCmd()
       if cmd == 'D' {
           x := readInt()
           tmp = append(tmp, x)
       } else {
           x := readInt()
           if x < l || x > r {
               fmt.Fprint(out, 0)
               return
           }
           flag := true
           if x == l && k1.Len() > 0 {
               heap.Pop(&k1)
               flag = false
           }
           if x == r && k2.Len() > 0 {
               heap.Pop(&k2)
               flag = false
           }
           for _, v := range tmp {
               if v == x {
                   continue
               }
               if v < x {
                   heap.Push(&k1, v)
               } else {
                   heap.Push(&k2, v)
               }
           }
           if flag {
               ans = ans * 2 % INF
           }
           if k1.Len() > 0 {
               l = k1[0]
           } else {
               l = 0
           }
           if k2.Len() > 0 {
               r = k2[0]
           } else {
               r = INF
           }
           tmp = tmp[:0]
       }
   }
   if len(tmp) == 0 {
       fmt.Fprintln(out, ans)
       return
   }
   sum := 0
   for _, v := range tmp {
       if v < l {
           heap.Push(&k1, v)
       } else if v > r {
           heap.Push(&k2, v)
       } else {
           sum++
       }
   }
   ans = ans * int64(sum+1) % INF
   fmt.Fprintln(out, ans)
}
