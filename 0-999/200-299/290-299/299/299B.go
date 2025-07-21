package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       fmt.Fprintln(os.Stderr, "failed to read n and k:", err)
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       fmt.Fprintln(os.Stderr, "failed to read road string:", err)
       return
   }
   cnt := 0
   for i := 0; i < n; i++ {
       if s[i] == '#' {
           cnt++
           if cnt >= k {
               fmt.Println("NO")
               return
           }
       } else {
           cnt = 0
       }
   }
   fmt.Println("YES")
}
