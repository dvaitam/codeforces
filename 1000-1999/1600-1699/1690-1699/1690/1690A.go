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

   var t, n int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       fmt.Fscan(reader, &n)
       a := (n + 1) / 3
       b := (n + 2) / 3 + 1
       c := n/3 - 1
       fmt.Fprintln(writer, a, b, c)
   }
}
