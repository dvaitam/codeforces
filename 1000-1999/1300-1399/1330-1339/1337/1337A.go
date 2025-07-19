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
       var a, b, c, d int64
       fmt.Fscan(reader, &a, &b, &c, &d)
       // x = a, y = c, z = c always forms valid triangle since a <= b <= c <= d
       fmt.Fprintln(writer, a, c, c)
   }
}
