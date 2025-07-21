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
   // find first zero to remove
   pos := -1
   for i := 0; i < n; i++ {
       if s[i] == '0' {
           pos = i
           break
       }
   }
   var res string
   if pos != -1 {
       res = s[:pos] + s[pos+1:]
   } else {
       // all ones, remove last char
       res = s[:n-1]
   }
   // trim leading zeros
   start := 0
   for start < len(res) && res[start] == '0' {
       start++
   }
   res = res[start:]
   if res == "" {
       res = "0"
   }
   fmt.Fprintln(os.Stdout, res)
}
