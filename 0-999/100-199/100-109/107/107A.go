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
   var n, p int
   fmt.Fscan(reader, &n, &p)
   to := make([][2]int, n+1)
   out := make([]bool, n+1)
   in := make([]bool, n+1)
   tank := make([]int, n+1)
   ans := 0
   for i := 0; i < p; i++ {
       var a, b, d int
       fmt.Fscan(reader, &a, &b, &d)
       to[a][0] = b
       to[a][1] = d
       out[a] = true
       in[b] = true
   }
   for i := 1; i <= n; i++ {
       if out[i] && !in[i] {
           tank[ans] = i
           ans++
       }
   }
   fmt.Fprintln(writer, ans)
   for i := 0; i < ans; i++ {
       minD := to[tank[i]][1]
       at := to[tank[i]][0]
       for out[at] {
           if minD > to[at][1] {
               minD = to[at][1]
           }
           at = to[at][0]
       }
       fmt.Fprintf(writer, "%d %d %d\n", tank[i], at, minD)
   }
}
