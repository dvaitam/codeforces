package main

import "fmt"

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   var n, m, k int64
   if _, err := fmt.Scan(&n, &m, &k); err != nil {
       return
   }
   tn, tm, tk := n, m, k
   tg := gcd(tn, tk)
   tn /= tg
   tk /= tg
   tg = gcd(tm, tk)
   tm /= tg
   tk /= tg
   if tk != 1 && tk != 2 {
       fmt.Println("NO")
       return
   }
   if tk == 2 {
       fmt.Println("YES")
       fmt.Println(0, 0)
       fmt.Println(tn, 0)
       fmt.Println(0, tm)
       return
   }
   // tk == 1
   f := false
   if tn*2 <= n {
       tn *= 2
       f = true
   }
   if !f && tm*2 <= m {
       tm *= 2
       f = true
   }
   if !f {
       fmt.Println("NO")
       return
   }
   fmt.Println("YES")
   fmt.Println(0, 0)
   fmt.Println(tn, 0)
   fmt.Println(0, tm)
}
