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
   heights := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &heights[i])
   }
   // initial sum of first k planks
   sum := 0
   for i := 0; i < k; i++ {
       sum += heights[i]
   }
   minSum := sum
   minPos := 0
   // slide window
   for i := k; i < n; i++ {
       sum += heights[i] - heights[i-k]
       if sum < minSum {
           minSum = sum
           minPos = i - k + 1
       }
   }
   // output 1-based index
   fmt.Println(minPos + 1)
}
