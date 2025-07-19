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

   var n, x int
   if _, err := fmt.Fscan(in, &n, &x); err != nil {
       return
   }
   // Single element
   if n == 1 {
       fmt.Fprintln(out, "YES")
       fmt.Fprintln(out, x)
       return
   }
   // Impossible for n==2 and x==0
   if n == 2 && x == 0 {
       fmt.Fprintln(out, "NO")
       return
   }
   // Compute XOR of 0..n-3 (i < n-2)
   y := 0
   for i := 0; i < n-2; i++ {
       y ^= i
   }
   // Set high bit to avoid collisions
   const bit = 1 << 17
   y |= bit
   x |= bit
   // Adjust if equal
   m := n - 2
   isAdjust := false
   for x == y && m <= 1000000 {
       isAdjust = true
       y ^= m
       m++
   }
   // Build result
   res := make([]int, 0, n)
   if !isAdjust {
       for i := 0; i < n-2; i++ {
           res = append(res, i)
       }
   } else {
       for i := 1; i < n-2; i++ {
           res = append(res, i)
       }
       // append the adjusted value
       res = append(res, m-1)
   }
   // append final two values
   res = append(res, y, x)

   // Output
   fmt.Fprintln(out, "YES")
   for i, v := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
