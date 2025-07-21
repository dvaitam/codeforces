package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   target := "hello"
   idx := 0
   for i := 0; i < len(s) && idx < len(target); i++ {
       if s[i] == target[idx] {
           idx++
       }
   }
   if idx == len(target) {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
