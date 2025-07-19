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
   var ans float64
   for i := 1; i <= n; i++ {
       ans += 1.0 / float64(i)
   }
   fmt.Printf("%.12f\n", ans)
}
