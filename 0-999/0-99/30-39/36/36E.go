package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents an edge in the multigraph
type Edge struct {
   to, id, rev int
   used        bool
}

var (
   adj  [][]*Edge
   ptr  []int
   path []int
)

// dfs performs Hierholzer's algorithm from vertex u
func dfs(u int) {
   for ptr[u] < len(adj[u]) {
       e := adj[u][ptr[u]]
       ptr[u]++
       if e.used {
           continue
       }
       e.used = true
       adj[e.to][e.rev].used = true
       dfs(e.to)
       path = append(path, e.id)
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var m int
   if _, err := fmt.Fscan(in, &m); err != nil {
       return
   }
   if m < 2 {
       fmt.Fprintln(out, -1)
       return
   }
   a := make([]int, m+1)
   b := make([]int, m+1)
   maxv := 0
   for i := 1; i <= m; i++ {
       fmt.Fscan(in, &a[i], &b[i])
       if a[i] > maxv {
           maxv = a[i]
       }
       if b[i] > maxv {
           maxv = b[i]
       }
   }
   adj = make([][]*Edge, maxv+1)
   // Build graph
   for i := 1; i <= m; i++ {
       u, v := a[i], b[i]
       eu := &Edge{to: v, id: i, rev: len(adj[v])}
       ev := &Edge{to: u, id: i, rev: len(adj[u])}
       adj[u] = append(adj[u], eu)
       adj[v] = append(adj[v], ev)
   }
   // Detect connected components (vertices with edges)
   visitedV := make([]bool, maxv+1)
   var compsVerts [][]int
   var compsEdges [][]int
   for v := 1; v <= maxv; v++ {
       if !visitedV[v] && len(adj[v]) > 0 {
           // BFS/DFS to collect component vertices
           stack := []int{v}
           visitedV[v] = true
           var compV []int
           for len(stack) > 0 {
               u := stack[len(stack)-1]
               stack = stack[:len(stack)-1]
               compV = append(compV, u)
               for _, e := range adj[u] {
                   if !visitedV[e.to] {
                       visitedV[e.to] = true
                       stack = append(stack, e.to)
                   }
               }
           }
           // Collect unique edge IDs in component
           usedE := make([]bool, m+1)
           var compE []int
           for _, u2 := range compV {
               for _, e := range adj[u2] {
                   if !usedE[e.id] {
                       usedE[e.id] = true
                       compE = append(compE, e.id)
                   }
               }
           }
           compsVerts = append(compsVerts, compV)
           compsEdges = append(compsEdges, compE)
       }
   }
   // If more than 2 components, impossible
   if len(compsVerts) > 2 {
       fmt.Fprintln(out, -1)
       return
   }
   // If exactly 2 components, build separate trails
   if len(compsVerts) == 2 {
       // Process each component separately
       var trails [][]int
       for ci := 0; ci < 2; ci++ {
           compV := compsVerts[ci]
           compE := compsEdges[ci]
           // Build local adjacency for component
           localAdj := make([][]*Edge, maxv+1)
           for _, id := range compE {
               u, v := a[id], b[id]
               eu := &Edge{to: v, id: id, rev: len(localAdj[v])}
               ev := &Edge{to: u, id: id, rev: len(localAdj[u])}
               localAdj[u] = append(localAdj[u], eu)
               localAdj[v] = append(localAdj[v], ev)
           }
           // Find local odd vertices
           var localOdds []int
           for _, u := range compV {
               if len(localAdj[u])%2 == 1 {
                   localOdds = append(localOdds, u)
               }
           }
           if len(localOdds) > 2 {
               fmt.Fprintln(out, -1)
               return
           }
           // Determine start vertex
           startLocal := compV[0]
           if len(localOdds) == 2 {
               startLocal = localOdds[0]
           }
           // Eulerian trail for component
           ptrComp := make([]int, maxv+1)
           var pathComp []int
           var dfsComp func(int)
           dfsComp = func(u int) {
               for ptrComp[u] < len(localAdj[u]) {
                   e := localAdj[u][ptrComp[u]]
                   ptrComp[u]++
                   if e.used {
                       continue
                   }
                   e.used = true
                   localAdj[e.to][e.rev].used = true
                   dfsComp(e.to)
                   pathComp = append(pathComp, e.id)
               }
           }
           dfsComp(startLocal)
           if len(pathComp) != len(compE) {
               fmt.Fprintln(out, -1)
               return
           }
           // Reverse pathComp
           for i, j := 0, len(pathComp)-1; i < j; i, j = i+1, j-1 {
               pathComp[i], pathComp[j] = pathComp[j], pathComp[i]
           }
           trails = append(trails, pathComp)
       }
       // Output two trails
       // First trail
       fmt.Fprintln(out, len(trails[0]))
       for i, id := range trails[0] {
           if i > 0 {
               fmt.Fprint(out, " ")
           }
           fmt.Fprint(out, id)
       }
       fmt.Fprintln(out)
       // Second trail
       fmt.Fprintln(out, len(trails[1]))
       for i, id := range trails[1] {
           if i > 0 {
               fmt.Fprint(out, " ")
           }
           fmt.Fprint(out, id)
       }
       fmt.Fprintln(out)
       return
   }
   // Find odd degree vertices
   odds := make([]int, 0, 4)
   for v := 1; v <= maxv; v++ {
       if len(adj[v])%2 == 1 {
           odds = append(odds, v)
       }
   }
   if len(odds) > 4 {
       fmt.Fprintln(out, -1)
       return
   }
   dummyID := m + 1
   dummyAdded := false
   // If 4 odd vertices, connect first two with dummy edge
   if len(odds) == 4 {
       u, v := odds[0], odds[1]
       eu := &Edge{to: v, id: dummyID, rev: len(adj[v])}
       ev := &Edge{to: u, id: dummyID, rev: len(adj[u])}
       adj[u] = append(adj[u], eu)
       adj[v] = append(adj[v], ev)
       dummyAdded = true
   }
   // Determine start vertex
   start := -1
   if dummyAdded {
       // After adding dummy, odds[2] and odds[3] remain odd
       start = odds[2]
   } else if len(odds) == 2 {
       start = odds[0]
   } else {
       // Eulerian circuit: start at any vertex with edges
       for v := 1; v <= maxv; v++ {
           if len(adj[v]) > 0 {
               start = v
               break
           }
       }
   }
   if start == -1 {
       fmt.Fprintln(out, -1)
       return
   }
   ptr = make([]int, maxv+1)
   dfs(start)
   totalEdges := m
   if dummyAdded {
       totalEdges++
   }
   if len(path) != totalEdges {
       // Not all edges are in one component or error
       fmt.Fprintln(out, -1)
       return
   }
   // Reverse path to correct order
   for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
       path[i], path[j] = path[j], path[i]
   }
   if dummyAdded {
       // Split at dummy edge
       var idx int
       for i, id := range path {
           if id == dummyID {
               idx = i
               break
           }
       }
       // First trail: edges before dummy
       L1 := idx
       // Second trail: edges after dummy
       L2 := len(path) - idx - 1
       if L1 == 0 || L2 == 0 {
           fmt.Fprintln(out, -1)
           return
       }
       fmt.Fprintln(out, L1)
       for i := 0; i < idx; i++ {
           if i > 0 {
               fmt.Fprint(out, " ")
           }
           fmt.Fprint(out, path[i])
       }
       fmt.Fprintln(out)
       fmt.Fprintln(out, L2)
       for i := idx + 1; i < len(path); i++ {
           if i > idx+1 {
               fmt.Fprint(out, " ")
           }
           fmt.Fprint(out, path[i])
       }
       fmt.Fprintln(out)
   } else {
       // No dummy: simple split into two non-empty parts
       // m >= 2 ensured
       fmt.Fprintln(out, 1)
       fmt.Fprintln(out, path[0])
       fmt.Fprintln(out, m-1)
       for i := 1; i < len(path); i++ {
           if i > 1 {
               fmt.Fprint(out, " ")
           }
           fmt.Fprint(out, path[i])
       }
       fmt.Fprintln(out)
   }
}
