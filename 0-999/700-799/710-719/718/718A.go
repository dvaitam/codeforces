package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// IntHeap is a min-heap of ints.
type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
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

   var n int
   var t int64
   fmt.Fscan(in, &n, &t)
   var s string
   fmt.Fscan(in, &s)
   // split into integer and fractional parts
   dot := 0
   for i := range s {
       if s[i] == '.' {
           dot = i
           break
       }
   }
   intStr := s[:dot]
   fracStr := s[dot+1:]
   // integer part digits
   intPart := make([]int, len(intStr))
   for i := range intStr {
       intPart[i] = int(intStr[i] - '0')
   }
   // fractional part digits
   fracPart := make([]int, len(fracStr))
   for i := range fracStr {
       fracPart[i] = int(fracStr[i] - '0')
   }
   fracLen := len(fracPart)
   // init heap of positions i with fracPart[i] >= 5
   h := &IntHeap{}
   heap.Init(h)
   for i, v := range fracPart {
       if v >= 5 {
           heap.Push(h, i)
       }
   }
   // perform rounds
   for t > 0 && h.Len() > 0 {
       i := heap.Pop(h).(int)
       if i >= fracLen {
           continue
       }
       // round at position i-1 (fraction index)
       t--
       rp := i - 1
       if rp < 0 {
           // round to integer
           carry := 1
           for k := len(intPart) - 1; k >= 0 && carry > 0; k-- {
               v := intPart[k] + carry
               intPart[k] = v % 10
               carry = v / 10
           }
           if carry > 0 {
               intPart = append([]int{carry}, intPart...)
           }
           fracLen = 0
           break
       }
       // carry in fractional part
       carry := 1
       p := rp
       for k := p; k >= 0 && carry > 0; k-- {
           v := fracPart[k] + carry
           fracPart[k] = v % 10
           carry = v / 10
           p = k
       }
       if carry > 0 {
           // propagate to integer
           carryInt := 1
           for k := len(intPart) - 1; k >= 0 && carryInt > 0; k-- {
               v := intPart[k] + carryInt
               intPart[k] = v % 10
               carryInt = v / 10
           }
           if carryInt > 0 {
               intPart = append([]int{carryInt}, intPart...)
           }
           fracLen = 0
           break
       }
       // truncate fractional to p+1
       fracLen = p + 1
       // if new digit at p is >=5, schedule next
       if fracPart[p] >= 5 {
           heap.Push(h, p)
       }
   }
   // remove trailing zeros in fractional part
   for fracLen > 0 && fracPart[fracLen-1] == 0 {
       fracLen--
   }
   // output result
   for _, d := range intPart {
       out.WriteByte(byte('0' + d))
   }
   if fracLen > 0 {
       out.WriteByte('.')
       for i := 0; i < fracLen; i++ {
           out.WriteByte(byte('0' + fracPart[i]))
       }
   }
}
