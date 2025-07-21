package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func modPow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // We need C(2n-1, n) = (2n-1)! / (n! * (n-1)!)
   N := 2*n
   fac := make([]int64, N+1)
   invFac := make([]int64, N+1)
   fac[0] = 1
   for i := 1; i <= N; i++ {
       fac[i] = fac[i-1] * int64(i) % mod
   }
   invFac[N] = modPow(fac[N], mod-2)
   for i := N; i > 0; i-- {
       invFac[i-1] = invFac[i] * int64(i) % mod
   }
   // compute C(2n-1, n)
   comb := fac[2*n-1] * invFac[n] % mod * invFac[n-1] % mod
   // total = 2*comb - n
   ans := (comb*2 - int64(n)) % mod
   if ans < 0 {
       ans += mod
   }
   fmt.Fprintln(os.Stdout, ans)
}
