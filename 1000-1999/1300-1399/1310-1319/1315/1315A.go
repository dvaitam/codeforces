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
       var a, b, x, y int
       fmt.Fscan(reader, &a, &b, &x, &y)

       left := x * b
       right := (a - x - 1) * b
       top := y * a
       bottom := (b - y - 1) * a

       ans := left
       if right > ans {
           ans = right
       }
       if top > ans {
           ans = top
       }
       if bottom > ans {
           ans = bottom
       }

       fmt.Fprintln(writer, ans)
   }
}
