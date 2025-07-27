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
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       to := make([]int, n+1)
       k := 0
       for i := 0; i < m; i++ {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           if x != y {
               to[x] = y
               k++
           }
       }
       state := make([]int, n+1) // 0=unvisited,1=visiting,2=visited
       cycles := 0
       for u := 1; u <= n; u++ {
           if state[u] != 0 || to[u] == 0 {
               continue
           }
           // detect cycle starting from u
           v := u
           for state[v] == 0 {
               state[v] = 1
               nxt := to[v]
               if nxt == 0 {
                   break
               }
               v = nxt
           }
           if state[v] == 1 {
               cycles++
           }
           // mark component as fully visited
           w := u
           for state[w] == 1 {
               state[w] = 2
               w = to[w]
               if w == 0 {
                   break
               }
           }
       }
       // result is non-diagonal rooks + number of cycles
       result := k + cycles
       fmt.Fprintln(writer, result)
   }
}
