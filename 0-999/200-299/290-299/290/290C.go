package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var sum, tmp, ans float64
   for i := 1; i <= n; i++ {
       if _, err := fmt.Fscan(reader, &tmp); err != nil {
           return
       }
       sum += tmp
       avg := sum / float64(i)
       if avg > ans {
           ans = avg
       }
   }
   fmt.Printf("%.20f\n", ans)
}
