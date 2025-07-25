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
   var max1, max2, cnt int
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(in, &x)
       if i >= 2 && x < max2 {
           cnt++
       }
       if x > max1 {
           max2 = max1
           max1 = x
       } else if x > max2 {
           max2 = x
       }
   }
   fmt.Println(cnt)
}
