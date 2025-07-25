package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)

   // Build the longest good string
   ans := make([]byte, 0, n)
   for i := 0; i < len(s); i++ {
       if len(ans)%2 == 0 {
           ans = append(ans, s[i])
       } else if ans[len(ans)-1] != s[i] {
           ans = append(ans, s[i])
       }
   }
   // Ensure even length
   if len(ans)%2 == 1 {
       ans = ans[:len(ans)-1]
   }

   deletions := n - len(ans)
   fmt.Fprintln(writer, deletions)
   if len(ans) > 0 {
       fmt.Fprintln(writer, string(ans))
   }
}
