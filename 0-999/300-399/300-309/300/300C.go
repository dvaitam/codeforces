package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

// modPow computes x^p mod mod
func modPow(x, p int64) int64 {
   res := int64(1)
   x %= mod
   for p > 0 {
       if p&1 == 1 {
           res = res * x % mod
       }
       x = x * x % mod
       p >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var a, b, n int
   if _, err := fmt.Fscan(in, &a, &b, &n); err != nil {
       return
   }

   // Precompute factorials and inverse factorials
   fact := make([]int64, n+1)
   invFact := make([]int64, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invFact[n] = modPow(fact[n], mod-2)
   for i := n; i > 0; i-- {
       invFact[i-1] = invFact[i] * int64(i) % mod
   }

   var res int64
   // check for each count of digit a
   for k := 0; k <= n; k++ {
       sum := k*a + (n-k)*b
       // check if sum's digits are only a or b
       x := sum
       good := true
       for x > 0 {
           d := x % 10
           if d != a && d != b {
               good = false
               break
           }
           x /= 10
       }
       if !good {
           continue
       }
       // add C(n, k)
       comb := fact[n] * invFact[k] % mod * invFact[n-k] % mod
       res = (res + comb) % mod
   }
   fmt.Fprint(out, res)
