package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   good := false
   for d := 1; d*4 < n && !good; d++ {
       for i := 0; i+4*d < n; i++ {
           ok := true
           for k := 0; k < 5; k++ {
               if s[i+k*d] != '*' {
                   ok = false
                   break
               }
           }
           if ok {
               good = true
               break
           }
       }
   }
   if good {
       fmt.Println("yes")
   } else {
       fmt.Println("no")
   }
}
