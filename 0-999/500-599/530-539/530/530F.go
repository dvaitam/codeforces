package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int) int {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N int
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   // points[0] is origin
   xs := make([]int, N+1)
   ys := make([]int, N+1)
   xs[0], ys[0] = 0, 0
   for i := 1; i <= N; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
   }
   reachable := make([]bool, N+1)
   // for each species i = 1..10
   for i := 1; i <= 10; i++ {
       visited := make([]bool, N+1)
       q := make([]int, 0, N+1)
       visited[0] = true
       q = append(q, 0)
       for head := 0; head < len(q); head++ {
           u := q[head]
           ux, uy := xs[u], ys[u]
           for v := 1; v <= N; v++ {
               if visited[v] {
                   continue
               }
               vx, vy := xs[v], ys[v]
               // aligned and distance == i
               if ux == vx && abs(uy-vy) == i || uy == vy && abs(ux-vx) == i {
                   visited[v] = true
                   q = append(q, v)
               }
           }
       }
       for v := 1; v <= N; v++ {
           if visited[v] {
               reachable[v] = true
           }
       }
   }
   maxDist := 0
   for i := 1; i <= N; i++ {
       if reachable[i] {
           d := abs(xs[i]) + abs(ys[i])
           if d > maxDist {
               maxDist = d
           }
       }
   }
   fmt.Fprint(writer, maxDist)
}
