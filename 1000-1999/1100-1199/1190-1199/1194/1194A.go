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
       var n, x int64
       fmt.Fscan(reader, &n, &x)
       // After removing all odd-positioned numbers, remaining are evens: 2,4,6,...
       // The x-th remaining number is 2*x.
       fmt.Fprintln(writer, 2*x)
   }
}
