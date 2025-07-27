package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < t; i++ {
       var s string
       fmt.Fscan(reader, &s)
       // simulate removals using a stack
       stack := make([]byte, 0, len(s))
       for j := 0; j < len(s); j++ {
           c := s[j]
           n := len(stack)
           if n > 0 {
               top := stack[n-1]
               if (top == 'A' && c == 'B') || (top == 'B' && c == 'B') {
                   // remove the substring
                   stack = stack[:n-1]
                   continue
               }
           }
           stack = append(stack, c)
       }
       fmt.Fprintln(writer, len(stack))
   }
}
