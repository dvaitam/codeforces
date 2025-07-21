package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   adj := make([][]int, n+1)
   origDeg := make([]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // original degrees
   for u := 1; u <= n; u++ {
       origDeg[u] = len(adj[u])
   }
   // detect simple cycles via DFS (vertex cactus guarantees at most one cycle per vertex)
   type frame struct{ u, parent, next int }
   state := make([]int8, n+1)      // 0=unvisited,1=visiting,2=done
   posInStack := make([]int, n+1)  // position in pathStack or -1
   for i := range posInStack {
       posInStack[i] = -1
   }
   pathStack := make([]int, 0, n)
   cycleMark := make([]bool, n+1)
   cycles := 0
   // iterative DFS from node 1 (graph is connected)
   dfsStack := []frame{{u: 1, parent: -1, next: 0}}
   state[1] = 1
   pathStack = append(pathStack, 1)
   posInStack[1] = 0
   for len(dfsStack) > 0 {
       f := &dfsStack[len(dfsStack)-1]
       u := f.u
       if f.next < len(adj[u]) {
           v := adj[u][f.next]
           f.next++
           if v == f.parent {
               continue
           }
           if state[v] == 0 {
               // tree edge
               state[v] = 1
               posInStack[v] = len(pathStack)
               pathStack = append(pathStack, v)
               dfsStack = append(dfsStack, frame{u: v, parent: u, next: 0})
           } else if state[v] == 1 {
               // back-edge, found a cycle
               cycles++
               // mark nodes from v to u in pathStack
               for i := posInStack[v]; i < len(pathStack); i++ {
                   cycleMark[pathStack[i]] = true
               }
           }
       } else {
           // backtrack
           state[u] = 2
           dfsStack = dfsStack[:len(dfsStack)-1]
           pathStack = pathStack[:len(pathStack)-1]
           posInStack[u] = -1
       }
   }
   inCycle := cycleMark
   // compute forest odd degrees
   odd := 0
   for u := 1; u <= n; u++ {
       d := origDeg[u]
       if inCycle[u] {
           d -= 2
       }
       if d&1 == 1 {
           odd++
       }
   }
   minLines := cycles + odd/2
   // maximum is every edge as a radial line
   maxLines := m
   fmt.Printf("%d %d\n", minLines, maxLines)
}
