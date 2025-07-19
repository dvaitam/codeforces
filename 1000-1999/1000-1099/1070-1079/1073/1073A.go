package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var s string
   fmt.Fscan(reader, &n, &s)
   ch := s[0]
   found := false
   var a byte
   for i := 1; i < n; i++ {
       if s[i] != ch {
           a = s[i]
           found = true
           break
       }
   }
   if found {
       fmt.Println("YES")
       fmt.Printf("%c%c\n", ch, a)
   } else {
       fmt.Println("NO")
   }
}
