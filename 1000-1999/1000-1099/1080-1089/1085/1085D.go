package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var s float64
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   // count degrees
   deg := make([]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       deg[u]++
       deg[v]++
   }
   // count leaves (degree == 1)
   leaf := 0
   for i := 1; i <= n; i++ {
       if deg[i] == 1 {
           leaf++
       }
   }
   if leaf == 0 {
       // degenerate: no leaf? but n>=2, so not possible
       leaf = 1
   }
   // minimal diameter = 2*s / leaf
   ans := 2.0 * s / float64(leaf)
   // print with precision
   fmt.Printf("%.10f\n", ans)
}
