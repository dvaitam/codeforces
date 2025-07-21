package main

import (
   "bufio"
   "fmt"
   "os"
)

func bfs(n int, g [][]int, start int) []int {
   dist := make([]int, n)
   for i := 0; i < n; i++ {
       dist[i] = -1
   }
   queue := make([]int, 0, n)
   queue = append(queue, start)
   dist[start] = 0
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       for _, v := range g[u] {
           if dist[v] == -1 {
               dist[v] = dist[u] + 1
               queue = append(queue, v)
           }
       }
   }
   return dist
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, d int
   if _, err := fmt.Fscan(in, &n, &m, &d); err != nil {
       return
   }
   p := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &p[i])
       p[i]--
   }
   g := make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       g[u] = append(g[u], v)
       g[v] = append(g[v], u)
   }
   // find one endpoint a of diameter among p
   d0 := bfs(n, g, p[0])
   a := p[0]
   for _, x := range p {
       if d0[x] > d0[a] {
           a = x
       }
   }
   // BFS from a
   da := bfs(n, g, a)
   // find other endpoint b
   b := p[0]
   for _, x := range p {
       if da[x] > da[b] {
           b = x
       }
   }
   // BFS from b
   db := bfs(n, g, b)

   // count nodes where max(dist to a, dist to b) <= d
   cnt := 0
   for i := 0; i < n; i++ {
       if da[i] <= d && db[i] <= d {
           cnt++
       }
   }
   fmt.Println(cnt)
}
