package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var r, g, b int64
   if _, err := fmt.Fscan(in, &r, &g, &b); err != nil {
       return
   }
   sum := r + g + b
   mx := r
   if g > mx {
       mx = g
   }
   if b > mx {
       mx = b
   }
   other := sum - mx
   var ans int64
   if mx > 2*other {
       ans = other
   } else {
       ans = sum / 3
   }
   fmt.Println(ans)
}
