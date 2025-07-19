package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxN = 1000005

type edge struct{ to, id int }

var (
   g     [maxN][]edge
   d     [maxN][]int
   used  [maxN]bool
   color [maxN]bool
   same  [maxN]int
   diff  [maxN]int
   h     [maxN]int
   cc    int
   kk    bool
   ans   bool
   vv    int
)

func cfs(v, p int) {
   used[v] = true
   for _, e := range g[v] {
       to, id := e.to, e.id
       if id == p {
           continue
       }
       if !used[to] {
           color[to] = !color[v]
           cfs(to, id)
           h[to] = h[v] + 1
           d[v] = append(d[v], to)
       } else {
           if h[to] > h[v] {
               if color[v] == color[to] {
                   same[to]++; same[v]--
                   cc++; kk = color[v]
               } else {
                   diff[to]++; diff[v]--
               }
           }
       }
   }
}

func dfs(v int) (int, int) {
   f0, f1 := same[v], diff[v]
   for _, to := range d[v] {
       a, b := dfs(to)
       f0 += a; f1 += b
   }
   if !ans && f0 == cc && f1 == 0 {
       ans = true
       vv = v
       kk = !color[v]
   }
   return f0, f1
}

func hu(v int) {
   color[v] = !color[v]
   for _, to := range d[v] {
       hu(to)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t, n, m int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       fmt.Fscan(in, &n, &m)
       // reset
       cc, kk, ans = 0, false, false
       for i := 1; i <= n; i++ {
           g[i] = g[i][:0]
           d[i] = d[i][:0]
           used[i] = false
           color[i] = false
           same[i] = 0
           diff[i] = 0
           h[i] = 0
       }
       for i := 0; i < m; i++ {
           var u, v0 int
           fmt.Fscan(in, &u, &v0)
           g[u] = append(g[u], edge{v0, i})
           g[v0] = append(g[v0], edge{u, i})
       }
       cfs(1, -1)
       // clear used for dfs
       for i := 1; i <= n; i++ {
           used[i] = false
       }
       if cc <= 1 {
           fmt.Fprintln(out, "YES")
           for i := 1; i <= n; i++ {
               bit := color[i] != (!kk)
               if bit {
                   out.WriteByte('1')
               } else {
                   out.WriteByte('0')
               }
           }
           out.WriteByte('\n')
           continue
       }
       dfs(1)
       if ans {
           hu(vv)
           fmt.Fprintln(out, "YES")
           for i := 1; i <= n; i++ {
               bit := color[i] != (!kk)
               if bit {
                   out.WriteByte('1')
               } else {
                   out.WriteByte('0')
               }
           }
           out.WriteByte('\n')
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
