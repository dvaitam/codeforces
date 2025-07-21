package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b, x, y int64
   if _, err := fmt.Fscan(reader, &a, &b, &x, &y); err != nil {
       return
   }
   g := gcd(x, y)
   x /= g
   y /= g
   // find maximum k such that k*x <= a and k*y <= b
   k1 := a / x
   k2 := b / y
   k := k1
   if k2 < k {
       k = k2
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if k > 0 {
       fmt.Fprint(writer, k*x, " ", k*y)
   } else {
       fmt.Fprint(writer, "0 0")
   }
}
