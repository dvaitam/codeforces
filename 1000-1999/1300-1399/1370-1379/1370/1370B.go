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

   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       even := make([]int, 0, n)
       odd := make([]int, 0, n)
       for i := 1; i <= 2*n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if x%2 == 0 {
               even = append(even, i)
           } else {
               odd = append(odd, i)
           }
       }
       pairs := n - 1
       cnt := 0
       for i := 0; i+1 < len(even) && cnt < pairs; i += 2 {
           fmt.Fprintf(writer, "%d %d\n", even[i], even[i+1])
           cnt++
       }
       for i := 0; i+1 < len(odd) && cnt < pairs; i += 2 {
           fmt.Fprintf(writer, "%d %d\n", odd[i], odd[i+1])
           cnt++
       }
   }
}
