package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   if 2*k < n || k == n {
       fmt.Fprintln(writer, -1)
       return
   }
   half := n >> 1
   prefix := (k - half + 1) >> 1
   var sb strings.Builder
   for i := 0; i < prefix; i++ {
       sb.WriteString("()")
   }
   for i := 0; i < n-k; i++ {
       sb.WriteByte('(')
   }
   for i := 0; i < n-k; i++ {
       sb.WriteByte(')')
   }
   suffix := (k - half) >> 1
   for i := 0; i < suffix; i++ {
       sb.WriteString("()")
   }
   sb.WriteByte('\n')
   writer.WriteString(sb.String())
}
