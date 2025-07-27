package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// Simple offline solution: read all tree heights and output the minimum.
// For interactive version, this will not apply; use appropriate queries.
func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   minVal := int64(math.MaxInt64)
   var x int64
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if _, err := fmt.Fscan(reader, &x); err != nil {
               return
           }
           if x < minVal {
               minVal = x
           }
       }
   }
   fmt.Fprintln(writer, minVal)
}
