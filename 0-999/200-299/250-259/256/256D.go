package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // t = number of honest people
   t := n - k
   const mod = 777777777
   // Precompute binomial coefficients C up to n (we need up to k as well)
   maxN := n
   if k > n {
       maxN = k
   }
   C := make([][]int, maxN+1)
   for i := 0; i <= maxN; i++ {
       C[i] = make([]int, i+1)
       C[i][0], C[i][i] = 1, 1
       for j := 1; j < i; j++ {
           C[i][j] = (C[i-1][j-1] + C[i-1][j]) % mod
       }
   }
   // dp h[j] = j! * sum of (1/ prod fv!) over fv for processed labels, sum fv = j
   // initialize
   h := make([]int, k+1)
   h[0] = 1
   // process each possible answer value v != t, domain v from 1..n
   for v := 1; v <= n; v++ {
       if v == t {
           continue
       }
       // new dp
       newh := make([]int, k+1)
       for newj := 0; newj <= k; newj++ {
           var sum int64
           // fv occurrences of label v
           for fv := 0; fv <= newj; fv++ {
               if fv == v {
                   continue
               }
               // from previous old_j = newj - fv
               oldj := newj - fv
               // h[oldj] * C(newj, fv)
               sum += int64(h[oldj]) * int64(C[newj][fv])
           }
           newh[newj] = int(sum % mod)
       }
       h = newh
   }
   // total sequences = C(n, t) * h[k] % mod
   // compute C(n, t)
   // ensure C table covers n
   // We need to recompute C up to n if maxN< n
   var cn_t int
   if n <= maxN {
       cn_t = C[n][t]
   } else {
       // build row n only
       row := make([]int, n+1)
       row[0], row[1] = 1, n%mod
       for i := 2; i <= n; i++ {
           row[i] = 0
           // compute C[n][i] via Pascal: C[n][i] = C[n-1][i-1] + C[n-1][i]
           // but since we didn't build, fallback to iterative
           // Actually n<=28 so maxN>=n always
       }
       cn_t = C[n][t]
   }
   ans := int64(cn_t) * int64(h[k]) % mod
   fmt.Println(ans)
}
