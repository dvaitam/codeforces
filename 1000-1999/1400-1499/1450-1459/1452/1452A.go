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
       var x, y int
       fmt.Fscan(reader, &x, &y)
       a, b := x, y
       if b > a {
           a, b = b, a
       }
       if a == b {
           fmt.Fprintln(writer, 2*a)
       } else {
           fmt.Fprintln(writer, 2*a-1)
       }
   }
