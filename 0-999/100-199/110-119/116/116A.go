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
   curr := 0
   maxc := 0
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       curr -= a
       curr += b
       if curr > maxc {
           maxc = curr
       }
   }
   fmt.Println(maxc)
}
