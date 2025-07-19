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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // Check feasibility
   if n%2 == 0 {
       if k > n/2-1 {
           fmt.Fprint(writer, -1)
           return
       }
   } else {
       if k > n/2 {
           fmt.Fprint(writer, -1)
           return
       }
   }
   // Total pairs
   total := n * k
   fmt.Fprintln(writer, total)
   // Generate pairs
   for i := 1; i <= n; i++ {
       for j := i + 1; j <= i+k; j++ {
           if j > n {
               t := j - n
               fmt.Fprintf(writer, "%d %d\n", i, t)
           } else {
               fmt.Fprintf(writer, "%d %d\n", i, j)
           }
       }
   }
}
