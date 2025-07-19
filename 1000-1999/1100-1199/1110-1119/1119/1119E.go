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
   var res int64
   var ans int64
   for i := 0; i < n; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       // use as many pairs of this x and previous leftovers
       tmpPairs := res
       half := x / 2
       if half < tmpPairs {
           tmpPairs = half
       }
       // remaining x after using pairs
       xRem := x - tmpPairs*2
       // use triples of x
       tmpTriples := xRem / 3
       tmp := tmpPairs + tmpTriples
       ans += tmp
       // update leftovers: previous plus current leftover sticks
       res = res + x - 3*tmp
   }
   fmt.Println(ans)
}
