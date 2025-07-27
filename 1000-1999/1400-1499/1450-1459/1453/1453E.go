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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       adj := make([][]int, n+1)
       for i := 0; i < n-1; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       // BFS to get parent and order
       parent := make([]int, n+1)
       order := make([]int, 0, n)
       queue := make([]int, 0, n)
       queue = append(queue, 1)
       parent[1] = -1
       for i := 0; i < len(queue); i++ {
           u := queue[i]
           order = append(order, u)
           for _, v := range adj[u] {
               if parent[v] == 0 {
                   parent[v] = u
                   queue = append(queue, v)
               }
           }
       }
       // DP arrays
       h := make([]int, n+1)
       ans := make([]int, n+1)
       // process in reverse BFS order
       for i := len(order) - 1; i >= 0; i-- {
           u := order[i]
           // gather child heights and ans
           maxH := 0
           secondH := 0
           bestAns := 0
           children := 0
           for _, v := range adj[u] {
               if parent[v] == u {
                   // child
                   children++
                   // update bestAns
                   if ans[v] > bestAns {
                       bestAns = ans[v]
                   }
                   hv := h[v] + 1
                   // update top two heights
                   if hv > maxH {
                       secondH = maxH
                       maxH = hv
                   } else if hv > secondH {
                       secondH = hv
                   }
               }
           }
           curAns := bestAns
           if children >= 2 {
               // transition between subtrees: second largest height + 2
               if secondH+2 > curAns {
                   curAns = secondH + 2
               }
           }
           if children >= 1 {
               // return to this node after last subtree: largest height + 1
               if maxH+1 > curAns {
                   curAns = maxH + 1
               }
           }
           ans[u] = curAns
           // compute h[u]
           if children >= 1 {
               h[u] = maxH
           } else {
               h[u] = 0
           }
       }
       // result at root
       fmt.Fprintln(writer, ans[1])
   }
}
