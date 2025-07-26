package main

import (
   "bufio"
   "os"
)

const LOG = 18

var (
   adj   [][]int
   up    [LOG][]int
   depth []int
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   readInt := func() int {
       b, _ := reader.ReadByte()
       for (b < '0' || b > '9') && b != '-' {
           b, _ = reader.ReadByte()
       }
       neg := false
       if b == '-' {
           neg = true
           b, _ = reader.ReadByte()
       }
       x := 0
       for b >= '0' && b <= '9' {
           x = x*10 + int(b-'0')
           b, _ = reader.ReadByte()
       }
       if neg {
           return -x
       }
       return x
   }

   n := readInt()
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       u := readInt()
       v := readInt()
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   depth = make([]int, n+1)
   for i := 0; i < LOG; i++ {
       up[i] = make([]int, n+1)
   }
   // root at 1
   depth[1] = 0
   up[0][1] = 1
   // BFS to set depth and parent
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range adj[u] {
           if v == up[0][u] {
               continue
           }
           depth[v] = depth[u] + 1
           up[0][v] = u
           queue = append(queue, v)
       }
   }
   // binary lifting
   for i := 1; i < LOG; i++ {
       for v := 1; v <= n; v++ {
           up[i][v] = up[i-1][ up[i-1][v] ]
       }
   }
   // process queries
   q := readInt()
   for i := 0; i < q; i++ {
       x := readInt()
       y := readInt()
       a := readInt()
       b := readInt()
       k := readInt()
       d1 := dist(a, b)
       d2 := dist(a, x) + 1 + dist(y, b)
       d3 := dist(a, y) + 1 + dist(x, b)
       ok := false
       if d1 <= k && (k-d1)%2 == 0 {
           ok = true
       }
       if d2 <= k && (k-d2)%2 == 0 {
           ok = true
       }
       if d3 <= k && (k-d3)%2 == 0 {
           ok = true
       }
       if ok {
           writer.WriteString("YES\n")
       } else {
           writer.WriteString("NO\n")
       }
   }
}

func lca(u, v int) int {
   if depth[u] < depth[v] {
       u, v = v, u
   }
   // lift u up to depth v
   diff := depth[u] - depth[v]
   for i := 0; i < LOG; i++ {
       if (diff>>i)&1 == 1 {
           u = up[i][u]
       }
   }
   if u == v {
       return u
   }
   for i := LOG - 1; i >= 0; i-- {
       if up[i][u] != up[i][v] {
           u = up[i][u]
           v = up[i][v]
       }
   }
   return up[0][u]
}

func dist(u, v int) int {
   w := lca(u, v)
   return depth[u] + depth[v] - 2*depth[w]
}
