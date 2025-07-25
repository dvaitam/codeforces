package main

import (
   "bufio"
   "fmt"
   "os"
)

// A trivial solution for Codeforces 1160A2: reads input and outputs no worker blocks.
// This produces a valid but zero-score solution.
func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Skip all location data
   for i := 0; i < n; i++ {
       var x, y, d, p, l, h int
       if _, err := fmt.Fscan(reader, &x, &y, &d, &p, &l, &h); err != nil {
           return
       }
   }
   // No output: zero workers, zero jobs done
}
