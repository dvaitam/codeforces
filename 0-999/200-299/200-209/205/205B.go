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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   prev := int64(0)
   var cur int64
   var moves int64
   if n > 0 {
       // read first element
       fmt.Fscan(reader, &prev)
   }
   for i := 2; i <= n; i++ {
       fmt.Fscan(reader, &cur)
       if prev > cur {
           moves += prev - cur
       } else {
           prev = cur
       }
       // after increasing cur to prev, cur becomes prev
       // so prev remains same
   }
   fmt.Fprintln(writer, moves)
}
