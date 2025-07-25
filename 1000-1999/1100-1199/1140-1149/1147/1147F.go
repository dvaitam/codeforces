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
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, n)
       for j := 0; j < n; j++ {
           fmt.Fscan(in, &a[i][j])
       }
   }
   // build functional graph: for left nodes pick max edge, for right pick min edge
   m := 2 * n
   f := make([]int, m)
   for i := 0; i < n; i++ {
       bestJ := 0
       for j := 1; j < n; j++ {
           if a[i][j] > a[i][bestJ] {
               bestJ = j
           }
       }
       f[i] = bestJ + n
   }
   for j := 0; j < n; j++ {
       bestI := 0
       for i := 1; i < n; i++ {
           if a[i][j] < a[bestI][j] {
               bestI = i
           }
       }
       f[j+n] = bestI
   }
   // find a cycle
   vis := make([]int, m)
   var cycle []int
   for i := 0; i < m; i++ {
       if vis[i] != 0 {
           continue
       }
       path := make(map[int]int)
       u := i
       step := 1
       for {
           if vis[u] != 0 {
               break
           }
           vis[u] = step
           path[u] = step
           step++
           u = f[u]
       }
       if startStep, ok := path[u]; ok {
           // extract cycle
           v := u
           for {
               cycle = append(cycle, v)
               v = f[v]
               if v == u {
                   break
               }
           }
           break
       }
   }
   // decide strategy
   // cycle length is even
   inc := true
   if len(cycle)%4 == 2 {
       inc = false
   }
   // starting node
   start := cycle[0] + 1
   // output choice
   if inc {
       fmt.Fprintf(out, "increasing %d\n", start)
   } else {
       fmt.Fprintf(out, "decreasing %d\n", start)
   }
   out.Flush()
   // interact: judge moves first
   for {
       var v int
       if _, err := fmt.Fscan(in, &v); err != nil || v < 0 {
           return
       }
       // move along functional edge
       to := f[v-1] + 1
       fmt.Fprintf(out, "%d\n", to)
       out.Flush()
   }
}
