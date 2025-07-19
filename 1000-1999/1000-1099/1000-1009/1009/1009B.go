package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   fmt.Fscan(reader, &s)
   cnt1 := 0
   for i := 0; i < len(s); i++ {
       if s[i] == '1' {
           cnt1++
       }
   }
   var sb strings.Builder
   for i := 0; i < len(s); i++ {
       if s[i] != '1' {
           sb.WriteByte(s[i])
       }
   }
   s2 := sb.String()
   idx := strings.IndexByte(s2, '2')
   if idx == -1 {
       idx = len(s2)
   }
   prefix := s2[:idx]
   rest := s2[idx:]
   writer := bufio.NewWriter(os.Stdout)
   fmt.Fprint(writer, prefix)
   for i := 0; i < cnt1; i++ {
       writer.WriteByte('1')
   }
   fmt.Fprint(writer, rest)
   writer.Flush()
}
