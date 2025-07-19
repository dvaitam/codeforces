package main

import (
   "bufio"
   "fmt"
   "os"
)

const N = 410000

var (
   adj  [N][]int
   path []int
   act  [N]bool
   seen [N]int
   nx   [N]int
   pr   [N]int
   typ  int
)

func dfs(u, bad int) int {
   if seen[u] == typ {
       return 0
   }
   seen[u] = typ
   if u == bad {
       return 0
   }
   if act[u] {
       return u
   }
   path = append(path, u)
   for _, v := range adj[u] {
       rec := dfs(v, bad)
       if rec != 0 {
           return rec
       }
   }
   path = path[:len(path)-1]
   return 0
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   // initialize globals
   path = make([]int, 0, N)
   typ = 17

   for T > 0 {
       T--
       var n, m, s, t int
       fmt.Fscan(reader, &n, &m, &s, &t)
       // clear adjacency
       for i := 1; i <= n; i++ {
           adj[i] = adj[i][:0]
       }
       // read edges
       edges := make([][2]int, m)
       for i := 0; i < m; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
           edges[i][0], edges[i][1] = u, v
       }
       // reset state
       for i := 1; i <= n; i++ {
           nx[i] = -1
           pr[i] = -1
           act[i] = false
       }
       nx[s] = t
       pr[t] = s
       pr[s] = -1
       nx[t] = -1
       act[s], act[t] = true, true
       // build path
       ok := true
       u := s
       for ok && u != t {
           for _, w := range adj[u] {
               if act[w] {
                   continue
               }
               // find path from w to an active node
               path = path[:0]
               typ++
               v := dfs(w, u)
               if v == 0 {
                   ok = false
                   break
               }
               // insert the path nodes between u and its next
               for len(path) > 0 {
                   x := path[len(path)-1]
                   path = path[:len(path)-1]
                   act[x] = true
                   // relink pointers: u->x->oldNext
                   pr[nx[u]] = x
                   nx[x] = nx[u]
                   pr[x] = u
                   nx[u] = x
               }
           }
           if !ok {
               break
           }
           u = nx[u]
       }
       if !ok {
           fmt.Fprintln(writer, "No")
           continue
       }
       // check all nodes are in the chain
       for i := 1; i <= n; i++ {
           if nx[i] == -1 && pr[i] == -1 {
               ok = false
               break
           }
       }
       if !ok {
           fmt.Fprintln(writer, "No")
           continue
       }
       // collect order
       ans := make([]int, 0, n)
       for x := s; x != -1; x = nx[x] {
           ans = append(ans, x)
       }
       rans := make([]int, n+1)
       for idx, x := range ans {
           if idx >= n {
               break
           }
           rans[x] = idx
       }
       // output
       fmt.Fprintln(writer, "Yes")
       for _, e := range edges {
           u1, v1 := e[0], e[1]
           if rans[u1] > rans[v1] {
               u1, v1 = v1, u1
           }
           fmt.Fprintf(writer, "%d %d\n", u1, v1)
       }
   }
}
