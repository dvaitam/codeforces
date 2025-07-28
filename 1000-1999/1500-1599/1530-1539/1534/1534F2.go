package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([][]byte, n)
   for i := 0; i < n; i++ {
       var row string
       fmt.Fscan(reader, &row)
       grid[i] = []byte(row)
   }
   a := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // assign id to each sand block
   id := make([][]int, n)
   for i := range id {
       id[i] = make([]int, m)
       for j := range id[i] {
           id[i][j] = -1
       }
   }
   // rows and cols of each node
   rows := make([]int, 0)
   cols := make([]int, 0)
   idx := 0
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '#' {
               id[i][j] = idx
               rows = append(rows, i)
               cols = append(cols, j)
               idx++
           }
       }
   }
   tot := idx
   // prepare helper arrays
   firstBelow := make([][]int, m)
   for j := 0; j < m; j++ {
       firstBelow[j] = make([]int, n)
       nearest := -1
       for i := n - 1; i >= 0; i-- {
           if id[i][j] != -1 {
               nearest = i
           }
           firstBelow[j][i] = nearest
       }
   }
   belowID := make([]int, tot)
   for j := 0; j < m; j++ {
       last := -1
       for i := n - 1; i >= 0; i-- {
           if id[i][j] != -1 {
               cur := id[i][j]
               belowID[cur] = last
               last = cur
           }
       }
   }
   // build directed graph
   adj := make([][]int, tot)
   rev := make([][]int, tot)
   addEdge := func(u, v int) {
       adj[u] = append(adj[u], v)
       rev[v] = append(rev[v], u)
   }
   for u := 0; u < tot; u++ {
       i, j := rows[u], cols[u]
       // up
       if i > 0 && id[i-1][j] != -1 {
           addEdge(u, id[i-1][j])
       }
       // down in same col
       if b := belowID[u]; b != -1 {
           addEdge(u, b)
       }
       // left neighbor via falling path
       if j > 0 {
           if r := firstBelow[j-1][i]; r != -1 {
               addEdge(u, id[r][j-1])
           }
       }
       // right neighbor
       if j+1 < m {
           if r := firstBelow[j+1][i]; r != -1 {
               addEdge(u, id[r][j+1])
           }
       }
   }
   // kosaraju for SCC
   visited := make([]bool, tot)
   order := make([]int, 0, tot)
   var dfs1 func(int)
   dfs1 = func(u int) {
       visited[u] = true
       for _, v := range adj[u] {
           if !visited[v] {
               dfs1(v)
           }
       }
       order = append(order, u)
   }
   for u := 0; u < tot; u++ {
       if !visited[u] {
           dfs1(u)
       }
   }
   comp := make([]int, tot)
   for i := range comp {
       comp[i] = -1
   }
   cnum := 0
   var dfs2 func(int)
   dfs2 = func(u int) {
       comp[u] = cnum
       for _, v := range rev[u] {
           if comp[v] == -1 {
               dfs2(v)
           }
       }
   }
   for i := tot - 1; i >= 0; i-- {
       u := order[i]
       if comp[u] == -1 {
           dfs2(u)
           cnum++
       }
   }
   if cnum == 0 {
       fmt.Fprintln(writer, 0)
       return
   }
   // build condensed DAG
   cadj := make([][]int, cnum)
   indeg := make([]int, cnum)
   for u := 0; u < tot; u++ {
       cu := comp[u]
       for _, v := range adj[u] {
           cv := comp[v]
           if cu != cv {
               cadj[cu] = append(cadj[cu], cv)
           }
       }
   }
   // topological sort on condensed DAG
   for u := 0; u < cnum; u++ {
       for _, v := range cadj[u] {
           indeg[v]++
       }
   }
   topo := make([]int, 0, cnum)
   q := make([]int, 0, cnum)
   for u := 0; u < cnum; u++ {
       if indeg[u] == 0 {
           q = append(q, u)
       }
   }
   for bi := 0; bi < len(q); bi++ {
       u := q[bi]
       topo = append(topo, u)
       for _, v := range cadj[u] {
           indeg[v]--
           if indeg[v] == 0 {
               q = append(q, v)
           }
       }
   }
   // identify target components
   target := make([]bool, cnum)
   for j := 0; j < m; j++ {
       need := a[j]
       if need <= 0 {
           continue
       }
       for i := n - 1; i >= 0 && need > 0; i-- {
           if id[i][j] != -1 {
               need--
               target[comp[id[i][j]]] = true
           }
       }
   }
   // greedy covering using topo order and BFS
   seen := make([]bool, cnum)
   covered := make([]bool, cnum)
   ans := 0
   for _, u := range topo {
       if !target[u] || covered[u] {
           continue
       }
       // need new operation
       ans++
       // bfs from u
       queue := []int{u}
       seen[u] = true
       if target[u] {
           covered[u] = true
       }
       for bi := 0; bi < len(queue); bi++ {
           x := queue[bi]
           for _, v := range cadj[x] {
               if !seen[v] {
                   seen[v] = true
                   if target[v] {
                       covered[v] = true
                   }
                   queue = append(queue, v)
               }
           }
       }
   }
   fmt.Fprintln(writer, ans)
}
