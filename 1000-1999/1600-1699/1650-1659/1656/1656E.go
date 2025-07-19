package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(in, &n)
       G := make([][]int, n)
       deg := make([]int, n)
       for i := 1; i < n; i++ {
           var u, v int
           fmt.Fscan(in, &u, &v)
           u--
           v--
           G[u] = append(G[u], v)
           G[v] = append(G[v], u)
           deg[u]++
           deg[v]++
       }
       color := make([]int, n)
       visited := make([]bool, n)
       queue := make([]int, 0, n)
       queue = append(queue, 0)
       visited[0] = true
       // BFS to color tree bipartitely
       for head := 0; head < len(queue); head++ {
           u := queue[head]
           for _, v := range G[u] {
               if !visited[v] {
                   visited[v] = true
                   color[v] = 1 - color[u]
                   queue = append(queue, v)
               }
           }
       }
       for i := 0; i < n; i++ {
           val := deg[i]
           if color[i] == 0 {
               val = -val
           }
           fmt.Fprintf(out, "%d ", val)
       }
       fmt.Fprintln(out)
   }
}
