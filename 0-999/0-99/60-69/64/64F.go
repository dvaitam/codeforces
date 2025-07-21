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
   n := len(s)
   // must not start or end with dot
   if n == 0 || s[0] == '.' || s[n-1] == '.' {
       fmt.Println("NO")
       return
   }
   lastDot := -1
   prevDot := false
   for i, c := range s {
       if c == '.' {
           if prevDot {
               fmt.Println("NO")
               return
           }
           prevDot = true
           lastDot = i
       } else {
           // only lowercase letters and digits allowed
           if (c < 'a' || c > 'z') && (c < '0' || c > '9') {
               fmt.Println("NO")
               return
           }
           prevDot = false
       }
   }
   // check last part length
   var lastLen int
   if lastDot == -1 {
       lastLen = n
   } else {
       lastLen = n - lastDot - 1
   }
   if lastLen == 2 || lastLen == 3 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
