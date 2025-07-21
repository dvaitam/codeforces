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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   ans := 0
   for i := 0; i < n; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       if a + b + c >= 2 {
           ans++
       }
   }
   fmt.Fprint(writer, ans)
}
