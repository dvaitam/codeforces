package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var q int
   if _, err := fmt.Fscan(in, &q); err != nil {
       return
   }
   for i := 0; i < q; i++ {
       var n int64
       fmt.Fscan(in, &n)
       cnt2, cnt3, cnt5 := 0, 0, 0
       for n%2 == 0 {
           cnt2++
           n /= 2
       }
       for n%3 == 0 {
           cnt3++
           n /= 3
       }
       for n%5 == 0 {
           cnt5++
           n /= 5
       }
       if n != 1 {
           fmt.Fprintln(out, -1)
       } else {
           // Each division by 2: 1 move
           // Each division by 3 (2n/3): removes one 3, adds one 2: counts as 2 moves total (1 for op, 1 for eventual /2)
           // Each division by 5 (4n/5): removes one 5, adds two 2s: counts as 3 moves total
           res := cnt2 + 2*cnt3 + 3*cnt5
           fmt.Fprintln(out, res)
       }
   }
}
