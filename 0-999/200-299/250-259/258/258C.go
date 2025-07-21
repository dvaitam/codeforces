package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func powmod(a, e int64) int64 {
   res := int64(1)
   a %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   maxa := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       if a[i] > maxa {
           maxa = a[i]
       }
   }
   cnt := make([]int, maxa+2)
   for _, v := range a {
       cnt[v]++
   }
   // s[x] = number of a[i] >= x
   s := make([]int, maxa+3)
   s[maxa+1] = 0
   for x := maxa; x >= 1; x-- {
       s[x] = s[x+1] + cnt[x]
   }
   // divisors
   divs := make([][]int, maxa+1)
   for d := 1; d <= maxa; d++ {
       for m := d; m <= maxa; m += d {
           divs[m] = append(divs[m], d)
       }
   }
   // find max number of divisors
   maxk := 0
   for m := 1; m <= maxa; m++ {
       if k := len(divs[m]); k > maxk {
           maxk = k
       }
   }
   inv := make([]int64, maxk+1)
   for k := 1; k <= maxk; k++ {
       inv[k] = powmod(int64(k), MOD-2)
   }
   var ans int64 = 0
   // for each M
   for M := 1; M <= maxa; M++ {
       D := divs[M]
       k := len(D)
       // compute H[M]
       var H int64 = 1
       // intervals between divisors
       for j := 1; j < k; j++ {
           lower := D[j-1]
           upper := D[j]
           cntj := int64(s[lower] - s[upper])
           if cntj > 0 {
               H = H * powmod(int64(j), cntj) % MOD
           }
       }
       // positions with a[i] >= M
       sm := int64(s[M])
       if sm > 0 {
           H = H * powmod(int64(k), sm) % MOD
       }
       // subtract sequences without any b[i] == M
       // ratio = (k-1)/k mod MOD
       var Fm int64
       if k > 0 {
           ratio := int64(k-1) * inv[k] % MOD
           rpow := powmod(ratio, sm)
           Fm = H * ((1 - rpow + MOD) % MOD) % MOD
       }
       ans = (ans + Fm) % MOD
   }
   fmt.Println(ans)
}
