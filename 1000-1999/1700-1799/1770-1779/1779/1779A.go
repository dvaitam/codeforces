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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       var s string
       fmt.Fscan(reader, &n, &s)
       ans := 0
       fl := true
       for i := 0; i+1 < n; i++ {
           if s[i] == 'R' && s[i+1] == 'L' {
               ans = 0
               fl = false
               break
           }
           if s[i] == 'L' && s[i+1] == 'R' && fl {
               ans = i + 1
               fl = false
           }
       }
       if fl {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, ans)
       }
   }
}
