package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var x, y int64
   if _, err := fmt.Fscan(reader, &x, &y); err != nil {
       return
   }
   // Starting from (x,x,x) to (y,y,y)
   // If x == y, zero operations (though per problem y < x always)
   if x == y {
       fmt.Println(0)
       return
   }
   // Build sequence f: f[0]=y, f[1]=2*y-1, f[n]=f[n-1]+f[n-2]-1
   // Find minimal t such that f[t] >= x, t>=1
   prev := y
   curr := 2*y - 1
   t := int64(1)
   for curr < x {
       next := curr + prev - 1
       prev = curr
       curr = next
       t++
   }
   // Total operations: first increase side to curr then two more to set other sides
   // Reverse operations count: t + 2
   fmt.Println(t + 2)
}
