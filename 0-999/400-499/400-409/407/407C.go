package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   const MOD = 1000000007
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   const K = 100
   maxN := n + K + 5
   // factorials and inverse factorials
   fac := make([]int, maxN)
   ifac := make([]int, maxN)
   fac[0] = 1
   for i := 1; i < maxN; i++ {
       fac[i] = int(int64(fac[i-1]) * int64(i) % MOD)
   }
   ifac[maxN-1] = modInv(fac[maxN-1], MOD)
   for i := maxN - 1; i > 0; i-- {
       ifac[i-1] = int(int64(ifac[i]) * int64(i) % MOD)
   }
   comb := func(n, k int) int {
       if k < 0 || k > n {
           return 0
       }
       return int(int64(fac[n]) * int64(ifac[k]) % MOD * int64(ifac[n-k]) % MOD)
   }
   // add[t][i] accumulates t-th order contributions at pos i
   // add[t][i] accumulates t-th order contributions at pos i (0..n)
   add := make([][]int, K+2)
   for t := range add {
       add[t] = make([]int, n+2)
   }
   for qi := 0; qi < m; qi++ {
       var l, r, k int
       fmt.Fscan(reader, &l, &r, &k)
       // convert to zero-based
       l--
       r--
       // start sequence of order k at l
       add[k][l] = (add[k][l] + 1) % MOD
       // subtract contributions after r
       for i := 0; i <= k; i++ {
           x := k - i
           // total length from l to r is (r-l)
           tval := comb(r-l+x, x)
           idx := r + 1
           add[i][idx] = (add[i][idx] - tval + MOD) % MOD
       }
   }
   // accumulate contributions
   // accumulate contributions: first prefix sums over i, then add next order
   for t := K; t >= 0; t-- {
       for i := 0; i <= n; i++ {
           if i > 0 {
               add[t][i] = (add[t][i] + add[t][i-1]) % MOD
           }
           add[t][i] = (add[t][i] + add[t+1][i]) % MOD
       }
   }
   // output final array
   for i := 0; i < n; i++ {
       res := (a[i] + add[0][i]) % MOD
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(res))
   }
   writer.WriteByte('\n')
}

// modInv returns modular inverse of a mod m, assuming m is prime
func modInv(a, m int) int {
   return modPow(a, m-2, m)
}

// modPow computes a^e mod m
func modPow(a, e, m int) int {
   res := 1
   base := a % m
   for e > 0 {
       if e&1 == 1 {
           res = int(int64(res) * int64(base) % int64(m))
       }
       base = int(int64(base) * int64(base) % int64(m))
       e >>= 1
   }
   return res
}
