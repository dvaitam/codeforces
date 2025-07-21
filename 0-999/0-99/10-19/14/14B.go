package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var x0 int
   if _, err := fmt.Fscan(reader, &n, &x0); err != nil {
       return
   }
   // initialize intersection range to full possible track
   const INF = 1000000000
   left, right := 0, INF
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       if a > b {
           a, b = b, a
       }
       if a > left {
           left = a
       }
       if b < right {
           right = b
       }
   }
   // check if intersection exists
   if left > right {
       fmt.Println(-1)
       return
   }
   // compute minimal distance from x0 to [left, right]
   if x0 < left {
       fmt.Println(left - x0)
   } else if x0 > right {
       fmt.Println(x0 - right)
   } else {
       fmt.Println(0)
   }
}
