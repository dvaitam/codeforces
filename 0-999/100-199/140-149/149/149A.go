package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var k int
   fmt.Fscan(reader, &k)
   // Read growth values for 12 months
   a := make([]int, 12)
   for i := 0; i < 12; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // If no growth needed
   if k == 0 {
       fmt.Fprintln(writer, 0)
       return
   }
   // Sort months by growth in descending order
   sort.Slice(a, func(i, j int) bool {
       return a[i] > a[j]
   })
   sum := 0
   for i, v := range a {
       sum += v
       if sum >= k {
           // i+1 months needed
           fmt.Fprintln(writer, i+1)
           return
       }
   }
   // Not possible to reach k
   fmt.Fprintln(writer, -1)
}
