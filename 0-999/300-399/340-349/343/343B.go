package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   cplus, cminus := 0, 0
   for _, ch := range s {
       if ch == '+' {
           cplus++
       } else if ch == '-' {
           cminus++
       }
   }
   net := cplus - cminus
   yes := false
   switch {
   case net < 0:
       yes = false
   case net == 0:
       yes = true
   default: // net > 0
       if net%4 != 0 {
           yes = true
       }
   }
   if yes {
       fmt.Println("Yes")
   } else {
       fmt.Println("No")
   }
}
