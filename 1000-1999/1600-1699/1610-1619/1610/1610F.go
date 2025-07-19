package main

import (
   "bufio"
   "fmt"
   "os"
)

type edgeRef struct{ to, id int }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   fmt.Fscan(in, &n, &m)
   // initial capacity for merged edges
   maxE := m*3 + 5
   u := make([]int, maxE)
   v := make([]int, maxE)
   w := make([]int, maxE)
   l := make([]int, maxE)
   r := make([]int, maxE)
   vs := make([]bool, maxE)
   ed := make([]int, maxE)
   // f[x][label] -> edge id
   f := make([][3]int, n+1)
   // read edges
   q := make([]int, maxE)
   ql, qr := 0, 0
   for i := 1; i <= m; i++ {
       fmt.Fscan(in, &u[i], &v[i], &w[i])
       q[qr] = i
       qr++
   }
   idx := m
   // merge matching labels
   for ql < qr {
       x := q[ql]; ql++
       var y int
       if f[u[x]][w[x]] != 0 {
           y = f[u[x]][w[x]]
       } else if f[v[x]][w[x]] != 0 {
           y = f[v[x]][w[x]]
       }
       if y != 0 {
           // remove old
           f[u[y]][w[y]] = 0
           f[v[y]][w[y]] = 0
           idx++
           z := idx
           l[z], r[z] = x, y
           // new endpoints
           if u[x] == u[y] {
               u[z] = v[x]
               v[z] = v[y]
           } else if u[x] == v[y] {
               u[z] = v[x]
               v[z] = u[y]
           } else if v[x] == u[y] {
               u[z] = u[x]
               v[z] = v[y]
           } else {
               u[z] = u[x]
               v[z] = u[y]
           }
           w[z] = w[x]
           vs[x], vs[y] = true, true
           q[qr] = z; qr++
       } else {
           f[u[x]][w[x]] = x
           f[v[x]][w[x]] = x
       }
   }
   // build graph of unmerged
   deg := make([]int, n+1)
   g := make([][]edgeRef, n+1)
   for i := 1; i <= idx; i++ {
       if !vs[i] {
           g[u[i]] = append(g[u[i]], edgeRef{v[i], i})
           g[v[i]] = append(g[v[i]], edgeRef{u[i], i})
           deg[u[i]]++
           deg[v[i]]++
       }
   }
   // dfs to set ed on tree-like parts
   seen := make([]bool, n+1)
   var dfs1 func(x int)
   dfs1 = func(x int) {
       seen[x] = true
       for _, eRef := range g[x] {
           y, id := eRef.to, eRef.id
           if seen[y] || vs[id] {
               continue
           }
           if u[id] == x {
               ed[id] = 1
           } else {
               ed[id] = 2
           }
           dfs1(y)
       }
   }
   // dfs2 for cycles
   var dfs2 func(x, pid int)
   dfs2 = func(x, pid int) {
       if seen[x] {
           return
       }
       seen[x] = true
       for _, eRef := range g[x] {
           y, id := eRef.to, eRef.id
           if id == pid || vs[id] {
               continue
           }
           if u[id] == x {
               ed[id] = 1
           } else {
               ed[id] = 2
           }
           dfs2(y, id)
           break
       }
   }
   for i := 1; i <= n; i++ {
       if !seen[i] && deg[i] <= 1 {
           dfs1(i)
       }
   }
   for i := 1; i <= n; i++ {
       if !seen[i] {
           dfs2(i, 0)
       }
   }
   // propagate directions
   var calc func(z int)
   calc = func(z int) {
       if l[z] == 0 {
           return
       }
       if ed[z] == 2 {
           // swap endpoints and children
           u[z], v[z] = v[z], u[z]
           l[z], r[z] = r[z], l[z]
       }
       // left child
       lx := l[z]
       if u[lx] == u[z] {
           ed[lx] = 1
       } else {
           ed[lx] = 2
       }
       // right child
       rx := r[z]
       if v[rx] == v[z] {
           ed[rx] = 1
       } else {
           ed[rx] = 2
       }
       calc(lx)
       calc(rx)
   }
   for i := 1; i <= idx; i++ {
       if !vs[i] {
           calc(i)
       }
   }
   // compute balances
   bal := make([]int, n+1)
   for i := 1; i <= m; i++ {
       if ed[i] == 1 {
           bal[u[i]] += w[i]
           bal[v[i]] -= w[i]
       } else {
           bal[v[i]] += w[i]
           bal[u[i]] -= w[i]
       }
   }
   // count +/-1
   cnt := 0
   for i := 1; i <= n; i++ {
       if bal[i] == 1 || bal[i] == -1 {
           cnt++
       }
   }
   fmt.Fprintln(out, cnt)
   // print directions for original edges
   for i := 1; i <= m; i++ {
       out.WriteByte(byte('0' + ed[i]))
   }
}
