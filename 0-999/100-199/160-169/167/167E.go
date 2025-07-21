package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var p int64
   if _, err := fmt.Fscan(in, &n, &m, &p); err != nil {
       return
   }
   adj := make([][]int, n)
   indeg := make([]int, n)
   outdeg := make([]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       indeg[v]++
       outdeg[u]++
   }
   // identify sources and sinks
   var src, sink []int
   for i := 0; i < n; i++ {
       if indeg[i] == 0 {
           src = append(src, i)
       }
       if outdeg[i] == 0 {
           sink = append(sink, i)
       }
   }
   k := len(src)
   if k == 0 {
       // empty matching: one empty set, even parity
       fmt.Println(1 % p)
       return
   }
   // topological order
   dq := make([]int, 0, n)
   indeg0 := make([]int, n)
   copy(indeg0, indeg)
   for i := 0; i < n; i++ {
       if indeg0[i] == 0 {
           dq = append(dq, i)
       }
   }
   topo := make([]int, 0, n)
   for idx := 0; idx < len(dq); idx++ {
       u := dq[idx]
       topo = append(topo, u)
       for _, v := range adj[u] {
           indeg0[v]--
           if indeg0[v] == 0 {
               dq = append(dq, v)
           }
       }
   }
   // build matrix
   M := make([][]int64, k)
   for i := range M {
       M[i] = make([]int64, k)
   }
   dp := make([]int64, n)
   for i := 0; i < k; i++ {
       // reset dp
       for j := 0; j < n; j++ {
           dp[j] = 0
       }
       dp[src[i]] = 1
       for _, u := range topo {
           if dp[u] != 0 {
               x := dp[u]
               for _, v := range adj[u] {
                   dp[v] += x
                   if dp[v] >= p {
                       dp[v] -= p
                   }
               }
           }
       }
       for j := 0; j < k; j++ {
           M[i][j] = dp[sink[j]]
       }
   }
   // determinant mod p
   det := int64(1)
   sign := int64(1)
   for i := 0; i < k; i++ {
       // find pivot
       piv := -1
       for r := i; r < k; r++ {
           if M[r][i] != 0 {
               piv = r
               break
           }
       }
       if piv == -1 {
           det = 0
           break
       }
       if piv != i {
           M[i], M[piv] = M[piv], M[i]
           sign = p - sign
       }
       inv := modInv(M[i][i], p)
       for r := i + 1; r < k; r++ {
           if M[r][i] != 0 {
               factor := M[r][i] * inv % p
               // subtract factor * row i
               for c := i; c < k; c++ {
                   M[r][c] = (M[r][c] - factor*M[i][c]) % p
                   if M[r][c] < 0 {
                       M[r][c] += p
                   }
               }
           }
       }
       det = det * M[i][i] % p
   }
   det = det * sign % p
   if det < 0 {
       det += p
   }
   fmt.Println(det)
}

// modInv returns modular inverse of a mod p, p prime
func modInv(a, p int64) int64 {
   return modPow(a, p-2, p)
}

func modPow(a, e, p int64) int64 {
   res := int64(1)
   a %= p
   for e > 0 {
       if e&1 != 0 {
           res = res * a % p
       }
       a = a * a % p
       e >>= 1
   }
   return res
}
