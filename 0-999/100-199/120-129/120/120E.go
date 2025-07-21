package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   var sb strings.Builder
   for i := 0; i < T; i++ {
       var n int
       fmt.Fscan(reader, &n)
       if n%2 == 1 {
           sb.WriteByte('0')
       } else {
           sb.WriteByte('1')
       }
       if i < T-1 {
           sb.WriteByte('\n')
       }
   }
   fmt.Fprint(os.Stdout, sb.String())
}
