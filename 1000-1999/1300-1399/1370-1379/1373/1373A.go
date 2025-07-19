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
       var a, b, c int64
       fmt.Fscan(reader, &a, &b, &c)
       x, y := solve(a, b, c)
       fmt.Fprintln(writer, x, y)
   }
}

func solve(a, b, c int64) (int64, int64) {
   if a*b <= c {
       if a == c && b == 1 {
           return -1, -1
       }
       return 1, -1
   }
   if a >= c {
       return -1, b
   }
   return 1, b
}
