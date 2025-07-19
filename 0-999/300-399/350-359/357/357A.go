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
   a := make([]int, n)
   var sum1 int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       sum1 += int64(a[i])
   }
   var low, high int64
   fmt.Fscan(reader, &low, &high)
   var sum2 int64
   for i := 0; i < n; i++ {
       sum2 += int64(a[i])
       sum1 -= int64(a[i])
       if sum2 >= low && sum2 <= high && sum1 >= low && sum1 <= high {
           // partition between i and i+1, second part starts at i+2 (1-based)
           fmt.Println(i + 2)
           return
       }
   }
   fmt.Println(0)
}
