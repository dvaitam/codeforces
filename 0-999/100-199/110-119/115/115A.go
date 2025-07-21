package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n)
   for i := 0; i < n; i++ {
       // read manager index, convert to 0-based; -1 stays as -1
       var x int
       fmt.Fscan(reader, &x)
       if x == -1 {
           p[i] = -1
       } else {
           p[i] = x - 1
       }
   }
   // depth array, 0 means uncomputed
   depth := make([]int, n)
   var maxDepth int

   // compute depth of i (1-based depth)
   var dfs func(int) int
   dfs = func(u int) int {
       if depth[u] != 0 {
           return depth[u]
       }
       if p[u] < 0 {
           depth[u] = 1
       } else {
           depth[u] = dfs(p[u]) + 1
       }
       return depth[u]
   }

   for i := 0; i < n; i++ {
       d := dfs(i)
       if d > maxDepth {
           maxDepth = d
       }
   }
   fmt.Println(maxDepth)
}
