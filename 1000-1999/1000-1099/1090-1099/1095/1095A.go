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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   ans := make([]byte, 0, len(s))
   step := 1
   for i := 0; i < len(s); i += step {
       ans = append(ans, s[i])
       step++
   }
   // Output result
   fmt.Println(string(ans))
}
