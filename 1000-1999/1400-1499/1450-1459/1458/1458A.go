package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd computes the greatest common divisor of a and b.
func gcd(a, b uint64) uint64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// absDiff returns the absolute difference of a and b.
func absDiff(a, b uint64) uint64 {
   if a > b {
       return a - b
   }
   return b - a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([]uint64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   b := make([]uint64, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(in, &b[j])
   }

   base := a[0]
   var g uint64
   for i := 1; i < n; i++ {
       diff := absDiff(a[i], base)
       g = gcd(g, diff)
   }

   for j := 0; j < m; j++ {
       ans := gcd(g, base + b[j])
       fmt.Fprintf(out, "%d", ans)
       if j+1 < m {
           out.WriteByte(' ')
       }
   }
   fmt.Fprintln(out)
}
