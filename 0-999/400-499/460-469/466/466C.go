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
   // Need at least 3 elements to split into three non-empty parts
   if n < 3 {
       fmt.Println(0)
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var total int64
   for _, v := range a {
       total += v
   }
   // total must be divisible by 3
   if total%3 != 0 {
       fmt.Println(0)
       return
   }
   target := total / 3
   var countT int64
   var ans int64
   var prefix int64
   // iterate up to n-1 to ensure third part non-empty
   for i := 0; i < n-1; i++ {
       prefix += a[i]
       // when prefix equals 2*target at position i (as end of second part)
       if i > 0 && prefix == 2*target {
           ans += countT
       }
       // count positions where prefix equals target (possible end of first part)
       if prefix == target {
           countT++
       }
   }
   fmt.Println(ans)
