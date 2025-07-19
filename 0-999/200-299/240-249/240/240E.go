package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF = 0x7f7f7f7f

// Edge represents an edge in the graph or auxiliary structure
type Edge struct {
   u, v    int
   w       int
   sg, use int
}

var (
	n, m, Rt         int
	ans             int64
	pre, sg, mi, li, vis []int
	a, b           []Edge
)

// solve performs one iteration of directed MST contraction.
// returns 0 if no arborescence exists, 1 if finished, 2 if contracted and should repeat
func solve() int {
   // reset labels
   for i := 1; i <= n; i++ {
       pre[i], sg[i], li[i], vis[i] = 0, 0, 0, 0
       mi[i] = INF
   }
   // find minimum incoming edge for each node
   for i := 1; i <= m; i++ {
       if b[i].u != b[i].v && b[i].w < mi[b[i].v] {
           pre[b[i].v] = b[i].u
           mi[b[i].v] = b[i].w
           li[b[i].v] = b[i].sg
       }
   }
   // check reachability
   for i := 1; i <= n; i++ {
       if i != Rt && pre[i] == 0 {
           return 0
       }
   }
   cnt := 0
   mi[Rt] = 0
   // mark used edges and detect cycles
   for i := 1; i <= n; i++ {
       if i != Rt {
           a[li[i]].use++
       }
       ans += int64(mi[i])
       now := i
       for vis[now] == 0 && now != Rt {
           vis[now] = i
           now = pre[now]
       }
       if now != Rt && vis[now] == i {
           cnt++
           k := now
           for {
               sg[now] = cnt
               now = pre[now]
               if sg[now] == cnt {
                   break
               }
           }
       }
   }
   if cnt == 0 {
       return 1
   }
   // assign component labels to acyclic nodes
   for i := 1; i <= n; i++ {
       if sg[i] == 0 {
           cnt++
           sg[i] = cnt
       }
   }
   // contract edges, build new a and update b
   oldA := len(a)
   for i := 1; i <= m; i++ {
       k := mi[b[i].v]
       l := b[i].v
       u2 := sg[b[i].u]
       v2 := sg[b[i].v]
       b[i].u = u2
       b[i].v = v2
       if u2 != v2 {
           b[i].w -= k
           // append auxiliary edge
           a = append(a, Edge{u: b[i].sg, v: li[l]})
           b[i].sg = len(a) - 1
       }
   }
   n = cnt
   Rt = sg[Rt]
   return 2
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m)
   Rt = 1
   // initialize slices
   pre = make([]int, n+5)
   sg = make([]int, n+5)
   mi = make([]int, n+5)
   li = make([]int, n+5)
   vis = make([]int, n+5)
   a = make([]Edge, m+1)
   b = make([]Edge, m+1)
   // read edges
   for i := 1; i <= m; i++ {
       var u, v, w int
       fmt.Fscan(in, &u, &v, &w)
       b[i].u, b[i].v, b[i].w = u, v, w
       b[i].sg = i
       a[i].w = w
   }
   // iterate contractions
   res := solve()
   for res == 2 {
       res = solve()
   }
   if res == 0 {
       fmt.Fprintln(out, -1)
       return
   }
   // output result
   fmt.Fprintln(out, ans)
   // propagate uses from auxiliary edges back to original
   for i := len(a) - 1; i > m; i-- {
       a[a[i].u].use += a[i].use
       a[a[i].v].use -= a[i].use
   }
   // print original edges used
   first := true
   for i := 1; i <= m; i++ {
       if a[i].w != 0 && a[i].use != 0 {
           if !first {
               fmt.Fprint(out, " ")
           }
           first = false
           fmt.Fprint(out, i)
       }
   }
   fmt.Fprintln(out)
}
