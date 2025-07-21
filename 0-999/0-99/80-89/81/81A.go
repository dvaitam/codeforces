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
   stack := make([]byte, 0, len(s))
   for i := 0; i < len(s); i++ {
       c := s[i]
       n := len(stack)
       if n > 0 && stack[n-1] == c {
           stack = stack[:n-1]
       } else {
           stack = append(stack, c)
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   writer.Write(stack)
   writer.WriteByte('\n')
}
