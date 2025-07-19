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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   ans := 1e18
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       cur := float64(m) * float64(a) / float64(b)
       if cur < ans {
           ans = cur
       }
   }
   fmt.Fprintf(writer, "%.20f\n", ans)
}
