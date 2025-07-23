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
   const target = "CODEFORCES"
   n := len(s)
   m := len(target)
   if n < m {
       fmt.Println("NO")
       return
   }
   for i := 0; i <= m; i++ {
       // Check prefix of length i and suffix of length m-i
       if s[:i] == target[:i] && s[n-(m-i):] == target[i:] {
           fmt.Println("YES")
           return
       }
   }
   fmt.Println("NO")
}
