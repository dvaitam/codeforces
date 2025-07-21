package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var name string
   if _, err := fmt.Fscan(reader, &name); err != nil {
       return
   }
   var n int64
   var h int
   fmt.Fscan(reader, &n, &h)
   // Approximate expected values: E[Bob|Alice=n] = n + 2^{h+1} - 2
   // and E[Alice|Bob=n] = n - (2^{h+1} - 2)
   // Compute offset = 2^{h+1} - 2
   var offset int64 = (1 << (h + 1)) - 2
   var ans float64
   if name == "Alice" {
       ans = float64(n + offset)
   } else {
       ans = float64(n - offset)
   }
   // Print with sufficient precision
   fmt.Printf("%.10f\n", ans)
}
