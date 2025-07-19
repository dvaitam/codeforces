package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, mo, FL, FR int64
   fmt.Fscan(reader, &n, &mo, &FL, &FR)
   var Exp float64
   Ll, LR := FL, FR
   var NL, NR int64
   for i := int64(2); i <= n; i++ {
       fmt.Fscan(reader, &NL, &NR)
       Exp = work(Exp, Ll, LR, NL, NR, mo)
       Ll, LR = NL, NR
   }
   Exp = work(Exp, FL, FR, NL, NR, mo)
   fmt.Printf("%.10f\n", Exp)
}

// work computes and accumulates the contribution between intervals [a,b] and [c,d]
func work(Exp float64, a, b, c, d, mo int64) float64 {
   A := float64(b/mo - (a-1)/mo)
   B := float64(d/mo - (c-1)/mo)
   C := A * B
   len1 := float64(b - a + 1)
   len2 := float64(d - c + 1)
   ans := (A*len2 + B*len1 - C) * 2000
   Exp += ans / len1 / len2
   return Exp
}
