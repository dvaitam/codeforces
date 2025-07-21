package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var l, r int64
   if _, err := fmt.Fscan(reader, &l, &r); err != nil {
       return
   }

   // Precompute powers of 10 up to 10^10
   pow10 := [...]int64{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 10000000000}

   var ans int64
   // Iterate over possible digit lengths L
   for L := 1; L < len(pow10); L++ {
       // interval for numbers with exactly L digits
       lo := l
       if lo < pow10[L-1] {
           lo = pow10[L-1]
       }
       hi := r
       if hi > pow10[L]-1 {
           hi = pow10[L] - 1
       }
       if lo > hi {
           continue
       }
       M := pow10[L] - 1
       // candidate points: ends and parabola peak
       c1 := M / 2
       c2 := c1 + 1
       candidates := []int64{lo, hi, c1, c2}
       for _, n := range candidates {
           if n < lo || n > hi {
               continue
           }
           prod := n * (M - n)
           if prod > ans {
               ans = prod
           }
       }
   }
   fmt.Fprint(writer, ans)
}
