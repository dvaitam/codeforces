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
   switch n {
   case 1:
       fmt.Fprint(writer, "1")
       return
   case 3:
       fmt.Fprint(writer, "1 1 3")
       return
   }

   x := n
   s := 1
   num := 0
   for x > 0 {
       x -= num
       num = (x + 1) >> 1
       if x == 1 {
           s >>= 1
           val := n - n%s
           fmt.Fprint(writer, val)
       } else {
           for i := 0; i < num; i++ {
               fmt.Fprint(writer, s)
               writer.WriteByte(' ')
           }
       }
       s <<= 1
   }
}
