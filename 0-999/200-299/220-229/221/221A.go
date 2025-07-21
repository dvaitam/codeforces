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
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if n >= 1 {
       // To satisfy the recursive swap sequence, output n followed by 1 to n-1
       fmt.Fprintf(writer, "%d", n)
       for i := 1; i < n; i++ {
           fmt.Fprintf(writer, " %d", i)
       }
   }
   fmt.Fprintln(writer)
}
