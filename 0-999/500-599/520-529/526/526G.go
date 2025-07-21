package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct{ to, w int }

var (
   n, q int
   adj [][]edge
   down1, down2, child1 []int
   up []int
   parent []int
   D1, D2 []int
   Bs [][]int64
   PBs [][]int64
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func dfs1(u, p int) {
   parent[u] = p
   down1[u], down2[u] = 0, 0
   for _, e := range adj[u] {
       v, w := e.to, e.w
       if v == p {
           continue
       }
       dfs1(v, u)
       d := down1[v] + w
       if d > down1[u] {
           down2[u], down1[u] = down1[u], d
           child1[u] = v
       } else if d > down2[u] {
           down2[u] = d
       }
   }
}

func dfs2(u, p int) {
   for _, e := range adj[u] {
       v, w := e.to, e.w
       if v == p {
           continue
       }
       // up for v: via u
       viaUp := up[u] + w
       viaSib := 0
       if child1[u] == v {
           viaSib = down2[u] + w
       } else {
           viaSib = down1[u] + w
       }
       up[v] = max(viaUp, viaSib)
       dfs2(v, u)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &q)
   adj = make([][]edge, n)
   for i := 0; i < n-1; i++ {
       var u, v, w int
       fmt.Fscan(in, &u, &v, &w)
       u--; v--
       adj[u] = append(adj[u], edge{v, w})
       adj[v] = append(adj[v], edge{u, w})
   }
   down1 = make([]int, n)
   down2 = make([]int, n)
   child1 = make([]int, n)
   up = make([]int, n)
   parent = make([]int, n)
   dfs1(0, -1)
   up[0] = 0
   dfs2(0, -1)
   D1 = make([]int, n)
   D2 = make([]int, n)
   Bs = make([][]int64, n)
   PBs = make([][]int64, n)
   // build D1,D2 and branch lists
   for u := 0; u < n; u++ {
       // collect depths from neighbors
       var ds []int
       ds = append(ds, up[u])
       for _, e := range adj[u] {
           v, w := e.to, e.w
           if v == parent[u] {
               continue
           }
           ds = append(ds, down1[v]+w)
       }
       // find top two
       a, b := 0, 0
       for _, d := range ds {
           if d > a {
               b = a; a = d
           } else if d > b {
               b = d
           }
       }
       D1[u], D2[u] = a, b
       // branches: all depths except a and b once
       // remove one a and one b
       remA, remB := 1, 1
       var bs []int64
       for _, d := range ds {
           if d == a && remA > 0 {
               remA--; continue
           }
           if d == b && remB > 0 {
               remB--; continue
           }
           bs = append(bs, int64(d))
       }
       // sort bs desc
       if len(bs) > 1 {
           // simple sort
           for i := 0; i < len(bs); i++ {
               for j := i+1; j < len(bs); j++ {
                   if bs[j] > bs[i] {
                       bs[i], bs[j] = bs[j], bs[i]
                   }
               }
           }
       }
       Bs[u] = bs
       ps := make([]int64, len(bs)+1)
       for i := 0; i < len(bs); i++ {
           ps[i+1] = ps[i] + bs[i]
       }
       PBs[u] = ps
   }
   var x, y, ansPrev int64
   for i := 0; i < q; i++ {
       fmt.Fscan(in, &x, &y)
       if i > 0 {
           x = ((x + ansPrev - 1) % int64(n)) + 1
           y = ((y + ansPrev - 1) % int64(n)) + 1
       }
       u := int(x - 1)
       yi := int(y)
       // use D1[u],D2[u]
       base := int64(D1[u] + D2[u])
       // extra branches y-1
       k := yi - 1
       if k < 0 {
           k = 0
       }
       if k > len(PBs[u])-1 {
           k = len(PBs[u]) - 1
       }
       ans := base
       if k > 0 {
           ans += PBs[u][k]
       }
       fmt.Fprintln(out, ans)
       ansPrev = ans
   }
}
