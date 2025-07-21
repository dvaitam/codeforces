package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // original graph
   g := make([][]int, n)
   grev := make([][]int, n)
   edges := make([][2]int, m)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--;
       b--;
       edges[i] = [2]int{a, b}
       g[a] = append(g[a], b)
       grev[b] = append(grev[b], a)
   }
   // 1) Kosaraju order
   visited := make([]bool, n)
   order := make([]int, 0, n)
   // iterative DFS1
   for i := 0; i < n; i++ {
       if visited[i] {
           continue
       }
       // stack of (v, next index)
       var stk []struct{v, idx int}
       stk = append(stk, struct{v, idx int}{i, 0})
       for len(stk) > 0 {
           top := &stk[len(stk)-1]
           v := top.v
           if !visited[v] {
               visited[v] = true
           }
           if top.idx < len(g[v]) {
               w := g[v][top.idx]
               top.idx++
               if !visited[w] {
                   stk = append(stk, struct{v, idx int}{w, 0})
               }
           } else {
               // done
               order = append(order, v)
               stk = stk[:len(stk)-1]
           }
       }
   }
   // 2) assign components
   comp := make([]int, n)
   for i := range comp {
       comp[i] = -1
   }
   cid := 0
   for oi := n-1; oi >= 0; oi-- {
       v := order[oi]
       if comp[v] != -1 {
           continue
       }
       // BFS/DFS on grev
       var stack []int
       stack = append(stack, v)
       comp[v] = cid
       for len(stack) > 0 {
           u := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           for _, w := range grev[u] {
               if comp[w] == -1 {
                   comp[w] = cid
                   stack = append(stack, w)
               }
           }
       }
       cid++
   }
   // build comp graph edges
   C := cid
   pairs := make([][2]int, 0, m)
   for _, e := range edges {
       u := comp[e[0]]
       v := comp[e[1]]
       if u != v {
           pairs = append(pairs, [2]int{u, v})
       }
   }
   // dedup pairs
   // sort by u,v
   // simple radix: use map per u
   adj := make([][]int, C)
   seen := make(map[int]struct{})
   for _, p := range pairs {
       key := p[0]<<32 | p[1]
       if _, ok := seen[key]; ok {
           continue
       }
       seen[key] = struct{}{}
       adj[p[0]] = append(adj[p[0]], p[1])
   }
   // record component sizes
   compSize := make([]int, C)
   for i := 0; i < n; i++ {
       compSize[comp[i]]++
   }
   // compute answer: pipes inside SCCs to make them strongly connected
   ans := 0
   for i := 0; i < C; i++ {
       if compSize[i] > 1 {
           ans += compSize[i]
       }
   }
   // transitive reduction on comp graph: count edges needed
   // visited stamping
   visitedTime := make([]int, C)
   stamp := 1
   // for each u, process its outgoing edges
   for u := 0; u < C; u++ {
       ds := adj[u]
       k := len(ds)
       if k <= 1 {
           ans += k
           continue
       }
       // BFS from u through paths length>=2
       // mark all nodes reachable via children-of-children etc.
       // initial seeds: neighbors of children
       // build a small queue
       q := make([]int, 0, 16)
       // mark via stamp
       // increase stamp
       stamp++
       // avoid overflow
       if stamp == 0 {
           stamp = 1
           for i := range visitedTime {
               visitedTime[i] = 0
           }
       }
       // initial
       for _, w := range ds {
           for _, x := range adj[w] {
               if visitedTime[x] != stamp {
                   visitedTime[x] = stamp
                   q = append(q, x)
               }
           }
       }
       // BFS
       for i := 0; i < len(q); i++ {
           v := q[i]
           for _, y := range adj[v] {
               if visitedTime[y] != stamp {
                   visitedTime[y] = stamp
                   q = append(q, y)
               }
           }
       }
       // count edges needed
       for _, v := range ds {
           if visitedTime[v] != stamp {
               ans++
           }
       }
   }
   // print
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, ans)
}
