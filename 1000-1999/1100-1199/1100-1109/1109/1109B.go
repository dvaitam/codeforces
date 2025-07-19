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
   // Check impossible: all chars same or only middle differs
   ok := false
   for i := 0; i < n; i++ {
       if s[i] != s[0] && i != n-1-i {
           ok = true
           break
       }
   }
   if !ok {
       fmt.Println("Impossible")
       return
   }
   ans := 2
   t := s + s
   for i := 1; i < n; i++ {
       u := t[i : i+n]
       if u == s {
           continue
       }
       pal := true
       for j := 0; j < n-1-j; j++ {
           if u[j] != u[n-1-j] {
               pal = false
               break
           }
       }
       if pal {
           ans = 1
           break
       }
   }
   fmt.Println(ans)
}
