package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// modExp computes base^exp % mod
func modExp(base, exp int64) int64 {
   res := int64(1)
   base %= mod
   for exp > 0 {
       if exp&1 == 1 {
           res = (res * base) % mod
       }
       base = (base * base) % mod
       exp >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var x int64
   if _, err := fmt.Fscan(in, &n, &x); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // sum of all exponents
   var sumA int64
   for _, v := range a {
       sumA += v
   }
   // maximum exponent
   maxA := a[n-1]
   // count differences d = maxA - a[i]
   cnt := make(map[int64]int64, n)
   for _, v := range a {
       d := maxA - v
       cnt[d]++
   }
   // compute v_x(T) where T = sum x^{d_i}
   var carry int64
   var vT int64
   for j := int64(0); ; j++ {
       total := carry
       if c, ok := cnt[j]; ok {
           total += c
       }
       if total % x != 0 {
           vT = j
           break
       }
       carry = total / x
   }
   // exponent of gcd = min(sumA, (sumA - maxA) + vT)
   exp := (sumA - maxA) + vT
   if exp > sumA {
       exp = sumA
   }
   // compute x^exp mod mod
   ans := modExp(x, exp)
   fmt.Println(ans)
}
