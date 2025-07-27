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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for i := 0; i < T; i++ {
       var n int64
       fmt.Fscan(reader, &n)
       res := (n + 1) / 2
       fmt.Fprintln(writer, res)
   }
}
