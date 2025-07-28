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
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       adj := make([][]int, n+1)
       for i := 0; i < n-1; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       parent := make([]int, n+1)
       children := make([][]int, n+1)
       // BFS to build rooted tree at 1
       queue := make([]int, n)
       head, tail := 0, 0
       queue[tail] = 1; tail++
       parent[1] = 0
       for head < tail {
           u := queue[head]; head++
           for _, v := range adj[u] {
               if v == parent[u] {
                   continue
               }
               parent[v] = u
               children[u] = append(children[u], v)
               queue[tail] = v; tail++
           }
       }
       // count leaves
       leaves := 0
       isLeaf := make([]bool, n+1)
       for u := 1; u <= n; u++ {
           if len(children[u]) == 0 {
               leaves++
               isLeaf[u] = true
           }
       }
       // count buds and their leaf-children
       buds := 0
       s := 0
       for u := 2; u <= n; u++ {
           ch := children[u]
           if len(ch) == 0 {
               continue
           }
           allLeaf := true
           for _, c := range ch {
               if !isLeaf[c] {
                   allLeaf = false
                   break
               }
           }
           if allLeaf {
               buds++
               s += len(ch)
           }
       }
       // minimal leaves = max(leaves - buds, s)
       without := leaves - buds
       if s > without {
           fmt.Fprintln(writer, s)
       } else {
           fmt.Fprintln(writer, without)
       }
   }
}
