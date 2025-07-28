package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = int64(1e9 + 7)

func modPow(a, e, m int) int {
   res := 1 % m
   a %= m
   for e > 0 {
       if e&1 == 1 {
           res = (res * a) % m
       }
       a = (a * a) % m
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   fmt.Fscan(in, &n, &m, &k)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // total elements N = n * m
   if k == 1 {
       // all segments of length < N plus one full cycle
       N := int64(n) * int64(m) % MOD
       ans := (N * ((N - 1 + MOD) % MOD)) % MOD
       ans = (ans + 1) % MOD
       fmt.Println(ans)
       return
   }
   // prefix sums modulo k for one period, positions 0..n-1
   presum := make([]int, n)
   // presum[0] = 0
   for i := 1; i < n; i++ {
       presum[i] = (presum[i-1] + a[i-1]) % k
   }
   // total sum of a mod k
   var S int
   for i := 0; i < n; i++ {
       S = (S + a[i]) % k
   }
   // count of prefix residues in one block
   cntPresum := make([]int64, k)
   for i := 0; i < n; i++ {
       cntPresum[presum[i]]++
   }
   // handle S == 0 mod k separately
   fullCycles := m / k
   rem := m % k
   countRes := make([]int64, k)
   if S % k == 0 {
       // prefix residues same every block
       for r := 0; r < k; r++ {
           countRes[r] = cntPresum[r] * int64(m) % MOD
       }
   } else {
       // S invertible mod k
       // build d array: d[i] = cntPresum[(i*S) mod k]
       d := make([]int64, k)
       for i := 0; i < k; i++ {
           idx := (i * S) % k
           d[i] = cntPresum[idx]
       }
       // prefix sums of d
       pref := make([]int64, k+1)
       for i := 0; i < k; i++ {
           pref[i+1] = pref[i] + d[i]
       }
       totalD := pref[k]
       // compute invS modulo k
       invS := modPow(S, k-2, k)
       base := int64(fullCycles) * int64(n) % MOD
       for r := 0; r < k; r++ {
           // u = r * invS mod k
           u := (r * invS) % k
           var extra int64
           if rem > 0 {
               // window [u-rem+1 .. u] in cyclic d
               l := u - rem + 1
               if l >= 0 {
                   extra = pref[u+1] - pref[l]
               } else {
                   // wrap around
                   extra = (pref[u+1] + totalD - pref[(l+k)])
               }
           }
           countRes[r] = (base + extra) % MOD
       }
   }
   // count segments of length < N: unordered P pairs (i,j) and (j,i) count both arcs
   var ans int64
   for r := 0; r < k; r++ {
       c := countRes[r]
       // each unordered P pair yields two segments (i->j and j->i), so count c*(c-1)
       ans = (ans + c*((c-1+MOD)%MOD)%MOD) % MOD
   }
   // add full cycle segment if total sum of b divisible by k
   if int64(m)*int64(S)%int64(k) == 0 {
       ans = (ans + 1) % MOD
   }
   fmt.Println(ans)
}
