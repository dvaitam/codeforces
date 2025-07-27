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
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       var n int64
       fmt.Fscan(reader, &n)
       x := n
       for {
           if isFair(x) {
               fmt.Fprintln(writer, x)
               break
           }
           x++
       }
   }
}

// isFair returns true if x is divisible by each of its non-zero digits.
func isFair(x int64) bool {
   y := x
   for y > 0 {
       d := y % 10
       y /= 10
       if d != 0 && x%d != 0 {
           return false
       }
   }
   return true
}
