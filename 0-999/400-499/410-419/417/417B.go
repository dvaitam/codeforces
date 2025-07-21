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
   const maxK = 100000
   // mex[k] is the next expected new solution index for participant k
   mex := make([]int, maxK+1)
   for i := 0; i < n; i++ {
       var x, k int
       fmt.Fscan(reader, &x, &k)
       // participants are 1-indexed up to maxK
       cur := mex[k]
       if x > cur {
           fmt.Fprintln(writer, "NO")
           return
       }
       if x == cur {
           // first time seeing this solution index, increment mex
           mex[k] = cur + 1
       }
       // else x < cur: duplicate of previous solution, ok
   }
   fmt.Fprintln(writer, "YES")
}
