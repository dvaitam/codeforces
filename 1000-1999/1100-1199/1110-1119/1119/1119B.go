package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var h int64
   if _, err := fmt.Fscan(in, &n, &h); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   result := 0
   // Try prefixes of length i
   for i := 1; i <= n; i++ {
       b := make([]int64, i)
       copy(b, a[:i])
       // sort descending
       sort.Slice(b, func(i, j int) bool {
           return b[i] > b[j]
       })
       var sum int64
       // sum b[0], b[2], b[4], ...
       for j := 0; j < i; j += 2 {
           sum += b[j]
           if sum > h {
               break
           }
       }
       if sum <= h {
           result = i
       } else {
           break
       }
   }
   fmt.Println(result)
}
