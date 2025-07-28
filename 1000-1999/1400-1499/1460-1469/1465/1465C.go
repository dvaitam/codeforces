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
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       out := make([]int, n+1)
       k := 0
       for i := 0; i < m; i++ {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           if x != y {
               out[x] = y
               k++
           }
       }
       visited := make([]bool, n+1)
       cycles := 0
       // detect cycles in directed graph with outdeg <=1
       for i := 1; i <= n; i++ {
           if visited[i] || out[i] == 0 {
               visited[i] = true
               continue
           }
           // traverse path
           cur := i
           index := make(map[int]int)
           step := 0
           for {
               if visited[cur] {
                   break
               }
               if _, ok := index[cur]; ok {
                   cycles++
                   break
               }
               index[cur] = step
               step++
               next := out[cur]
               if next == 0 {
                   break
               }
               cur = next
           }
           // mark all in path as visited
           for u := range index {
               visited[u] = true
           }
       }
       // result is k moves for each non-diagonal + one extra per cycle
       res := k + cycles
       fmt.Fprintln(writer, res)
   }
}
