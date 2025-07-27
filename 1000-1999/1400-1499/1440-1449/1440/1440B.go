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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, k int
       fmt.Fscan(reader, &n, &k)
       total := n * k
       a := make([]int64, total)
       for i := 0; i < total; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // c = ceil(n/2)
       // floor(n/2) = n - c
       skip := n/2 + 1 // floor(n/2)+1 == (n - ceil(n/2)) + 1
       // initial position (1-based): pos = total - floor(n/2)
       // convert to 0-based index:
       idx := total - (n/2) - 1
       var sum int64
       for i := 0; i < k; i++ {
           sum += a[idx]
           idx -= skip
       }
       fmt.Fprintln(writer, sum)
   }
}
