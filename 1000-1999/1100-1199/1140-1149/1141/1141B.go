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
   cnt, p, mx := 0, 0, 0
   firstZero := true
   for i := 0; i < n; i++ {
       var d int
       if _, err := fmt.Fscan(reader, &d); err != nil {
           return
       }
       if d == 1 {
           cnt++
           if cnt > mx {
               mx = cnt
           }
       } else {
           if firstZero {
               p = cnt
               firstZero = false
           }
           cnt = 0
       }
   }
   if cnt+p > mx {
       mx = cnt + p
   }
   fmt.Fprint(writer, mx)
}
