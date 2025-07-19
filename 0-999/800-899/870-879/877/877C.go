package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Read n
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }

   // Calculate number of bombs
   var m int
   if n%2 == 0 {
       m = n/2*3
   } else {
       m = n/2*3 + 1
   }

   // Prepare sequence
   var sb strings.Builder
   sb.Grow(m * 4) // approximate

   if n%2 == 0 {
       // odd, even, odd
       for i := 1; i <= n; i += 2 {
           sb.WriteString(strconv.Itoa(i))
           sb.WriteByte(' ')
       }
       for i := 2; i <= n; i += 2 {
           sb.WriteString(strconv.Itoa(i))
           sb.WriteByte(' ')
       }
       for i := 1; i <= n; i += 2 {
           sb.WriteString(strconv.Itoa(i))
           sb.WriteByte(' ')
       }
   } else {
       // even, odd, even
       for i := 2; i <= n; i += 2 {
           sb.WriteString(strconv.Itoa(i))
           sb.WriteByte(' ')
       }
       for i := 1; i <= n; i += 2 {
           sb.WriteString(strconv.Itoa(i))
           sb.WriteByte(' ')
       }
       for i := 2; i <= n; i += 2 {
           sb.WriteString(strconv.Itoa(i))
           sb.WriteByte(' ')
       }
   }

   // Output
   fmt.Fprintln(writer, m)
   writer.WriteString(sb.String())
}
