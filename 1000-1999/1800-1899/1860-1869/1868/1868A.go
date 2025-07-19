package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   in  = bufio.NewReader(os.Stdin)
   out = bufio.NewWriter(os.Stdout)
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func solve() {
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   if m == 1 {
       fmt.Fprintln(out, 0)
   } else if n > m-1 {
       fmt.Fprintln(out, m)
   } else {
       fmt.Fprintln(out, n+1)
   }
   for i := 0; i < min(m-1, n); i++ {
       for j := 0; j < m; j++ {
           fmt.Fprint(out, (j+i)%m, " ")
       }
       fmt.Fprintln(out)
   }
   if n > m-1 {
       for i := m-1; i < n; i++ {
           for j := 0; j < m; j++ {
               fmt.Fprint(out, j, " ")
           }
           fmt.Fprintln(out)
       }
   }
}

func main() {
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       solve()
       t--
   }
}
