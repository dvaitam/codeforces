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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       x := 0
       for i := 1; i <= n; i++ {
           var ai int
           fmt.Fscan(reader, &ai)
           if ai == 0 {
               x = i
           }
       }
       // If no zero found, place n+1 at beginning
       if x == 0 {
           fmt.Fprintf(writer, "%d ", n+1)
       }
       // Print sequence with insertion at x
       for i := 1; i <= n; i++ {
           fmt.Fprintf(writer, "%d ", i)
           if i == x {
               fmt.Fprintf(writer, "%d ", n+1)
           }
       }
       fmt.Fprintln(writer)
   }
}
