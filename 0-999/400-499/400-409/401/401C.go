package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var n, m int
   if _, err := fmt.Fscan(os.Stdin, &n, &m); err != nil {
       return
   }
   // Check feasibility: no more than one consecutive zero, no more than two consecutive ones
   if n > m+1 || n*2+2 < m {
       fmt.Print(-1)
       return
   }
   p, k := 1, 0
   out := make([]byte, 0, n+m)
   for n > 0 || m > 0 {
       if p == 0 || (m > n && k < 2) {
           m--
           p = 1
           k++
       } else {
           n--
           p = 0
           k = 0
       }
       // append '0' or '1'
       out = append(out, "01"[p])
   }
   w := bufio.NewWriter(os.Stdout)
   w.Write(out)
   w.Flush()
}
