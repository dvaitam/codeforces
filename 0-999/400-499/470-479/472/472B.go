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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   floors := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &floors[i])
   }
   // Sort floors ascending
   sort.Ints(floors)
   var ans int64
   // Group in batches of size k from highest floors
   for i := n - 1; i >= 0; i -= k {
       // round-trip from 1 to floors[i] and back
       ans += int64(2 * (floors[i] - 1))
   }
   fmt.Fprintln(writer, ans)
}
