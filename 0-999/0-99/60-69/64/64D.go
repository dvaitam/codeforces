package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var prices [3]int
   for i := 0; i < 3; i++ {
       if _, err := fmt.Fscan(reader, &prices[i]); err != nil {
           return
       }
   }
   // Pair prices with original indices
   presents := make([]struct { price, idx int }, 3)
   for i := 0; i < 3; i++ {
       presents[i] = struct{ price, idx int }{prices[i], i}
   }
   // Sort descending by price
   sort.Slice(presents, func(i, j int) bool {
       return presents[i].price > presents[j].price
   })
   // Map each present to sister number: highest -> 1, next -> 2, lowest -> 3
   ans := make([]int, 3)
   for rank, p := range presents {
       // sister numbers are 1-based by rank
       ans[p.idx] = rank + 1
   }
   // Output mapping for each present in original order
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%d %d %d", ans[0], ans[1], ans[2])
}
