package main

import (
   "fmt"
   "os"
)

// primeFactors returns the map of prime -> exponent for n
func primeFactors(n int) map[int]int {
   m := make(map[int]int)
   d := 2
   for d*d <= n {
       for n%d == 0 {
           m[d]++
           n /= d
       }
       d++
   }
   if n > 1 {
       m[n]++
   }
   return m
}

func main() {
   var b, d int
   if _, err := fmt.Fscan(os.Stdin, &b, &d); err != nil {
       return
   }
   // factorize
   fb := primeFactors(b)
   // for b-1 and b+1, we just need divisibility tests
   // check 2-type: d divides some b^k
   fd := primeFactors(d)
   ok2 := true
   k := 0
   for p, ed := range fd {
       eb := fb[p]
       if eb == 0 {
           ok2 = false
           break
       }
       // need k * eb >= ed => k >= ceil(ed/eb)
       need := (ed + eb - 1) / eb
       if need > k {
           k = need
       }
   }
   if ok2 {
       fmt.Println("2-type")
       fmt.Println(k)
       return
   }
   // 3-type: d divides b-1
   if (b-1) % d == 0 {
       fmt.Println("3-type")
       return
   }
   // 11-type: d divides b+1
   if (b+1) % d == 0 {
       fmt.Println("11-type")
       return
   }
   // 6-type: all prime factors in b, b-1, or b+1
   ok6 := true
   for p := range fd {
       if fb[p] > 0 {
           continue
       }
       if (b-1) % p == 0 {
           continue
       }
       if (b+1) % p == 0 {
           continue
       }
       ok6 = false
       break
   }
   if ok6 {
       fmt.Println("6-type")
       return
   }
   // otherwise
   fmt.Println("7-type")
}
