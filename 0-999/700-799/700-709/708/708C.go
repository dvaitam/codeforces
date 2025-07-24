package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   parent := make([]int, n+1)
   order := make([]int, 0, n)
   // build parent and order (preorder)
   stk := make([]int, 0, n)
   stk = append(stk, 1)
   parent[1] = 0
   for len(stk) > 0 {
       v := stk[len(stk)-1]
       stk = stk[:len(stk)-1]
       order = append(order, v)
       for _, u := range adj[v] {
           if u == parent[v] {
               continue
           }
           parent[u] = v
           stk = append(stk, u)
       }
   }
   sz := make([]int, n+1)
   // child max1, max2
   childMax1 := make([]int, n+1)
   childMax2 := make([]int, n+1)
   childMaxNode := make([]int, n+1)
   // postorder to compute sz and child max
   for i := n - 1; i >= 0; i-- {
       v := order[i]
       sz[v] = 1
       m1, m2, mn := 0, 0, 0
       for _, u := range adj[v] {
           if u == parent[v] {
               continue
           }
           sz[v] += sz[u]
           // track top two child subtree sizes
           if sz[u] > m1 {
               m2 = m1
               m1 = sz[u]
               mn = u
           } else if sz[u] > m2 {
               m2 = sz[u]
           }
       }
       childMax1[v] = m1
       childMax2[v] = m2
       childMaxNode[v] = mn
   }
   // upSz and comp max for each
   upSz := make([]int, n+1)
   compMax1 := make([]int, n+1)
   compMax2 := make([]int, n+1)
   compMaxNode := make([]int, n+1)
   half := n / 2
   // compute upSz and comp max
   for _, v := range order {
       // upSz
       if v == 1 {
           upSz[v] = 0
       } else {
           upSz[v] = n - sz[v]
       }
       // combine childMax and upSz
       // start with child values
       m1, m2, mn := childMax1[v], childMax2[v], childMaxNode[v]
       // consider upSz
       if upSz[v] > m1 {
           m2 = m1
           m1 = upSz[v]
           mn = parent[v]
       } else if upSz[v] > m2 {
           m2 = upSz[v]
       }
       compMax1[v] = m1
       compMax2[v] = m2
       compMaxNode[v] = mn
   }
   // answer
   out := make([]byte, n)
   for v := 1; v <= n; v++ {
       m := compMax1[v]
       if m <= half {
           out[v-1] = '1'
           continue
       }
       w := compMaxNode[v]
       // local max at w excluding v
       var local int
       if compMaxNode[w] != v {
           local = compMax1[w]
       } else {
           local = compMax2[w]
       }
       if 2*local >= m {
           out[v-1] = '1'
       } else {
           out[v-1] = '0'
       }
   }
   // print
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < n; i++ {
       writer.WriteByte(out[i])
       if i+1 < n {
           writer.WriteByte(' ')
       }
   }
}
