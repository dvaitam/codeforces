package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var p, n, M, m int
   if _, err := fmt.Fscan(in, &p, &n, &M, &m); err != nil {
       return
   }
   // intervals
   l := make([]int, n+1)
   r := make([]int, n+1)
   // constraints pairs
   type pair struct{ x, y int }
   pairs1 := make([]pair, p)
   for i := 0; i < p; i++ {
       fmt.Fscan(in, &pairs1[i].x, &pairs1[i].y)
   }
   // read intervals
   aL := make([]int, n)
   aR := make([]int, n)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &l[i], &r[i])
       aL[i-1] = l[i]
       aR[i-1] = r[i]
   }
   // later m pairs
   pairs2 := make([]pair, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &pairs2[i].x, &pairs2[i].y)
   }
   // total nodes: 2*n for boolean, plus 2*M for interval graph
   N := 2*n + 2*M
   adj := make([][]int, N+1)
   // helper to add edge
   add := func(u, v int) {
       if u >= 1 && u <= N {
           adj[u] = append(adj[u], v)
       }
   }
   // initial 2-SAT edges from pairs1: x false => y true, y false => x true
   for _, pr := range pairs1 {
       x, y := pr.x, pr.y
       add(x+n, y)
       add(y+n, x)
   }
   // interval constraints: prepare order by l and r with ids
   ids := make([]int, n)
   for i := 1; i <= n; i++ {
       ids[i-1] = i
   }
   // sort by l ascending
   // sort by l
   {
       sl := ids
       // copy
       tmp := make([]int, n)
       copy(tmp, sl)
       // sort tmp
       sort.Slice(tmp, func(i, j int) bool { return l[tmp[i]] < l[tmp[j]] })
       sl = tmp
       // chain suffix edges 1->2,2->3,...
       for i := 1; i < M; i++ {
           add(i+2*n, i+2*n+1)
       }
       for _, ii := range sl {
           add(l[ii]+2*n, ii+n)
       }
   }
   // sort by r
   {
       sl := ids
       tmp := make([]int, n)
       copy(tmp, sl)
       sort.Slice(tmp, func(i, j int) bool { return r[tmp[i]] < r[tmp[j]] })
       sl = tmp
       for i := 2; i <= M; i++ {
           add(i+2*n+M, i+2*n+M-1)
       }
       for _, ii := range sl {
           add(r[ii]+2*n+M, ii+n)
       }
   }
   // link variable to interval endpoints
   for i := 1; i <= n; i++ {
       if r[i] != M {
           add(i, r[i]+2*n+1)
       }
       if l[i] != 1 {
           add(i, l[i]-1+2*n+M)
       }
   }
   // pairs2 constraints
   for _, pr := range pairs2 {
       x, y := pr.x, pr.y
       add(x, y+n)
       add(y, x+n)
   }
   // Kosaraju
   vis := make([]bool, N+1)
   order := make([]int, 0, N)
   var stack []struct{ u, i int }
   // first pass
   for u := 1; u <= N; u++ {
       if !vis[u] {
           stack = append(stack, struct{ u, i int }{u, 0})
           vis[u] = true
           for len(stack) > 0 {
               top := &stack[len(stack)-1]
               if top.i < len(adj[top.u]) {
                   v := adj[top.u][top.i]
                   top.i++
                   if !vis[v] {
                       vis[v] = true
                       stack = append(stack, struct{ u, i int }{v, 0})
                   }
               } else {
                   order = append(order, top.u)
                   stack = stack[:len(stack)-1]
               }
           }
       }
   }
   // build reverse graph
   radj := make([][]int, N+1)
   for u, nbrs := range adj {
       for _, v := range nbrs {
           radj[v] = append(radj[v], u)
       }
   }
   comp := make([]int, N+1)
   cid := 0
   vis2 := make([]bool, N+1)
   // second pass
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       if !vis2[u] {
           cid++
           stack = append(stack, struct{ u, i int }{u, 0})
           vis2[u] = true
           comp[u] = cid
           for len(stack) > 0 {
               top := &stack[len(stack)-1]
               if top.i < len(radj[top.u]) {
                   v := radj[top.u][top.i]
                   top.i++
                   if !vis2[v] {
                       vis2[v] = true
                       comp[v] = cid
                       stack = append(stack, struct{ u, i int }{v, 0})
                   }
               } else {
                   stack = stack[:len(stack)-1]
               }
           }
       }
   }
   // check and collect
   res := make([]int, 0, n)
   K := 0
   L := 1
   R := M
   for i := 1; i <= n; i++ {
       if comp[i] == comp[i+n] {
           fmt.Fprintln(out, -1)
           return
       }
       if comp[i] > comp[i+n] {
           K++
           if l[i] > L {
               L = l[i]
           }
           if r[i] < R {
               R = r[i]
           }
           res = append(res, i)
       }
   }
   // output
   fmt.Fprintf(out, "%d %d\n", K, L)
   for i, v := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprintf(out, "%d", v)
   }
   out.WriteByte('\n')
}
