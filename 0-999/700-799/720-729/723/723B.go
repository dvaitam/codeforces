package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)

   inside := false
   current := 0
   maxOutside := 0
   countInside := 0

   for i := 0; i < len(s); i++ {
       ch := s[i]
       if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
           current++
       } else {
           if current > 0 {
               if inside {
                   countInside++
               } else if current > maxOutside {
                   maxOutside = current
               }
               current = 0
           }
           if ch == '(' {
               inside = true
           } else if ch == ')' {
               inside = false
           }
       }
   }
   // process last word
   if current > 0 {
       if inside {
           countInside++
       } else if current > maxOutside {
           maxOutside = current
       }
   }
   fmt.Printf("%d %d\n", maxOutside, countInside)
}
