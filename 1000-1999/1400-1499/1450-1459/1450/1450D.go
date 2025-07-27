package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       pos := make([]int, n+1)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
           pos[a[i]] = i
       }
       // compute nearest smaller on left and right for each index
       L := make([]int, n)
       R := make([]int, n)
       // left
       stack := make([]int, 0, n)
       for i := 0; i < n; i++ {
           for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
               stack = stack[:len(stack)-1]
           }
           if len(stack) == 0 {
               L[i] = -1
           } else {
               L[i] = stack[len(stack)-1]
           }
           stack = append(stack, i)
       }
       // right
       stack = stack[:0]
       for i := n - 1; i >= 0; i-- {
           for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
               stack = stack[:len(stack)-1]
           }
           if len(stack) == 0 {
               R[i] = n
           } else {
               R[i] = stack[len(stack)-1]
           }
           stack = append(stack, i)
       }
       // span for each value x: maximum window size where x can be min
       span := make([]int, n+1)
       for x := 1; x <= n; x++ {
           idx := pos[x]
           span[x] = R[idx] - L[idx] - 1
       }
       // prefix minima of spans
       minSpan := make([]int, n+2)
       minSpan[0] = n + 1
       for i := 1; i <= n; i++ {
           minSpan[i] = min(minSpan[i-1], span[i])
       }
       // build answer for k=1..n
       res := make([]byte, n)
       for k := 1; k <= n; k++ {
           m := n - k + 1
           if minSpan[m] >= k {
               res[k-1] = '1'
           } else {
               res[k-1] = '0'
           }
       }
       fmt.Fprintln(out, string(res))
   }
}
