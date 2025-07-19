package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var l int
   if _, err := fmt.Fscan(os.Stdin, &l); err != nil {
       return
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   if l%2 != 0 {
       w.WriteString("-1\n")
       return
   }
   for i := 0; i < l; i++ {
       for j := 0; j < l; j++ {
           for k := 0; k < l; k++ {
               v := ((j/2)&1) ^ ((k/2)&1) ^ (i & 1)
               if v != 0 {
                   w.WriteByte('w')
               } else {
                   w.WriteByte('b')
               }
           }
           w.WriteByte('\n')
       }
       w.WriteByte('\n')
   }
}
