package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   ans := int64(0)
   for i := 0; i < n; i++ {
       var k, a int64
       fmt.Fscan(in, &k, &a)
       // compute minimal d such that 4^d >= a
       d := int64(0)
       cap := int64(1)
       for cap < a {
           cap <<= 2 // multiply by 4
           d++
       }
       if k+d > ans {
           ans = k + d
       }
   }
   fmt.Println(ans)
}
