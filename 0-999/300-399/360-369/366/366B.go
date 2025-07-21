package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   sums := make([]int64, k)
   for i := 0; i < n; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       sums[i%k] += a
   }
   minSum := sums[0]
   minR := 0
   for r := 1; r < k; r++ {
       if sums[r] < minSum {
           minSum = sums[r]
           minR = r
       }
   }
   // starting task index is residue+1
   fmt.Println(minR + 1)
}
