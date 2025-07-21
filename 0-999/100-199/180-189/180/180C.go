package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, "Error reading input:", err)
       return
   }
   s := strings.TrimSpace(line)
   n := len(s)
   totalUpper := 0
   for i := 0; i < n; i++ {
       if s[i] >= 'A' && s[i] <= 'Z' {
           totalUpper++
       }
   }

   prefixLower := 0
   prefixUpper := 0
   // ans: minimum operations
   ans := n // worst case: change all
   for i := 0; i <= n; i++ {
       // operations: change lowercase in prefix to upper, and uppercase in suffix to lower
       ops := prefixLower + (totalUpper - prefixUpper)
       if ops < ans {
           ans = ops
       }
       if i < n {
           if s[i] >= 'a' && s[i] <= 'z' {
               prefixLower++
           } else {
               prefixUpper++
           }
       }
   }
   fmt.Println(ans)
}
