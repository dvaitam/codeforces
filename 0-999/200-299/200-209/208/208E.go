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

   var n int
   fmt.Fscan(in, &n)
   parent := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &parent[i])
   }
   // build children
   children := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       p := parent[i]
       if p > 0 {
           children[p] = append(children[p], i)
       }
   }
   // prepare depth, tin/tout, depthList
   depth := make([]int, n+1)
   tin := make([]int, n+1)
   tout := make([]int, n+1)
   depthList := make([][]int, n+2)
   t := 0
   var dfs func(int)
   dfs = func(u int) {
       t++
       tin[u] = t
       d := depth[u]
       depthList[d] = append(depthList[d], t)
       for _, v := range children[u] {
           depth[v] = d + 1
           dfs(v)
       }
       tout[u] = t
   }
   // run DFS from each root
   for i := 1; i <= n; i++ {
       if parent[i] == 0 {
           depth[i] = 0
           dfs(i)
       }
   }
   // binary lifting
   const LOG = 18
   up := make([][]int, LOG)
   up[0] = make([]int, n+1)
   for i := 1; i <= n; i++ {
       up[0][i] = parent[i]
   }
   for k := 1; k < LOG; k++ {
       up[k] = make([]int, n+1)
       for i := 1; i <= n; i++ {
           up[k][i] = up[k-1][ up[k-1][i] ]
       }
   }
   // process queries
   var m int
   fmt.Fscan(in, &m)
   res := make([]int, m)
   for qi := 0; qi < m; qi++ {
       var v, p int
       fmt.Fscan(in, &v, &p)
       // find p-th ancestor
       u := v
       for k := 0; k < LOG && u > 0; k++ {
           if (p>>k)&1 == 1 {
               u = up[k][u]
           }
       }
       if u == 0 {
           res[qi] = 0
           continue
       }
       // count descendants at depth depth[u]+p in subtree of u
       D := depth[u] + p
       if D >= len(depthList) {
           res[qi] = 0
           continue
       }
       arr := depthList[D]
       l := sort.Search(len(arr), func(i int) bool { return arr[i] >= tin[u] })
       r := sort.Search(len(arr), func(i int) bool { return arr[i] > tout[u] })
       cnt := r - l
       if cnt > 0 {
           cnt-- // exclude v itself
       }
       res[qi] = cnt
   }
   // output
   for i, v := range res {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, v)
   }
}
