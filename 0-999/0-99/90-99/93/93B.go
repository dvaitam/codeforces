package main

import (
   "bufio"
   "fmt"
   "os"
)

const outFmt = "%.16f"

func main() {
   r := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   var n, W, m int
   if _, err := fmt.Fscan(r, &n, &W, &m); err != nil {
       return
   }
   if n < m && m%(m-n) > 0 {
       fmt.Fprintln(w, "NO")
       return
   }
   fmt.Fprintln(w, "YES")
   cur, used := 1, 0
   for i := 0; i < m; i++ {
       sum := 0
       printed := false
       for sum < n {
           cnt := n - sum
           if m-used < cnt {
               cnt = m - used
           }
           if printed {
               fmt.Fprint(w, " ")
           }
           length := float64(cnt)/float64(m) * float64(W)
           fmt.Fprintf(w, "%d "+outFmt, cur, length)
           printed = true
           used += cnt
           sum += cnt
           if used == m {
               cur++
               used = 0
           }
       }
       fmt.Fprintln(w)
   }
}
