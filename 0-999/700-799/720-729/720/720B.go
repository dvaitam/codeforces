package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   scanner := bufio.NewScanner(reader)
   scanner.Split(bufio.ScanWords)
   next := func() string {
       scanner.Scan()
       return scanner.Text()
   }
   // read n, m
   n, _ := strconv.Atoi(next())
   m, _ := strconv.Atoi(next())
   u := make([]int, m)
   v := make([]int, m)
   c := make([]int, m)
   for i := 0; i < m; i++ {
       u[i], _ = strconv.Atoi(next())
       v[i], _ = strconv.Atoi(next())
       c[i], _ = strconv.Atoi(next())
       u[i]--
       v[i]--
   }
   // build adjacency
   adj := make([][]struct{to, eid int}, n)
   for i := 0; i < m; i++ {
       adj[u[i]] = append(adj[u[i]], struct{to, eid int}{v[i], i})
       adj[v[i]] = append(adj[v[i]], struct{to, eid int}{u[i], i})
   }
   visited := make([]bool, n)
   parent := make([]int, n)
   parentEdge := make([]int, n)
   depth := make([]int, n)
   var cycles [][]int
   var dfs func(int, int)
   dfs = func(x, p int) {
       visited[x] = true
       for _, e := range adj[x] {
           y := e.to; eid := e.eid
           if y == p {
               continue
           }
           if !visited[y] {
               parent[y] = x
               parentEdge[y] = eid
               depth[y] = depth[x] + 1
               dfs(y, x)
           } else if depth[y] < depth[x] {
               // found back-edge x->y
               cycle := []int{eid}
               w := x
               for w != y {
                   cycle = append(cycle, parentEdge[w])
                   w = parent[w]
               }
               cycles = append(cycles, cycle)
           }
       }
   }
   depth[0] = 0
   dfs(0, -1)
   // count colors occurrences
   colorCount := make(map[int]int)
   for i := 0; i < m; i++ {
       colorCount[c[i]]++
   }
   // mark removed edges
   removed := make([]bool, m)
   for _, cycle := range cycles {
       chosen := -1
       for _, eid := range cycle {
           if colorCount[c[eid]] > 1 {
               chosen = eid
               break
           }
       }
       if chosen == -1 {
           chosen = cycle[0]
       }
       removed[chosen] = true
   }
   // count distinct colors in resulting tree
   outColor := make(map[int]struct{})
   for i := 0; i < m; i++ {
       if !removed[i] {
           outColor[c[i]] = struct{}{}
       }
   }
   fmt.Println(len(outColor))
}
