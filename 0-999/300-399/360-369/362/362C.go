package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   p := make([]int, n)
   for i := range p {
       fmt.Fscan(in, &p[i])
   }
   // C[i][v] = number of elements p[0..i-1] with value < v
   C := make([][]uint16, n+1)
   for i := 0; i <= n; i++ {
       C[i] = make([]uint16, n+1)
   }
   for i := 1; i <= n; i++ {
       pi := p[i-1]
       prev := C[i-1]
       cur := C[i]
       for v := 1; v <= n; v++ {
           cur[v] = prev[v]
           if pi < v {
               cur[v]++
           }
       }
   }
   // initial inversion count
   var inv0 int64
   for i := 0; i < n; i++ {
       inv0 += int64(i) - int64(C[i][p[i]])
   }
   // find minimal delta and count
   minDelta := int64(1<<60)
   var count int64
   for i := 0; i < n; i++ {
       xi := p[i]
       for j := i + 1; j < n; j++ {
           yj := p[j]
           low, high := xi, yj
           if low > high {
               low, high = high, low
           }
           diff := high - low
           // count in p[0..i-1]
           cntLeft := int64(C[i][high] - C[i][low+1])
           // count in p[i+1..j-1]
           cntMid := int64((C[j][high] - C[j][low+1]) - (C[i+1][high] - C[i+1][low+1]))
           // K = diff - 2*cntLeft + cntMid
           K := int64(diff) - 2*cntLeft + cntMid
           var delta int64
           if xi < yj {
               delta = K
           } else {
               delta = -K
           }
           if delta < minDelta {
               minDelta = delta
               count = 1
           } else if delta == minDelta {
               count++
           }
       }
   }
   fmt.Fprintln(out, inv0+minDelta, count)
}
