package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Determine stairway breaks: value 1 marks the start of a new stairway
   res := make([]int, 0, n)
   // For each new stairway (except the first), the previous value is the length of the prior stairway
   for i := 2; i <= n; i++ {
       if a[i] == 1 {
           res = append(res, a[i-1])
       }
   }
   // The last stairway ends at the last pronounced number
   if n >= 1 {
       res = append(res, a[n])
   }
   // Output
   fmt.Fprintln(writer, len(res))
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
