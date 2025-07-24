package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   p := make([]int, n)
   q := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &q[i])
   }
   mod := 998244353
   // Build partial graph with known edges
   out := make([]int, n+1)
   in := make([]int, n+1)
   for i := 0; i < n; i++ {
       if p[i] != 0 && q[i] != 0 {
           out[p[i]] = q[i]
           in[q[i]] = p[i]
       }
   }
   // Count initial cycles
   vis := make([]bool, n+1)
   cycles := 0
   for v := 1; v <= n; v++ {
       if !vis[v] {
           u := v
           seen := make(map[int]bool)
           for u != 0 && !vis[u] {
               seen[u] = true
               vis[u] = true
               u = out[u]
               if u != 0 && seen[u] {
                   cycles++
                   break
               }
           }
       }
   }
   // Count zero slots
   zeroZero := 0
   mz1, mz2 := 0, 0
   for i := 0; i < n; i++ {
       if p[i] == 0 && q[i] == 0 {
           zeroZero++
       } else if p[i] == 0 {
           mz1++
       } else if q[i] == 0 {
           mz2++
       }
   }
   // Number of zero slots of type p==0 should equal type q==0
   // Precompute factorials up to mz1+zeroZero
   size := mz1 + zeroZero
   fact := make([]int, size+1)
   fact[0] = 1
   for i := 1; i <= size; i++ {
       fact[i] = int(int64(fact[i-1]) * int64(i) % int64(mod))
   }
   fillWays := int(int64(fact[size]) * int64(fact[size]) % int64(mod))
   // Derangements
   D := make([]int, zeroZero+1)
   D[0] = 1
   if zeroZero >= 1 {
       D[1] = 0
   }
   for i := 2; i <= zeroZero; i++ {
       D[i] = int((int64(i-1) * (int64(D[i-1]) + int64(D[i-2]))) % int64(mod))
   }
   // Binomial coefficients for zeroZero
   C := make([][]int, zeroZero+1)
   for i := 0; i <= zeroZero; i++ {
       C[i] = make([]int, i+1)
       C[i][0] = 1
       for j := 1; j <= i; j++ {
           if j == i {
               C[i][j] = 1
           } else {
               C[i][j] = (C[i-1][j-1] + C[i-1][j]) % mod
           }
       }
   }
   ans := make([]int, n)
   for k := 0; k < n; k++ {
       u := n - cycles - k
       if u < 0 || u > zeroZero {
           ans[k] = 0
       } else {
           ways := int64(fillWays)
           ways = ways * int64(C[zeroZero][u]) % int64(mod)
           ways = ways * int64(D[zeroZero-u]) % int64(mod)
           ans[k] = int(ways)
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 0; i < n; i++ {
       if i > 0 {
           w.WriteByte(' ')
       }
       w.WriteString(strconv.Itoa(ans[i]))
   }
   w.WriteByte('\n')
}
