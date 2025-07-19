package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   if a > b {
       return gcd(b, a)
   }
   if a == 0 {
       return b
   }
   return gcd(b%a, a)
}

func lcm(a, b int64) int64 {
   return a / gcd(a, b) * b
}

type Edge struct {
   to    int
   g, l  int64
}

var (
   n, m       int
   adj        [][]Edge
   A          []int64
   visited    []bool
   order      []int
   okFlag     bool
)

func dfs(u int) {
   visited[u] = true
   order = append(order, u)
   for _, e := range adj[u] {
       v := e.to
       if visited[v] {
           if lcm(A[u], A[v]) != e.l || gcd(A[u], A[v]) != e.g {
               okFlag = false
           }
           continue
       }
       if e.l % A[u] != 0 {
           okFlag = false
           continue
       }
       A[v] = e.l / A[u] * e.g
       dfs(v)
       if !okFlag {
           return
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m)
   adj = make([][]Edge, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       var g, l int64
       fmt.Fscan(reader, &x, &y, &g, &l)
       adj[x] = append(adj[x], Edge{y, g, l})
       adj[y] = append(adj[y], Edge{x, g, l})
   }
   A = make([]int64, n+1)
   visited = make([]bool, n+1)
   // process each component
   for i := 1; i <= n; i++ {
       if visited[i] || len(adj[i]) == 0 {
           continue
       }
       // compute tmp = gcd of all l's at i, tmp2 = lcm of all g's
       var tmp, tmp2 int64
       for j, e := range adj[i] {
           if j == 0 {
               tmp = e.l
               tmp2 = e.g
           } else {
               tmp = gcd(tmp, e.l)
               tmp2 = lcm(tmp2, e.g)
           }
       }
       if tmp2 > tmp {
           fmt.Println("NO")
           return
       }
       // find candidates D dividing tmp and D % tmp2 == 0
       var cands []int64
       for d := int64(1); d*d <= tmp; d++ {
           if tmp % d == 0 {
               if d % tmp2 == 0 {
                   cands = append(cands, d)
               }
               oth := tmp / d
               if oth != d && oth % tmp2 == 0 {
                   cands = append(cands, oth)
               }
           }
       }
       ok := false
       for _, D := range cands {
           // try candidate
           A[i] = D
           order = order[:0]
           okFlag = true
           dfs(i)
           if okFlag {
               ok = true
               break
           }
           // reset visited for this attempt
           for _, u := range order {
               visited[u] = false
           }
       }
       if !ok {
           fmt.Println("NO")
           return
       }
   }
   // assign 1 to isolated or untouched
   for i := 1; i <= n; i++ {
       if A[i] == 0 {
           A[i] = 1
       }
   }
   fmt.Println("YES")
   for i := 1; i <= n; i++ {
       fmt.Printf("%d ", A[i])
   }
   fmt.Println()
}
