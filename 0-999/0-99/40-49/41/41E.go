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
   a := n >> 1
   b := n - a
   fmt.Fprintln(writer, a*b)
   for i := 0; i < a; i++ {
       for j := a; j < n; j++ {
           fmt.Fprintln(writer, i+1, j+1)
       }
   }
}
