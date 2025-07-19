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

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   vis := make([]int, n+1)
   var ss int
   for i := 0; i < k; i++ {
       var x int
       fmt.Fscan(reader, &x)
       vis[x] = 1
       if i == 0 {
           vis[x] = 2
           ss = x
       }
   }
   a := make([]int, 0, n)
   b := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if vis[i] != 2 {
           a = append(a, i)
       }
       if vis[i] == 0 {
           b = append(b, i)
       }
   }
   tot := (n-k) + ((n-1)*(n-2))/2
   if m > tot || k == n {
       fmt.Fprintln(writer, -1)
   } else {
       // print edges among non-root nodes until one edge remains
       for i := 0; i < len(a) && m > 1; i++ {
           for j := i + 1; j < len(a) && m > 1; j++ {
               fmt.Fprintln(writer, a[i], a[j])
               m--
           }
       }
       // connect root to non-special nodes for remaining edges
       for i := 0; i < m; i++ {
           fmt.Fprintln(writer, ss, b[i])
       }
   }
}
