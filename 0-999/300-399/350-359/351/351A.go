package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   total := 2 * n
   k := 0
   sumRem := 0 // sum of remainders in thousandths
   for i := 0; i < total; i++ {
       var s string
       fmt.Fscan(reader, &s)
       // s is like ddd.ddd
       // find decimal point
       var intPart, fracPart int
       // parse manually to avoid float errors
       for j := 0; j < len(s); j++ {
           if s[j] == '.' {
               // integer part
               fmt.Sscanf(s[:j], "%d", &intPart)
               // fractional part exact 3 digits
               fmt.Sscanf(s[j+1:], "%d", &fracPart)
               break
           }
       }
       rem := fracPart
       if rem != 0 {
           k++
           sumRem += rem
       }
   }
   // we need to choose c ceils among fractional elements
   // c in [L, R]
   L := k - n
   if L < 0 {
       L = 0
   }
   R := k
   if R > n {
       R = n
   }
   // best c rounds sumRem/1000
   // c0 = (sumRem + 500) / 1000
   c0 := (sumRem + 500) / 1000
   // clamp c0 to [L, R]
   c := c0
   if c < L {
       c = L
   } else if c > R {
       c = R
   }
   // compute diff = |1000*c - sumRem|
   diff := c*1000 - sumRem
   if diff < 0 {
       diff = -diff
   }
   // print diff/1000 with 3 decimals
   fmt.Printf("%d.%03d\n", diff/1000, diff%1000)
}
