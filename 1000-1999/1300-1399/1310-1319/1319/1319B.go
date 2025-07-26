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
   // map from key = b_i - i to sum of b_i
   sums := make(map[int]int64, n)
   var best int64
   for i := 1; i <= n; i++ {
       var b int
       fmt.Fscan(reader, &b)
       key := b - i
       sums[key] += int64(b)
       if sums[key] > best {
           best = sums[key]
       }
   }
   fmt.Fprint(writer, best)
}
