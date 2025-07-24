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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   // BFS to build tree parent and depth, and order
   parent := make([]int, n+1)
   depth := make([]int, n+1)
   queue := make([]int, n)
   var bfsList []int
   head, tail := 0, 0
   queue[tail] = 1
   tail++
   parent[1] = 0
   depth[1] = 0
   bfsList = append(bfsList, 1)
   for head < tail {
       u := queue[head]
       head++
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           depth[v] = depth[u] + 1
           queue[tail] = v
           tail++
           bfsList = append(bfsList, v)
       }
   }
   leafCount := make([]int, n+1)
   maxDepthLeaf := 0
   ans := 0
   // process bottom-up
   for i := len(bfsList) - 1; i >= 0; i-- {
       u := bfsList[i]
       // leaf (except root)
       if u != 1 && len(adj[u]) == 1 {
           leafCount[u] = 1
           if depth[u] > maxDepthLeaf {
               maxDepthLeaf = depth[u]
           }
       } else {
           sum := 0
           for _, v := range adj[u] {
               if v == parent[u] {
                   continue
               }
               sum += leafCount[v]
           }
           leafCount[u] = sum
       }
       // for non-root with more than one leaf in subtree, candidate time
       if u != 1 && leafCount[u] > 1 {
           t := depth[u] + leafCount[u]
           if t > ans {
               ans = t
           }
       }
   }
   if maxDepthLeaf > ans {
       ans = maxDepthLeaf
   }
   fmt.Fprintln(writer, ans)
}
