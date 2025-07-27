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
       var a, b int64
       fmt.Fscan(reader, &a, &b)
       d := a - b
       if d < 0 {
           d = -d
       }
       moves := d / 10
       if d%10 != 0 {
           moves++
       }
       fmt.Fprintln(writer, moves)
   }
}
