package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N int64
   if _, err := fmt.Fscan(reader, &N); err != nil {
       return
   }
   // Compute minimal sum of sides using floor(sqrt(N))
   a := math.Floor(math.Sqrt(float64(N)))
   res := math.Ceil(a + float64(N)/a)
   fmt.Fprintln(writer, int64(res))
}
