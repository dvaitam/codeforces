package main

import (
   "fmt"
)

// sumDigits returns the sum of decimal digits of x.
func sumDigits(x uint64) uint64 {
   var s uint64
   for x > 0 {
       s += x % 10
       x /= 10
   }
   return s
}

// integerSqrt returns floor(sqrt(n)) for uint64 n.
func integerSqrt(n uint64) uint64 {
   var low, high uint64 = 0, 2000000000
   for low < high {
       mid := (low + high + 1) >> 1
       if mid*mid <= n {
           low = mid
       } else {
           high = mid - 1
       }
   }
   return low
}

func main() {
   var n uint64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var ans uint64
   found := false
   // s(x) is at most 9*19, but we check upto 200
   for s := uint64(1); s <= 200; s++ {
       // x^2 + s*x - n = 0 => x = (sqrt(s^2+4n) - s) / 2
       D := s*s + 4*n
       r := integerSqrt(D)
       if r*r != D {
           continue
       }
       if r < s || (r-s)&1 != 0 {
           continue
       }
       x := (r - s) >> 1
       if x == 0 {
           continue
       }
       if sumDigits(x) != s {
           continue
       }
       if !found || x < ans {
           ans = x
           found = true
       }
   }
   if !found {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}
