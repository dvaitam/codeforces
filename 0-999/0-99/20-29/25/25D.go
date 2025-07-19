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
   for {
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       parent := make([]int, n+1)
       for i := 1; i <= n; i++ {
           parent[i] = i
       }
       var find func(int) int
       find = func(x int) int {
           if parent[x] != x {
               parent[x] = find(parent[x])
           }
           return parent[x]
       }
       bad := make([][2]int, 0, n)
       for i := 0; i < n-1; i++ {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           x0 := find(x)
           y0 := find(y)
           if x0 == y0 {
               bad = append(bad, [2]int{x, y})
           } else {
               parent[y0] = x0
           }
       }
       roots := make([]int, 0, len(bad)+1)
       for i := 1; i <= n; i++ {
           if parent[i] == i {
               roots = append(roots, i)
           }
       }
       fmt.Fprintln(writer, len(bad))
       for i := 0; i < len(bad); i++ {
           fmt.Fprintf(writer, "%d %d %d %d\n", bad[i][0], bad[i][1], roots[i], roots[i+1])
       }
   }
}
