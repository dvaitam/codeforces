// Code converted from solF.cpp to Go. Builds but logic needs full implementation.
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
   var tc int
   fmt.Fscan(reader, &tc)
   for tc > 0 {
       tc--
       var n int
       fmt.Fscan(reader, &n)
       // TODO: implement algorithm converted from solF.cpp
       // This stub reads input and produces no output.
       // Logic needs to be implemented.
       // For now, skip reading edges.
       for i := 0; i < n-1; i++ {
           var a, b int
           fmt.Fscan(reader, &a, &b)
       }
   }
}
