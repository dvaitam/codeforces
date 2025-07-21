package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   ks := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &ks[i])
   }
   minTime := -1
   for i := 0; i < n; i++ {
       k := ks[i]
       total := 0
       for j := 0; j < k; j++ {
           var m int
           fmt.Fscan(in, &m)
           total += m*5 + 15
       }
       if minTime < 0 || total < minTime {
           minTime = total
       }
   }
   fmt.Println(minTime)
}
