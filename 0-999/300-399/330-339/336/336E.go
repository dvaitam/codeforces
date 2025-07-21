package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modPow(a, e int64) int64 {
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
   var n, k int64
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   // Special case n == 0: no triangles
   if n == 0 {
       if k == 0 {
           fmt.Println(1)
       } else {
           fmt.Println(0)
       }
       return
   }
   // cannot choose more than 4 non-overlapping triangles (one per quadrant chain)
   if k > 4 {
       fmt.Println(0)
       return
   }
   // number of ways to pick a sequence of k distinct quadrant chains: P(4, k)
   perm4 := int64(1)
   for i := int64(0); i < k; i++ {
       perm4 = perm4 * (4 - i) % MOD
   }
   // each chain has (n+1) triangles
   pw := modPow(n+1, k)
   ans := perm4 * pw % MOD
   fmt.Println(ans)
}
