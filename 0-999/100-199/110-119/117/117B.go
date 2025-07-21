package main

import (
   "fmt"
)

// extended gcd: returns gcd(a,b) and x,y such that ax+by=gcd
func extGCD(a, b int64) (int64, int64, int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := extGCD(b, a%b)
   x := y1
   y := x1 - (a/b)*y1
   return g, x, y
}

func main() {
   var a, b, mod int64
   if _, err := fmt.Scanf("%d %d %d", &a, &b, &mod); err != nil {
       return
   }
   // Quick lose if b >= mod-1
   if b >= mod-1 {
       fmt.Println(2)
       return
   }
   // compute P = 10^9 mod mod
   const base = int64(1000000000)
   P := base % mod
   // compute G = gcd(P, mod)
   g, _, _ := extGCD(P, mod)
   G := g
   // reduced modulus
   m := mod / G
   // reduced multiplier
   Pm := P / G
   // inverse of Pm modulo m
   _, invX, _ := extGCD(Pm, m)
   invP := invX % m
   if invP < 0 {
       invP += m
   }
   // range threshold
   R := mod - b
   // search for minimal s1
   maxI := (R - 1) / G
   var ans int64 = -1
   for i := int64(1); i <= maxI; i++ {
       // residue y = i*G, want s0 = (y/G)*invP mod m == i*invP mod m
       s0 := (i * invP) % m
       if s0 <= a {
           if s0 > 0 { // skip r1==0 case
               if ans == -1 || s0 < ans {
                   ans = s0
                   if ans == 1 {
                       break
                   }
               }
           }
       }
   }
   if ans == -1 {
       fmt.Println(2)
   } else {
       // print winning move padded to 9 digits
       fmt.Printf("1 %09d\n", ans)
   }
}
