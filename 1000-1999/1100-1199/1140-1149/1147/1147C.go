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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   same := true
   first := a[0]
   for i := 1; i < n; i++ {
       if a[i] != first {
           same = false
           break
       }
   }
   if same {
       fmt.Fprintln(writer, "Bob")
   } else {
       fmt.Fprintln(writer, "Alice")
   }
}
