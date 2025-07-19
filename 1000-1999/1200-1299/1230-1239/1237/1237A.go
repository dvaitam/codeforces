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
   f1, f2 := true, true
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       if a > 0 && a%2 != 0 {
           if f1 {
               a++
           }
           f1 = !f1
       }
       if a < 0 && a%2 != 0 {
           if f2 {
               a--
           }
           f2 = !f2
       }
       fmt.Fprintln(writer, a/2)
   }
}
