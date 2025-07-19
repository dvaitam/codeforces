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
   result := ""
   for len(s) > 0 {
       l := len(s)
       var idx int
       if l%2 == 0 {
           idx = l/2 - 1
       } else {
           idx = l/2
       }
       result += string(s[idx])
       s = s[:idx] + s[idx+1:]
   }
   fmt.Print(result)
}
