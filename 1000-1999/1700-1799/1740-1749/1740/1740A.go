package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for i := 0; i < t; i++ {
       var n int
       fmt.Fscan(in, &n)
       fmt.Fprintln(out, 7)
   }
}
