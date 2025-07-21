package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k, l int
   if _, err := fmt.Fscan(in, &n, &k, &l); err != nil {
       return
   }
   b := make([]bool, n+2)
   x := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &x[i])
       if x[i] >= 1 && x[i] <= n {
           b[x[i]] = true
       }
   }
   a := make([]int, l)
   for i := 0; i < l; i++ {
       fmt.Fscan(in, &a[i])
   }
   // compute transition nodes D where b[i] != b[i-1]
   D := make([]int, 0, 20)
   prev := false
   for i := 1; i <= n; i++ {
       if b[i] != prev {
           D = append(D, i)
       }
       prev = b[i]
   }
   // check end at n+1
   if prev {
       D = append(D, n+1)
   }
   m := len(D)
   if m == 0 {
       fmt.Println(0)
       return
   }
   if m%2 == 1 {
       fmt.Println(-1)
       return
   }
   // build graph nodes 1..n+1
   N := n + 2
   adj := make([][]int, N)
   for s := 1; s <= n+1; s++ {
       for _, ai := range a {
           t := s + ai
           if t <= n+1 {
               adj[s] = append(adj[s], t)
               adj[t] = append(adj[t], s)
           }
       }
   }
   // compute pairwise distances between D nodes via BFS
   const INF = 1e9
   dist := make([][]int, m)
   for i := 0; i < m; i++ {
       // BFS from D[i]
       d := make([]int, N)
       for j := range d {
           d[j] = -1
       }
       qi := 0
       qj := 0
       q := make([]int, N)
       start := D[i]
       d[start] = 0
       q[qj] = start
       qj++
       for qi < qj {
           u := q[qi]
           qi++
           for _, v := range adj[u] {
               if d[v] == -1 {
                   d[v] = d[u] + 1
                   q[qj] = v
                   qj++
               }
           }
       }
       dist[i] = make([]int, m)
       for j := 0; j < m; j++ {
           dval := d[D[j]]
           if dval >= 0 {
               dist[i][j] = dval
           } else {
               dist[i][j] = INF
           }
       }
   }
   // DP on matching
   maxMask := 1 << m
   dp := make([]int, maxMask)
   for i := 1; i < maxMask; i++ {
       dp[i] = INF
   }
   for mask := 1; mask < maxMask; mask++ {
       // only even bits
       if bitsCount(mask)%2 == 1 {
           continue
       }
       // find first bit
       var i0 int
       for bit := 0; bit < m; bit++ {
           if (mask>>bit)&1 == 1 {
               i0 = bit
               break
           }
       }
       // pair i0 with any j>i0
       for j := i0 + 1; j < m; j++ {
           if (mask>>j)&1 == 1 {
               m2 := mask ^ (1 << i0) ^ (1 << j)
               cost := dp[m2] + dist[i0][j]
               if cost < dp[mask] {
                   dp[mask] = cost
               }
           }
       }
   }
   ans := dp[maxMask-1]
   if ans >= INF/2 {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}

func bitsCount(x int) int {
   // builtin popcount
   cnt := 0
   for x > 0 {
       x &= x - 1
       cnt++
   }
   return cnt
}
