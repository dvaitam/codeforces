package main

import (
   "bufio"
   "fmt"
   "os"
)

func isFair(x uint64) bool {
   y := x
   for y > 0 {
       d := y % 10
       if d != 0 && x % d != 0 {
           return false
       }
       y /= 10
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       var n uint64
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
