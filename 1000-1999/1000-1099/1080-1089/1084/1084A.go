package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs64(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   arr := make([]int64, n+1)
   for i := int64(1); i <= n; i++ {
       fmt.Fscan(in, &arr[i])
   }
   best := int64(1<<63 - 1)
   for j := int64(1); j <= n; j++ {
       var sum int64
       for k := int64(1); k <= n; k++ {
           // two trips per person: down and up
           cost := 2*abs64(k-j) + 2*(k-1) + 2*(j-1)
           sum += arr[k] * cost
       }
       if sum < best {
           best = sum
       }
   }
   fmt.Println(best)
}
