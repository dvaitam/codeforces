package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n       int
   con     [][]int
   sz      []int
   rt      int
   last    int
   curkey  int
   cmds    [][3]int
)

// dfs computes subtree sizes
func dfs(x, pre int) {
   sz[x] = 1
   for _, y := range con[x] {
       if y == pre {
           continue
       }
       dfs(y, x)
       sz[x] += sz[y]
   }
}

// dfs2 generates transformations for subtree rooted at x
func dfs2(x, pre int) {
   if pre != rt {
       cmds = append(cmds, [3]int{rt, last, x})
       cmds = append(cmds, [3]int{x, pre, curkey})
       last = x
   }
   for _, y := range con[x] {
       if y == pre {
           continue
       }
       dfs2(y, x)
   }
}

// solve processes node u, skipping neighbor v (0 means none)
func solve(u, v int) {
   rt = u
   for _, y := range con[u] {
       if y == v {
           continue
       }
       last = y
       curkey = y
       dfs2(y, u)
       cmds = append(cmds, [3]int{rt, last, curkey})
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n)
   con = make([][]int, n+1)
   sz = make([]int, n+1)
   for i := 0; i < n-1; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       con[x] = append(con[x], y)
       con[y] = append(con[y], x)
   }

   // initial dfs from node 1 to compute sizes
   dfs(1, 0)
   u := 1
   // find centroid u
   for i := 1; i <= n; i++ {
       if sz[i]*2 >= n && sz[i] < sz[u] {
           u = i
       }
   }
   // recompute sizes from centroid
   dfs(u, 0)
   v := 0
   // check if two centroids
   for _, y := range con[u] {
       if sz[y]*2 == n {
           v = y
           break
       }
   }
   // generate operations
   if v == 0 {
       solve(u, 0)
   } else {
       solve(u, v)
       solve(v, u)
   }
   // output
   fmt.Fprintln(writer, len(cmds))
   for _, op := range cmds {
       fmt.Fprintln(writer, op[0], op[1], op[2])
   }
}
