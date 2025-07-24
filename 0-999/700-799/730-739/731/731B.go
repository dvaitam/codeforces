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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   coupons := 0
   for i := 0; i < n-1; i++ {
       // use coupons from previous day
       if coupons > a[i] {
           coupons = a[i]
       }
       rem := a[i] - coupons
       // decide how many new coupons to buy (max for tomorrow)
       maxC := a[i+1]
       newC := rem
       if newC > maxC {
           newC = maxC
       }
       // ensure remaining pizzas covered by pair discounts is even
       if (rem-newC)%2 != 0 {
           if newC > 0 {
               newC--
           } else {
               fmt.Fprintln(writer, "NO")
               return
           }
       }
       // after this, rem-newC is even >=0
       coupons = newC
   }
   // last day: only use remaining coupons and pair discounts
   last := a[n-1]
   if coupons > last {
       coupons = last
   }
   rem := last - coupons
   if rem%2 != 0 {
       fmt.Fprintln(writer, "NO")
   } else {
       fmt.Fprintln(writer, "YES")
   }
}
