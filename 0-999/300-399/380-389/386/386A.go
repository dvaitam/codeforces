package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   bids := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &bids[i])
   }
   // Find maximum and second maximum
   maxVal, secondVal := -1, -1
   maxIdx := -1
   for i, v := range bids {
       if v > maxVal {
           secondVal = maxVal
           maxVal = v
           maxIdx = i + 1 // 1-based index
       } else if v > secondVal {
           secondVal = v
       }
   }
   // Output winner index and price to pay (second highest bid)
   fmt.Printf("%d %d", maxIdx, secondVal)
}
