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
   arr := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }
   var total int64
   for _, v := range arr {
       total += v
   }
   if total%3 != 0 {
       fmt.Println(0)
       return
   }
   target := total / 3
   var prefix int64
   var countFirst int64
   var ways int64
   // iterate up to n-1 to leave non-empty third segment
   for i := 0; i < n-1; i++ {
       prefix += arr[i]
       if prefix == 2*target {
           ways += countFirst
       }
       if prefix == target {
           countFirst++
       }
   }
   fmt.Println(ways)
}
