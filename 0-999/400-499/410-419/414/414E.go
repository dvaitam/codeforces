package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   children [][]int
   parent []int
   depth  []int
)

func recomputeDepth() {
   // BFS from root=1
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   depth[1] = 0
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, v := range children[u] {
           depth[v] = depth[u] + 1
           queue = append(queue, v)
       }
   }
}

func lca(u, v int) int {
   // bring to same depth
   if depth[u] < depth[v] {
       u, v = v, u
   }
   for depth[u] > depth[v] {
       u = parent[u]
   }
   for u != v {
       u = parent[u]
       v = parent[v]
   }
   return u
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   children = make([][]int, n+1)
   parent = make([]int, n+1)
   depth = make([]int, n+1)
   for i := 1; i <= n; i++ {
       var li int
       fmt.Fscan(reader, &li)
       if li > 0 {
           children[i] = make([]int, li)
           for j := 0; j < li; j++ {
               fmt.Fscan(reader, &children[i][j])
               parent[children[i][j]] = i
           }
       } else {
           children[i] = []int{}
       }
   }
   // root has no parent
   parent[1] = 0
   recomputeDepth()
   for qi := 0; qi < m; qi++ {
       var typ int
       fmt.Fscan(reader, &typ)
       if typ == 1 {
           var u, v int
           fmt.Fscan(reader, &v, &u)
           w := lca(u, v)
           dist := depth[u] + depth[v] - 2*depth[w]
           fmt.Fprintln(writer, dist)
       } else if typ == 2 {
           var v, h int
           fmt.Fscan(reader, &v, &h)
           // find new parent: ancestor at distance h above v
           p := v
           for i := 0; i < h; i++ {
               p = parent[p]
           }
           // remove v from old parent
           old := parent[v]
           // find v in children[old]
           lst := children[old]
           for i, x := range lst {
               if x == v {
                   // remove index i
                   children[old] = append(lst[:i], lst[i+1:]...)
                   break
               }
           }
           // add to new parent
           children[p] = append(children[p], v)
           parent[v] = p
           // recompute depths
           recomputeDepth()
       } else if typ == 3 {
           var k int
           fmt.Fscan(reader, &k)
           // DFS preorder, track latest with depth k
           ans := -1
           var dfs func(u int)
           dfs = func(u int) {
               if depth[u] == k {
                   ans = u
               }
               for _, v := range children[u] {
                   dfs(v)
               }
           }
           dfs(1)
           fmt.Fprintln(writer, ans)
       }
   }
}
