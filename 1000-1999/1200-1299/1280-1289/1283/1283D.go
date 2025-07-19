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
   v := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v[i])
   }
   // BFS from initial positions
   type node struct{ pos, dist int }
   queue := make([]node, 0, n+m)
   visited := make(map[int]struct{}, n+m)
   for _, a := range v {
       if _, ok := visited[a]; !ok {
           visited[a] = struct{}{}
           queue = append(queue, node{pos: a, dist: 0})
       }
   }
   ans := make([]int, 0, m)
   var sum int64
   // BFS
   for head := 0; head < len(queue) && len(ans) < m; head++ {
       cur := queue[head]
       for _, d := range []int{-1, 1} {
           np := cur.pos + d
           if _, ok := visited[np]; ok {
               continue
           }
           visited[np] = struct{}{}
           nd := cur.dist + 1
           queue = append(queue, node{pos: np, dist: nd})
           sum += int64(nd)
           ans = append(ans, np)
           if len(ans) == m {
               break
           }
       }
   }
   // output
   fmt.Fprintln(writer, sum)
   for i, x := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, x)
   }
   fmt.Fprintln(writer)
}
