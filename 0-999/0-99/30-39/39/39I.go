package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
   }
   // BFS from 1 to get distances
   dist := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = -1
   }
   q := make([]int, 0, n)
   dist[1] = 0
   q = append(q, 1)
   for qi := 0; qi < len(q); qi++ {
       u := q[qi]
       for _, v := range adj[u] {
           if dist[v] == -1 {
               dist[v] = dist[u] + 1
               q = append(q, v)
           }
       }
   }
   // compute gcd of cycle lengths
   g := 0
   for u := 1; u <= n; u++ {
       if dist[u] < 0 {
           continue
       }
       for _, v := range adj[u] {
           if dist[v] >= 0 {
               // delta = dist[u] + 1 - dist[v]
               d := dist[u] + 1 - dist[v]
               if d < 0 {
                   d = -d
               }
               g = gcd(g, d)
           }
       }
   }
   // t is g
   t := g
   if t <= 0 {
       t = 1
   }
   // choose nodes with dist % t == 0
   cams := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if dist[i] >= 0 && dist[i]%t == 0 {
           cams = append(cams, i)
       }
   }
   // output
   fmt.Fprintln(out, t)
   fmt.Fprintln(out, len(cams))
   for i, v := range cams {
       if i > 0 {
           out.WriteString(" ")
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
