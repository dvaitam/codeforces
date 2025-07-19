package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
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
       if n == 1 {
           fmt.Fprintln(writer, "a")
           continue
       }
       half := n / 2
       // Build the string
       s := strings.Repeat("a", half) + "b" + strings.Repeat("a", half-1)
       if n%2 != 0 {
           s += "c"
       }
       fmt.Fprintln(writer, s)
   }
}
