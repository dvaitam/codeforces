package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s, t string
   fmt.Fscan(reader, &s)
   fmt.Fscan(reader, &t)
   pos := 0
   for i := 0; i < len(t); i++ {
       if pos < len(s) && s[pos] == t[i] {
           pos++
       }
   }
   // output 1-based position
   fmt.Println(pos + 1)
}
