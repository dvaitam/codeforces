package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var a, b, c int64
   if _, err := fmt.Fscan(in, &a, &b, &c); err != nil {
       return
   }
   counts := []int64{a, b, c}
   ans := int64(-1)
   // Try each color as target
   for i := 0; i < 3; i++ {
       u := counts[(i+1)%3]
       v := counts[(i+2)%3]
       // parity condition: only reachable if sum of other two is even
       if (u+v)&1 == 1 {
           continue
       }
       // minimal fights to unify to this target is max(u, v)
       f := u
       if v > u {
           f = v
       }
       if ans == -1 || f < ans {
           ans = f
       }
   }
   fmt.Fprintln(os.Stdout, ans)
}
