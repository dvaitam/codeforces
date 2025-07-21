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
   counts := make([]int, 5)
   for i := 0; i < n; i++ {
       var s int
       fmt.Fscan(reader, &s)
       if s >= 1 && s <= 4 {
           counts[s]++
       }
   }
   taxis := counts[4]
   // Pair groups of 3 with groups of 1
   taxis += counts[3]
   if counts[1] > counts[3] {
       counts[1] -= counts[3]
   } else {
       counts[1] = 0
   }
   // Pair groups of 2 together
   taxis += counts[2] / 2
   if counts[2]%2 != 0 {
       // one group of 2 remains
       taxis++
       // it can take up to two 1s
       if counts[1] > 2 {
           counts[1] -= 2
       } else {
           counts[1] = 0
       }
   }
   // Remaining groups of 1
   if counts[1] > 0 {
       taxis += (counts[1] + 3) / 4
   }
   fmt.Fprintln(os.Stdout, taxis)
}
