package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   P := make([]int, m+1)
   for i := 0; i <= m; i++ {
       fmt.Fscan(in, &P[i])
   }
   start := P[0]
   // Build edges from path
   type edge struct{ u, v int }
   edges := make([]edge, m)
   // adjacency lists: for each node, list of (neighbor, edgeID)
   adj := make([][]struct{ v, id int }, n+1)
   for i := 0; i < m; i++ {
       u, v := P[i], P[i+1]
       edges[i] = edge{u, v}
       adj[u] = append(adj[u], struct{ v, id int }{v, i})
       adj[v] = append(adj[v], struct{ v, id int }{u, i})
   }
   // sort adjacency by neighbor
   for u := 1; u <= n; u++ {
       // simple insertion sort as deg small
       a := adj[u]
       for i := 1; i < len(a); i++ {
           x := a[i]
           j := i
           for j > 0 && a[j-1].v > x.v {
               a[j] = a[j-1]
               j--
           }
           a[j] = x
       }
       adj[u] = a
   }
   // compute original degrees
   origDeg := make([]int, n+1)
   for _, e := range edges {
       origDeg[e.u]++
       origDeg[e.v]++
   }
   // map edge to its use index in orig path
   useIndex := make([]int, m)
   for i := 0; i < m; i++ {
       useIndex[i] = i
   }
   // helper BFS connectivity
   var bfsReach = func(src int, removed []bool) []bool {
       vis := make([]bool, n+1)
       q := make([]int, 0, n)
       vis[src] = true
       q = append(q, src)
       for qi := 0; qi < len(q); qi++ {
           u := q[qi]
           for _, ne := range adj[u] {
               if removed[ne.id] || vis[ne.v] {
                   continue
               }
               vis[ne.v] = true
               q = append(q, ne.v)
           }
       }
       return vis
   }
   // try positions
   var result []int
   for pos := m - 1; pos >= 0; pos-- {
       u := P[pos]
       oldNext := P[pos+1]
       // find old edge index in adj[u]
       var k int
       for i, ne := range adj[u] {
           if ne.v == oldNext && ne.id == pos {
               k = i
               break
           }
       }
       // try next neighbors
       for j := k + 1; j < len(adj[u]); j++ {
           ne := adj[u][j]
           v := ne.v
           eid := ne.id
           // skip if edge used in prefix
           if useIndex[eid] < pos {
               continue
           }
           // compute remaining degrees
           remDeg := make([]int, n+1)
           copy(remDeg, origDeg)
           for i := 0; i < pos; i++ {
               e := edges[i]
               remDeg[e.u]--
               remDeg[e.v]--
           }
           remDeg[u]--
           remDeg[v]--
           // check odd degrees
           odd := make([]int, 0, 2)
           for x := 1; x <= n; x++ {
               if remDeg[x]&1 == 1 {
                   odd = append(odd, x)
               }
           }
           okOdd := false
           if len(odd) == 0 && v == start {
               okOdd = true
           } else if len(odd) == 2 && ((odd[0] == start && odd[1] == v) || (odd[1] == start && odd[0] == v)) {
               okOdd = true
           }
           if !okOdd {
               continue
           }
           // removed edges array
           removed := make([]bool, m)
           for i := 0; i < pos; i++ {
               removed[useIndex[i]] = true
           }
           removed[eid] = true
           // check connectivity
           var src = v
           if len(odd) == 0 {
               src = start
           }
           vis := bfsReach(src, removed)
           // ensure all remaining edges endpoints visited
           goodConn := true
           for e := 0; e < m; e++ {
               if removed[e] {
                   continue
               }
               if !vis[edges[e].u] || !vis[edges[e].v] {
                   goodConn = false
                   break
               }
           }
           if !goodConn {
               continue
           }
           // build new path prefix
           result = make([]int, 0, m+1)
           result = append(result, P[:pos+1]...)
           result = append(result, v)
           // now build rest by Fleury's algorithm
           usedTrail := make([]bool, m)
           copy(usedTrail, removed)
           cur := v
           for step := pos + 1; step < m; step++ {
               // gather available edges
               cand := make([]struct{ v, id int }, 0, len(adj[cur]))
               for _, e2 := range adj[cur] {
                   if !usedTrail[e2.id] {
                       cand = append(cand, e2)
                   }
               }
               var pick struct{ v, id int }
               if len(cand) == 1 {
                   pick = cand[0]
               } else {
                   // try in order
                   for _, e2 := range cand {
                       // test if bridge
                       usedTrail[e2.id] = true
                       vis2 := bfsReach(cur, usedTrail)
                       usedTrail[e2.id] = false
                       if vis2[e2.v] {
                           pick = e2
                           break
                       }
                   }
               }
               usedTrail[pick.id] = true
               cur = pick.v
               result = append(result, cur)
           }
           // complete circuit: result should be length m+1
           if len(result) == m+1 {
               out := bufio.NewWriter(os.Stdout)
               for i, x := range result {
                   if i > 0 {
                       out.WriteByte(' ')
                   }
                   fmt.Fprint(out, x)
               }
               out.WriteByte('\n')
               out.Flush()
               return
           }
       }
   }
   fmt.Println("No solution")
}
