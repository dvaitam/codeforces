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
   for i := 0; i < t; i++ {
       var n string
       fmt.Fscan(reader, &n)
       // The number of distinct g(x) values is the maximum k such that 10^k <= n, plus 1
       // That equals the number of digits in n
       fmt.Fprintln(writer, len(n))
   }
}
