package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   ans := 0
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       cnt := 0
       for x > 0 {
           d := x % 10
           if d == 4 || d == 7 {
               cnt++
           }
           x /= 10
       }
       if cnt <= k {
           ans++
       }
   }
   fmt.Println(ans)
}
