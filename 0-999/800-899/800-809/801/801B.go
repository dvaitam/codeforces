package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x, y string
   if _, err := fmt.Fscan(reader, &x); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &y); err != nil {
       return
   }
   if len(x) != len(y) {
       fmt.Println(-1)
       return
   }
   for i := 0; i < len(x); i++ {
       if x[i] < y[i] {
           fmt.Println(-1)
           return
       }
   }
   // any valid z is y itself
   fmt.Print(y)
}
