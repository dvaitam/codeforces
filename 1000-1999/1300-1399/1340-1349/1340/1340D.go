package main

import (
   "bufio"
   "fmt"
   "os"
)

// pair represents a node and its associated value
type pair struct { x, y int }

var (
   to  [][]int
   ans []pair
)

// dfs performs the traversal analogous to the C++ solution
func dfs(x, f, v int) {
   ans = append(ans, pair{x, v})
   d := len(to[x])
   if f == 0 {
       for _, i := range to[x] {
           v++
           dfs(i, x, v)
           ans = append(ans, pair{x, v})
       }
   } else {
       ov := v
       if v >= d {
           ans = append(ans, pair{x, v-d})
           v -= d
           for _, i := range to[x] {
               if i == f {
                   continue
               }
               v++
               dfs(i, x, v)
               ans = append(ans, pair{x, v})
           }
       } else {
           for _, i := range to[x] {
               if i == f {
                   continue
               }
               if v >= d {
                   ans = append(ans, pair{x, 0})
                   v = 0
               }
               v++
               dfs(i, x, v)
               ans = append(ans, pair{x, v})
           }
       }
       if v >= ov {
           ans = append(ans, pair{x, ov-1})
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   to = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       to[u] = append(to[u], v)
       to[v] = append(to[v], u)
   }
   dfs(1, 0, 0)
   fmt.Fprintln(out, len(ans))
   for _, p := range ans {
       fmt.Fprintln(out, p.x, p.y)
   }
}
