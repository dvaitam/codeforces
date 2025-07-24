package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
   "sort"
)

func main() {
   var Ts string
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &Ts)
   T := new(big.Int)
   T.SetString(Ts, 10)
   // grid size
   n, m := 50, 50
   // compute dp1
   dp1 := make([][]*big.Int, n+2)
   dp2 := make([][]*big.Int, n+2)
   for i := range dp1 {
       dp1[i] = make([]*big.Int, m+2)
       dp2[i] = make([]*big.Int, m+2)
       for j := range dp1[i] {
           dp1[i][j] = new(big.Int)
           dp2[i][j] = new(big.Int)
       }
   }
   dp1[1][1].SetInt64(1)
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if i == 1 && j == 1 {
               continue
           }
           dp1[i][j].Add(dp1[i-1][j], dp1[i][j-1])
       }
   }
   dp2[n][m].SetInt64(1)
   for i := n; i >= 1; i-- {
       for j := m; j >= 1; j-- {
           if i == n && j == m {
               continue
           }
           dp2[i][j].Add(dp2[i+1][j], dp2[i][j+1])
       }
   }
   total := new(big.Int).Set(dp1[n][m])
   // carry = total - T
   carry := new(big.Int).Sub(total, T)
   type Edge struct{ i, j, ni, nj int; contrib *big.Int }
   var edges []Edge
   // collect edges
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           if i < n {
               contrib := new(big.Int).Mul(dp1[i][j], dp2[i+1][j])
               edges = append(edges, Edge{i, j, i + 1, j, contrib})
           }
           if j < m {
               contrib := new(big.Int).Mul(dp1[i][j], dp2[i][j+1])
               edges = append(edges, Edge{i, j, i, j + 1, contrib})
           }
       }
   }
   // sort edges by contrib descending
   sort.Slice(edges, func(a, b int) bool {
       return edges[a].contrib.Cmp(edges[b].contrib) > 0
   })
   // select edges to block
   var blocked []Edge
   zero := big.NewInt(0)
   for _, e := range edges {
       if carry.Cmp(zero) <= 0 {
           break
       }
       if e.contrib.Cmp(carry) <= 0 {
           // block this edge
           blocked = append(blocked, e)
           carry.Sub(carry, e.contrib)
       }
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintf(w, "%d %d\n", n, m)
   fmt.Fprintf(w, "%d\n", len(blocked))
   for _, e := range blocked {
       fmt.Fprintf(w, "%d %d %d %d\n", e.i, e.j, e.ni, e.nj)
   }
}
